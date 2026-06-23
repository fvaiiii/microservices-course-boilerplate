package part

import (
	"context"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/repository/converter"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/repository/record"
)

func (r *repository) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	if len(filter.UUIDs) > 0 {
		return r.listByUUIDs(ctx, filter.UUIDs)
	}

	return r.listByPartType(ctx, filter.PartType)
}

func (r *repository) listByUUIDs(ctx context.Context, uuids []string) ([]model.Part, error) {
	ids := make([]uuid.UUID, 0, len(uuids))

	for _, rawUUID := range uuids {
		id, err := uuid.Parse(rawUUID)
		if err != nil {
			return nil, errs.ErrInvalidUUID
		}
		ids = append(ids, id)
	}

	const query = `
		SELECT uuid, name, description, part_type, price, stock_quantity, created_at, updated_at
		FROM parts
		WHERE uuid = ANY($1)
	`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}
	defer rows.Close()

	recs, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.Part])
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}

	if len(recs) != len(ids) {
		return nil, errs.ErrPartNotFound
	}

	byUUID := make(map[string]record.Part, len(recs))
	for _, rec := range recs {
		byUUID[rec.UUID] = rec
	}

	parts := make([]model.Part, 0, len(uuids))
	for _, rawUUID := range uuids {
		rec, ok := byUUID[rawUUID]
		if !ok {
			return nil, errs.ErrPartNotFound
		}
		parts = append(parts, converter.PartToModel(rec))
	}

	return parts, nil
}

func (r *repository) listByPartType(ctx context.Context, partType model.PartType) ([]model.Part, error) {
	const baseQuery = `
		SELECT uuid, name, description, part_type, price, stock_quantity, created_at, updated_at
		FROM parts
	`

	var (
		rows pgx.Rows
		err  error
	)

	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	if partType == model.PartTypeUnspecified {
		rows, err = db.Query(ctx, baseQuery)
	} else {
		rows, err = db.Query(ctx, baseQuery+` WHERE part_type = $1`, string(partType))
	}
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}
	defer rows.Close()

	recs, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.Part])
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}

	parts := make([]model.Part, 0, len(recs))
	for _, rec := range recs {
		parts = append(parts, converter.PartToModel(rec))
	}

	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Name < parts[j].Name
	})

	return parts, nil
}
