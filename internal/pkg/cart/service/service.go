package service

type CartRepository interface{}

type CartService struct {
	repository CartRepository
}

func NewService(repository CartService) *CartService {
	return &CartService{repository: repository}
}
