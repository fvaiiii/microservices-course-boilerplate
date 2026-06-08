package v1

import (
	"context"

	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/converter"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/model"
	inventoryv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/inventory/v1"
)

func (s *inventoryServer) ListParts(
	ctx context.Context,
	req *inventoryv1.ListPartsRequest,
) (*inventoryv1.ListPartsResponse, error) {
	parts, err := s.partService.List(ctx, model.PartFilter{
		UUIDs:    req.GetUuids(),
		PartType: converter.PartTypeToModel(req.GetPartType()),
	})
	if err != nil {
		return nil, err
	}

	partsProto := make([]*inventoryv1.Part, 0, len(parts))
	for _, part := range parts {
		partsProto = append(partsProto, converter.PartToProto(part))
	}

	return &inventoryv1.ListPartsResponse{
		Parts: partsProto,
	}, nil
}
