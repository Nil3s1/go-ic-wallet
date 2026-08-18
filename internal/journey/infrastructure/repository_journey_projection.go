package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Nil3s1/go-ic-wallet/internal/journey/application"
	"github.com/Nil3s1/go-ic-wallet/internal/kernel/infrastructure/database"
)

var journeyProjectionMapping = database.Mapping[application.JourneyProjection]{
	Table:     "journey.tbl_Journey",
	IDColumns: []string{"JourneyReferenceId"},
	Columns: []string{
		"JourneyReferenceId", "MediaId", "StartStation", "StartTime",
		"EndStation", "EndTime", "DistanceTravelled", "Fare", "Status",
	},
	Scan: func(s database.Scanner) (application.JourneyProjection, error) {
		var (
			journey  application.JourneyProjection
			endTime  sql.NullTime
			distance int64
			fare     int64
		)
		if err := s.Scan(
			&journey.JourneyReferenceId,
			&journey.MediaId,
			&journey.StartStation,
			&journey.StartTime,
			&journey.EndStation,
			&endTime,
			&distance,
			&fare,
			&journey.Status,
		); err != nil {
			return application.JourneyProjection{}, err
		}
		if endTime.Valid {
			journey.EndTime = &endTime.Time
		}
		journey.DistanceTravelled = uint(distance)
		journey.Fare = uint(fare)
		return journey, nil
	},
	// The driver rejects uint, so distance and fare travel as int64.
	Values: func(j application.JourneyProjection) []any {
		endTime := sql.NullTime{}
		if j.EndTime != nil {
			endTime = sql.NullTime{Time: *j.EndTime, Valid: true}
		}
		return []any{
			j.JourneyReferenceId, j.MediaId, j.StartStation, j.StartTime,
			j.EndStation, endTime, int64(j.DistanceTravelled), int64(j.Fare), j.Status,
		}
	},
	ID: func(j application.JourneyProjection) []any { return []any{j.JourneyReferenceId} },
}

var _ application.JourneyProjectionRepository = (*JourneyProjectionRepository)(nil)

type JourneyProjectionRepository struct {
	repo *database.Repository[application.JourneyProjection]
}

func NewJourneyProjectionRepository(exec database.Executor) *JourneyProjectionRepository {
	return &JourneyProjectionRepository{repo: database.NewRepository(exec, journeyProjectionMapping)}
}

func (r *JourneyProjectionRepository) WithTx(exec database.Executor) *JourneyProjectionRepository {
	return &JourneyProjectionRepository{repo: r.repo.WithTx(exec)}
}

func (r *JourneyProjectionRepository) GetByKey(ctx context.Context, journeyReferenceId string) (application.JourneyProjection, error) {
	return r.repo.GetByKey(ctx, journeyReferenceId)
}

// ListByMedia orders by JourneyReferenceId as a tiebreaker so paging stays stable on equal timestamps.
func (r *JourneyProjectionRepository) ListByMedia(ctx context.Context, mediaId string, limit int, offset int) ([]application.JourneyProjection, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("journey: limit must be greater than 0")
	}
	if offset < 0 {
		return nil, fmt.Errorf("journey: offset must not be negative")
	}

	const query = `SELECT [JourneyReferenceId], [MediaId], [StartStation], [StartTime],
       [EndStation], [EndTime], [DistanceTravelled], [Fare], [Status]
FROM [journey].[tbl_Journey]
WHERE [MediaId] = @p1
ORDER BY [StartTime] DESC, [JourneyReferenceId]
OFFSET @p2 ROWS FETCH NEXT @p3 ROWS ONLY`

	return r.repo.ExecuteQuery(ctx, query, mediaId, offset, limit)
}

// Upsert keeps event replays idempotent; the application interface only reads.
func (r *JourneyProjectionRepository) Upsert(ctx context.Context, journey application.JourneyProjection) error {
	return r.repo.Upsert(ctx, journey)
}
