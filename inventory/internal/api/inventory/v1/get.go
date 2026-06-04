package v1

import (
	"context"

	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/converter"
	inventoryv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/inventory/v1"
)

func (s *inventoryServer) GetPart(
	ctx context.Context,
	req *inventoryv1.GetPartRequest,
) (*inventoryv1.GetPartResponse, error) {
	part, err := s.partService.Get(ctx, req.GetUuid())
	if err != nil {
		return nil, err
	}

	return &inventoryv1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
