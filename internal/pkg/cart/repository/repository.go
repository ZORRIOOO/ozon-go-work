package repository

import (
	"cart/internal/pkg/cart/model"
)

type CartStorage = map[model.UserId]map[int64]model.CartItem

type CartRepository struct {
	storage CartStorage
}

func NewCartRepository(capacity int) *CartRepository {
	return &CartRepository{storage: make(CartStorage, capacity)}
}

func (r *CartRepository) AddItem(item model.CartItem) (*model.CartItem, error) {
	if r.storage[item.UserId] == nil {
		r.storage[item.UserId] = make(map[int64]model.CartItem)
	}

	if existingItem, exists := r.storage[item.UserId][item.SKU]; exists {
		existingItem.Count += item.Count
		r.storage[item.UserId][item.SKU] = existingItem

		return &existingItem, nil
	} else {
		r.storage[item.UserId][item.SKU] = item

		return &item, nil
	}
}

func (r *CartRepository) DeleteItem(_ model.CartItem) (*model.CartItem, error) {
	return nil, nil
}
