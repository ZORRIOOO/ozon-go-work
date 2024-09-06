package repository

import (
	"errors"
	"fmt"
	"homework/cart/internal/pkg/cart/model"
	"net/http"
)

type CartStorage = map[model.UserId]map[int64]model.CartItem

type CartRepository struct {
	storage CartStorage
}

func NewCartRepository(capacity int) *CartRepository {
	return &CartRepository{storage: make(CartStorage, capacity)}
}

func (r *CartRepository) AddItem(item model.CartItem) (*model.CartItem, error, int) {
	if r.storage[item.UserId] == nil {
		r.storage[item.UserId] = make(map[int64]model.CartItem)
	}

	if existingItem, exists := r.storage[item.UserId][item.SKU]; exists {
		existingItem.Count += item.Count
		r.storage[item.UserId][item.SKU] = existingItem

		return &existingItem, nil, http.StatusOK
	} else {
		r.storage[item.UserId][item.SKU] = item

		return &item, nil, http.StatusOK
	}
}

func (r *CartRepository) DeleteItem(params model.DeleteCartParameters) (*model.CartItem, error) {
	if r.storage[params.UserId] == nil {
		return &model.CartItem{SKU: params.SKU, UserId: params.UserId}, nil
	}

	item, exists := r.storage[params.UserId][params.SKU]
	if !exists {
		return &model.CartItem{SKU: params.SKU, UserId: params.UserId}, nil
	}

	delete(r.storage[params.UserId], params.SKU)

	return &item, nil
}

func (r *CartRepository) DeleteItemsByUser(userId model.UserId) (*model.UserId, error) {
	if r.storage[userId] == nil || len(r.storage[userId]) == 0 {
		return &userId, nil
	}

	delete(r.storage, userId)

	return &userId, nil
}

func (r *CartRepository) GetItemsByUser(userId model.UserId) ([]model.CartItem, error, int) {
	if r.storage[userId] == nil {
		message := fmt.Sprintf("There is no such a cart")
		return []model.CartItem{}, errors.New(message), http.StatusNotFound
	}

	storageItems := r.storage[userId]
	items := make([]model.CartItem, 0, len(storageItems))
	for _, item := range storageItems {
		items = append(items, item)
	}

	return items, nil, http.StatusOK
}
