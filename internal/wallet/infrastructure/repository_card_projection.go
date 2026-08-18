package infrastructure

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel/infrastructure/database"
	"github.com/Nil3s1/go-ic-wallet/internal/wallet/application"
)

var cardProjectionMapping = database.Mapping[application.CardProjection]{
	Table:     "wallet.tbl_Card",
	IDColumns: []string{"CardNo"},
	Columns:   []string{"CardNo", "ValidTo", "CurrentBalance"},
	Scan: func(s database.Scanner) (application.CardProjection, error) {
		var (
			card    application.CardProjection
			balance int64
		)
		if err := s.Scan(&card.CardNo, &card.ValidTo, &balance); err != nil {
			return application.CardProjection{}, err
		}
		card.CurrentBalance = uint(balance)
		return card, nil
	},
	// The driver rejects uint, so the balance travels as int64.
	Values: func(c application.CardProjection) []any {
		return []any{c.CardNo, c.ValidTo, int64(c.CurrentBalance)}
	},
	ID: func(c application.CardProjection) []any { return []any{c.CardNo} },
}

var _ application.CardProjectionRepository = (*CardProjectionRepository)(nil)

type CardProjectionRepository struct {
	repo *database.Repository[application.CardProjection]
}

func NewCardProjectionRepository(exec database.Executor) *CardProjectionRepository {
	return &CardProjectionRepository{repo: database.NewRepository(exec, cardProjectionMapping)}
}

func (r *CardProjectionRepository) WithTx(exec database.Executor) *CardProjectionRepository {
	return &CardProjectionRepository{repo: r.repo.WithTx(exec)}
}

func (r *CardProjectionRepository) GetByKey(ctx context.Context, cardNo string) (application.CardProjection, error) {
	return r.repo.GetByKey(ctx, cardNo)
}

func (r *CardProjectionRepository) Add(ctx context.Context, card application.CardProjection) error {
	return r.repo.Add(ctx, card)
}

func (r *CardProjectionRepository) Update(ctx context.Context, card application.CardProjection) error {
	return r.repo.Update(ctx, card)
}

// Upsert keeps event replays idempotent; the application interface only knows Add and Update.
func (r *CardProjectionRepository) Upsert(ctx context.Context, card application.CardProjection) error {
	return r.repo.Upsert(ctx, card)
}
