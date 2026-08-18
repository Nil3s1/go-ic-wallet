package infrastructure

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/journey/application"
	"github.com/Nil3s1/go-ic-wallet/internal/kernel/infrastructure/database"
)

var journeyLogProjectionMapping = database.Mapping[application.JourneyLogProjection]{
	Table:     "journey.tbl_JourneyLog",
	IDColumns: []string{"MediaId"},
	Columns:   []string{"MediaId", "LastJourneyReferenceId", "IsOnJourney", "LastStation", "UpdatedAt"},
	Scan: func(s database.Scanner) (application.JourneyLogProjection, error) {
		var log application.JourneyLogProjection
		if err := s.Scan(
			&log.MediaId,
			&log.LastJourneyReferenceId,
			&log.IsOnJourney,
			&log.LastStation,
			&log.UpdatedAt,
		); err != nil {
			return application.JourneyLogProjection{}, err
		}
		return log, nil
	},
	Values: func(l application.JourneyLogProjection) []any {
		return []any{l.MediaId, l.LastJourneyReferenceId, l.IsOnJourney, l.LastStation, l.UpdatedAt}
	},
	ID: func(l application.JourneyLogProjection) []any { return []any{l.MediaId} },
}

var _ application.JourneyLogProjectionRepository = (*JourneyLogProjectionRepository)(nil)

type JourneyLogProjectionRepository struct {
	repo *database.Repository[application.JourneyLogProjection]
}

func NewJourneyLogProjectionRepository(exec database.Executor) *JourneyLogProjectionRepository {
	return &JourneyLogProjectionRepository{repo: database.NewRepository(exec, journeyLogProjectionMapping)}
}

func (r *JourneyLogProjectionRepository) WithTx(exec database.Executor) *JourneyLogProjectionRepository {
	return &JourneyLogProjectionRepository{repo: r.repo.WithTx(exec)}
}

func (r *JourneyLogProjectionRepository) GetByKey(ctx context.Context, mediaId string) (application.JourneyLogProjection, error) {
	return r.repo.GetByKey(ctx, mediaId)
}

// Upsert keeps event replays idempotent; the application interface only reads.
func (r *JourneyLogProjectionRepository) Upsert(ctx context.Context, log application.JourneyLogProjection) error {
	return r.repo.Upsert(ctx, log)
}
