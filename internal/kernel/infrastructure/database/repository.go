package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var ErrNotFound = errors.New("database: no rows found")

// Scanner abstracts *sql.Row and *sql.Rows so one mapper serves both read paths.
type Scanner interface {
	Scan(dest ...any) error
}

// Executor is satisfied by both *sql.DB and *sql.Tx, which is what makes every
// repository usable inside a transaction without duplicating its methods.
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

var (
	_ Executor = (*sql.DB)(nil)
	_ Executor = (*sql.Tx)(nil)
)

// Mapping carries the bounded-context knowledge the generic repository needs about T.
type Mapping[T any] struct {
	Table     string   // optionally schema-qualified, e.g. "wallet.CardProjection"
	IDColumns []string // the key columns, in the order GetByKey expects them
	Columns   []string // every persisted column, including the key
	Scan      func(Scanner) (T, error)
	Values    func(T) []any // aligned with Columns
	ID        func(T) []any // aligned with IDColumns
}

// Repository is a generic MSSQL-backed read-side store, reused by every bounded context.
type Repository[T any] struct {
	exec    Executor
	mapping Mapping[T]
}

func NewRepository[T any](exec Executor, mapping Mapping[T]) *Repository[T] {
	return &Repository[T]{exec: exec, mapping: mapping}
}

// WithTx returns a copy bound to exec; the receiver stays untouched because
// repositories are shared across concurrent requests.
func (r *Repository[T]) WithTx(exec Executor) *Repository[T] {
	return &Repository[T]{exec: exec, mapping: r.mapping}
}

// Executor exposes the current executor for queries that do not map onto T.
func (r *Repository[T]) Executor() Executor {
	return r.exec
}

func (r *Repository[T]) ExecuteQuery(ctx context.Context, query string, args ...any) ([]T, error) {
	rows, err := r.exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database: query: %w", err)
	}
	defer rows.Close()

	var result []T
	for rows.Next() {
		entity, err := r.mapping.Scan(rows)
		if err != nil {
			return nil, fmt.Errorf("database: scan row: %w", err)
		}
		result = append(result, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("database: iterate rows: %w", err)
	}

	return result, nil
}

func (r *Repository[T]) ExecuteQuerySingle(ctx context.Context, query string, args ...any) (T, error) {
	var zero T

	entity, err := r.mapping.Scan(r.exec.QueryRowContext(ctx, query, args...))
	if errors.Is(err, sql.ErrNoRows) {
		return zero, ErrNotFound
	}
	if err != nil {
		return zero, fmt.Errorf("database: query single: %w", err)
	}

	return entity, nil
}

func (r *Repository[T]) ExecuteCommand(ctx context.Context, query string, args ...any) (int64, error) {
	result, err := r.exec.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("database: exec: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("database: rows affected: %w", err)
	}

	return affected, nil
}

func (r *Repository[T]) GetByKey(ctx context.Context, key ...any) (T, error) {
	var zero T

	if err := r.checkKey(key); err != nil {
		return zero, err
	}

	where, _ := r.whereKey(0)
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		r.columnList(), quoteTable(r.mapping.Table), where)

	return r.ExecuteQuerySingle(ctx, query, key...)
}

func (r *Repository[T]) GetAll(ctx context.Context) ([]T, error) {
	query := fmt.Sprintf("SELECT %s FROM %s", r.columnList(), quoteTable(r.mapping.Table))

	return r.ExecuteQuery(ctx, query)
}

func (r *Repository[T]) Add(ctx context.Context, entity T) error {
	values, err := r.values(entity)
	if err != nil {
		return err
	}

	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = fmt.Sprintf("@p%d", i+1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		quoteTable(r.mapping.Table), r.columnList(), strings.Join(placeholders, ", "))

	_, err = r.ExecuteCommand(ctx, query, values...)
	return err
}

func (r *Repository[T]) Update(ctx context.Context, entity T) error {
	values, err := r.values(entity)
	if err != nil {
		return err
	}

	key := r.mapping.ID(entity)
	if err := r.checkKey(key); err != nil {
		return err
	}

	var (
		assignments []string
		args        []any
	)
	for i, column := range r.mapping.Columns {
		if r.isKeyColumn(column) {
			continue
		}
		args = append(args, values[i])
		assignments = append(assignments, fmt.Sprintf("%s = @p%d", quoteIdentifier(column), len(args)))
	}

	if len(assignments) == 0 {
		return fmt.Errorf("database: nothing to update, %q has no non-key columns", r.mapping.Table)
	}

	// The key placeholders continue after the SET arguments.
	where, _ := r.whereKey(len(args))
	args = append(args, key...)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		quoteTable(r.mapping.Table), strings.Join(assignments, ", "), where)

	affected, err := r.ExecuteCommand(ctx, query, args...)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}

	return nil
}

