package payment

type paymentService struct{}

func New() *paymentService {
	return &paymentService{}
}
