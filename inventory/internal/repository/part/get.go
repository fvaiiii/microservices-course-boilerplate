package part

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	errs "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/repository/converter"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/repository/record"
)

func (r *repository) Get(ctx context.Context, rawUUID string) (model.Part, error) {
	partRec, err := r.getPart(ctx, rawUUID)
	if err != nil {
		return model.Part{}, err
	}
	return converter.PartToModel(partRec), nil
}

func (r *repository) getPart(ctx context.Context, id string) (record.Part, error) {
	const query = `
		SELECT uuid, name, description, part_type, price, stock_quantity, created_at, updated_at
		FROM parts
		WHERE uuid = $1
	`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, id)
	if err != nil {
		return record.Part{}, fmt.Errorf("получить заказ: %w", err)
	}
	defer rows.Close()

	partRec, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[record.Part])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return record.Part{}, errs.ErrPartNotFound
		}
		return record.Part{}, fmt.Errorf("получить заказ: %w", err)
	}

	return partRec, nil
}
