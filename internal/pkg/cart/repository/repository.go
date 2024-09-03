package repository

import (
	"cart/internal/pkg/cart/model"
)

type CartStorage = map[model.UserId][]model.CartItem

type CartRepository struct {
	storage CartStorage
}

func NewCartRepository(capacity int) *CartRepository {
	return &CartRepository{storage: make(CartStorage, capacity)}
}

func (r *CartRepository) AddItem(item model.CartItem) (*model.CartItem, error) {
	return nil, nil
}

func (r *CartRepository) DeleteItem(_ model.CartItem) (*model.CartItem, error) {
	return nil, nil
}
