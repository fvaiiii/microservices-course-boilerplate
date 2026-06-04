package part

type service struct {
	partRepo PartRepository
}

func New(partRepo PartRepository) *service {
	return &service{
		partRepo: partRepo,
	}
}
