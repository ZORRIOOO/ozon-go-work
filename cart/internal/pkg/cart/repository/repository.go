package repository

import (
	"errors"
	"homework/cart/internal/pkg/cart/model"
	"sync"
)

type (
	CartStorage = map[model.UserId]map[int64]model.CartItem

	CartRepository struct {
		storage CartStorage
		mx      sync.Mutex
	}
)

func NewCartRepository(capacity int) *CartRepository {
	return &CartRepository{storage: make(CartStorage, capacity)}
}

func (r *CartRepository) AddItem(item model.CartItem) (*model.CartItem, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

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

func (r *CartRepository) DeleteItem(params model.DeleteCartParameters) (*model.CartItem, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

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
	r.mx.Lock()
	defer r.mx.Unlock()

	if len(r.storage[userId]) == 0 {
		return &userId, nil
	}

	delete(r.storage, userId)

	return &userId, nil
}

func (r *CartRepository) GetItemsByUser(userId model.UserId) ([]model.CartItem, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	if r.storage[userId] == nil {
		return nil, errors.New("message=There is no such a cart, status=404")
	}

	storageItems := r.storage[userId]
	items := make([]model.CartItem, 0, len(storageItems))
	for _, item := range storageItems {
		items = append(items, item)
	}

	return items, nil
}
