package order

type service struct {
	orderRepo     OrderRepository
	inventoryRepo InventoryClient
	paymentRepo   PaymentClient
}

func New(
	orderRepo OrderRepository,
	inventoryRepo InventoryClient,
	paymentRepo PaymentClient,
) *service {
	return &service{
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
		paymentRepo:   paymentRepo,
	}
}
