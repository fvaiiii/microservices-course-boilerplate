package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	inventoryapi "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/api/inventory/v1"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/interceptor"
	partrepo "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/repository/part"
	partsvc "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/service/part"
	inventoryv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/inventory/v1"
)

func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
	}
}

func RegisterServices(grpcServer *grpc.Server, pool *pgxpool.Pool, txManager partrepo.TxManager) {
	repo := partrepo.New(pool, txManager)
	svc := partsvc.New(repo)
	api := inventoryapi.New(svc)
	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)
}
