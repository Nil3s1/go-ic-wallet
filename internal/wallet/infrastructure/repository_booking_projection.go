package infrastructure

import (
	"context"
	"fmt"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel/infrastructure/database"
	"github.com/Nil3s1/go-ic-wallet/internal/wallet/application"
)

var bookingProjectionMapping = database.Mapping[application.BookingProjection]{
	Table:     "wallet.tbl_Booking",
	IDColumns: []string{"CardNo", "ReferenceId"},
	Columns:   []string{"CardNo", "ReferenceId", "Amount", "Direction", "BookedAt", "BalanceAfter"},
	Scan: func(s database.Scanner) (application.BookingProjection, error) {
		var (
			booking      application.BookingProjection
			amount       int64
			balanceAfter int64
		)
		if err := s.Scan(
			&booking.CardNo,
			&booking.ReferenceId,
			&amount,
			&booking.Direction,
			&booking.BookedAt,
			&balanceAfter,
		); err != nil {
			return application.BookingProjection{}, err
		}
		booking.Amount = uint(amount)
		booking.BalanceAfter = uint(balanceAfter)
		return booking, nil
	},
	// The driver rejects uint, so the amounts travel as int64.
	Values: func(b application.BookingProjection) []any {
		return []any{b.CardNo, b.ReferenceId, int64(b.Amount), b.Direction, b.BookedAt, int64(b.BalanceAfter)}
	},
	ID: func(b application.BookingProjection) []any { return []any{b.CardNo, b.ReferenceId} },
}

var _ application.BookingProjectionRepository = (*BookingProjectionRepository)(nil)

type BookingProjectionRepository struct {
	repo *database.Repository[application.BookingProjection]
}

func NewBookingProjectionRepository(exec database.Executor) *BookingProjectionRepository {
	return &BookingProjectionRepository{repo: database.NewRepository(exec, bookingProjectionMapping)}
}

func (r *BookingProjectionRepository) WithTx(exec database.Executor) *BookingProjectionRepository {
	return &BookingProjectionRepository{repo: r.repo.WithTx(exec)}
}

func (r *BookingProjectionRepository) GetByReference(ctx context.Context, cardNo string, referenceId string) (application.BookingProjection, error) {
	return r.repo.GetByKey(ctx, cardNo, referenceId)
}

// ListByCard orders by ReferenceId as a tiebreaker so paging stays stable on equal timestamps.
func (r *BookingProjectionRepository) ListByCard(ctx context.Context, cardNo string, limit int, offset int) ([]application.BookingProjection, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("wallet: limit must be greater than 0")
	}
	if offset < 0 {
		return nil, fmt.Errorf("wallet: offset must not be negative")
	}

	const query = `SELECT [CardNo], [ReferenceId], [Amount], [Direction], [BookedAt], [BalanceAfter]
FROM [wallet].[tbl_Booking]
WHERE [CardNo] = @p1
ORDER BY [BookedAt] DESC, [ReferenceId]
OFFSET @p2 ROWS FETCH NEXT @p3 ROWS ONLY`

	return r.repo.ExecuteQuery(ctx, query, cardNo, offset, limit)
}

// Upsert keeps event replays idempotent; the application interface only reads.
func (r *BookingProjectionRepository) Upsert(ctx context.Context, booking application.BookingProjection) error {
	return r.repo.Upsert(ctx, booking)
}
