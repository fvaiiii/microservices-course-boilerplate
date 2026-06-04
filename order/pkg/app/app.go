package app

import (
	"net/http"

	orderapi "github.com/fvaiiii/microservices-course-boilerplate/order/internal/api/order/v1"
	inventoryclient "github.com/fvaiiii/microservices-course-boilerplate/order/internal/client/grpc/inventory/v1"
	paymentclient "github.com/fvaiiii/microservices-course-boilerplate/order/internal/client/grpc/payment/v1"
	orderrepo "github.com/fvaiiii/microservices-course-boilerplate/order/internal/repository/order"
	ordersvc "github.com/fvaiiii/microservices-course-boilerplate/order/internal/service/order"
	orderv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/payment/v1"
)

func NewHTTPHandler(
	inventoryClient inventoryv1.InventoryServiceClient,
	paymentClient paymentv1.PaymentServiceClient,
) (http.Handler, error) {
	inventoryRepo := inventoryclient.New(inventoryClient)
	paymentRepo := paymentclient.New(paymentClient)

	orderRepo := orderrepo.New()
	orderService := ordersvc.New(orderRepo, inventoryRepo, paymentRepo)
	api := orderapi.New(orderService)

	server, err := orderv1.NewServer(api)
	if err != nil {
		return nil, err
	}

	return server, nil
}
