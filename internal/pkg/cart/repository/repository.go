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
