package v1

type api struct {
	orderService OrderService
}

func New(orderService OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
