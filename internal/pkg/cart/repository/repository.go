package repository

import (
	"cart/internal/pkg/cart/model"
	"context"
)

type CartStorage = map[model.UserId][]model.CartItem

type CartRepository struct {
	storage CartStorage
}

func NewCartRepository(capacity int) *CartRepository {
	return &CartRepository{storage: make(CartStorage, capacity)}
}

func (r *CartRepository) AddItem(_ context.Context, item model.CartItem) (*model.CartItem, error) {
	return &item, nil
}
