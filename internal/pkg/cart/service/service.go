package service

type CartRepository interface{}

type CartService struct {
	repository CartRepository
}

func NewCartService(repository CartRepository) *CartService {
	return &CartService{repository: repository}
}