// Upsert inserts or updates in one statement, which keeps projection replays idempotent.
func (r *Repository[T]) Upsert(ctx context.Context, entity T) error {
	values, err := r.values(entity)
	if err != nil {
		return err
	}

	var (
		sourceColumns []string
		placeholders  []string
		assignments   []string
		insertValues  []string
		matches       []string
	)
	for i, column := range r.mapping.Columns {
		quoted := quoteIdentifier(column)
		sourceColumns = append(sourceColumns, quoted)
		placeholders = append(placeholders, fmt.Sprintf("@p%d", i+1))
		insertValues = append(insertValues, "source."+quoted)
		if r.isKeyColumn(column) {
			matches = append(matches, fmt.Sprintf("target.%s = source.%s", quoted, quoted))
		} else {
			assignments = append(assignments, fmt.Sprintf("target.%s = source.%s", quoted, quoted))
		}
	}

	if len(matches) != len(r.mapping.IDColumns) {
		return fmt.Errorf("database: key columns of %q are missing from Columns", r.mapping.Table)
	}
	if len(assignments) == 0 {
		return fmt.Errorf("database: cannot upsert %q, it has no non-key columns", r.mapping.Table)
	}

	// HOLDLOCK serializes the match check against concurrent inserts of the same key.
	query := fmt.Sprintf(`MERGE INTO %s WITH (HOLDLOCK) AS target
USING (VALUES (%s)) AS source (%s)
ON %s
WHEN MATCHED THEN UPDATE SET %s
WHEN NOT MATCHED THEN INSERT (%s) VALUES (%s);`,
		quoteTable(r.mapping.Table),
		strings.Join(placeholders, ", "),
		strings.Join(sourceColumns, ", "),
		strings.Join(matches, " AND "),
		strings.Join(assignments, ", "),
		strings.Join(sourceColumns, ", "),
		strings.Join(insertValues, ", "))

	_, err = r.ExecuteCommand(ctx, query, values...)
	return err
}

// Delete is idempotent: a row that is already gone is not an error.
func (r *Repository[T]) Delete(ctx context.Context, key ...any) error {
	if err := r.checkKey(key); err != nil {
		return err
	}

	where, _ := r.whereKey(0)
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", quoteTable(r.mapping.Table), where)

	_, err := r.ExecuteCommand(ctx, query, key...)
	return err
}

// whereKey renders the key predicate with placeholders numbered after startIndex.
func (r *Repository[T]) whereKey(startIndex int) (string, int) {
	predicates := make([]string, len(r.mapping.IDColumns))
	for i, column := range r.mapping.IDColumns {
		predicates[i] = fmt.Sprintf("%s = @p%d", quoteIdentifier(column), startIndex+i+1)
	}
	return strings.Join(predicates, " AND "), startIndex + len(r.mapping.IDColumns)
}

func (r *Repository[T]) checkKey(key []any) error {
	if len(key) != len(r.mapping.IDColumns) {
		return fmt.Errorf("database: %q needs %d key values, got %d",
			r.mapping.Table, len(r.mapping.IDColumns), len(key))
	}
	return nil
}

func (r *Repository[T]) isKeyColumn(column string) bool {
	for _, key := range r.mapping.IDColumns {
		if strings.EqualFold(column, key) {
			return true
		}
	}
	return false
}

func (r *Repository[T]) values(entity T) ([]any, error) {
	values := r.mapping.Values(entity)
	if len(values) != len(r.mapping.Columns) {
		return nil, fmt.Errorf("database: mapping for %q yields %d values for %d columns",
			r.mapping.Table, len(values), len(r.mapping.Columns))
	}
	return values, nil
}

func (r *Repository[T]) columnList() string {
	quoted := make([]string, len(r.mapping.Columns))
	for i, column := range r.mapping.Columns {
		quoted[i] = quoteIdentifier(column)
	}
	return strings.Join(quoted, ", ")
}

// quoteIdentifier bracket-quotes a T-SQL identifier; a literal ] is escaped by doubling it.
func quoteIdentifier(name string) string {
	return "[" + strings.ReplaceAll(name, "]", "]]") + "]"
}

func quoteTable(name string) string {
	parts := strings.Split(name, ".")
	for i, part := range parts {
		parts[i] = quoteIdentifier(part)
	}
	return strings.Join(parts, ".")
}
