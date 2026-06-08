package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/fvaiiii/microservices-course-boilerplate/order/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	inventoryv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/inventory/v1"
)

type client struct {
	api inventoryv1.InventoryServiceClient
}

func New(api inventoryv1.InventoryServiceClient) *client {
	return &client{
		api: api,
	}
}

func (c *client) ListParts(ctx context.Context, uuids []uuid.UUID) ([]model.Part, error) {
	ids := make([]string, 0, len(uuids))
	for _, id := range uuids {
		ids = append(ids, id.String())
	}
	resp, err := c.api.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids: ids,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, fmt.Errorf("получить список деталей: %w", errs.ErrPartNotFound)
		}
		return nil, fmt.Errorf("list parts grpc: %w", err)
	}
	return converter.ToParts(resp.Parts)
}
