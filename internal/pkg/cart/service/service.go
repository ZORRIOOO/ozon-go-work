package service

import (
	"cart/internal/pkg/cart/model"
	"context"
	"fmt"
)

type CartRepository interface {
	AddItem(context context.Context, item model.CartItem) (*model.CartItem, error)
}

type CartService struct {
	repository CartRepository
}

func (cartService CartService) DeleteItem(ctx context.Context, cartItem model.CartItem) (*model.CartItem, error) {
	fmt.Println(ctx, cartItem)
	return nil, nil
}

func (cartService CartService) AddItem(ctx context.Context, cartItem model.CartItem) (*model.CartItem, error) {
	fmt.Println(ctx, cartItem)
	return nil, nil
}

func NewCartService(repository CartRepository) *CartService {
	return &CartService{repository: repository}
}
