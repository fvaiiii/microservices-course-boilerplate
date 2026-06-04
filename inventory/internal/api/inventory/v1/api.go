package v1

import inventoryv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/inventory/v1"

type inventoryServer struct {
	inventoryv1.UnimplementedInventoryServiceServer
	partService PartService
}

func New(partService PartService) *inventoryServer {
	return &inventoryServer{
		partService: partService,
	}
}
