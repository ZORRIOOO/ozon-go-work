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

func (r *CartRepository) DeleteItem(skuId int64, userId model.UserId) (*model.CartItem, error) {
	if r.storage[userId] == nil {
		return &model.CartItem{SKU: skuId, UserId: userId, Count: 0}, nil
	}

	item, exists := r.storage[userId][skuId]
	if !exists {
		return &model.CartItem{SKU: skuId, UserId: userId, Count: 0}, nil
	}

	delete(r.storage[userId], skuId)

	return &item, nil
}

func (r *CartRepository) DeleteItemsByUser(userId model.UserId) (*model.UserId, error) {
	if r.storage[userId] == nil || len(r.storage[userId]) == 0 {
		return &userId, nil
	}

	delete(r.storage, userId)

	return &userId, nil
}

func (r *CartRepository) GetItemsByUser(userId model.UserId) ([]model.CartItem, error) {
	if r.storage[userId] == nil {
		return []model.CartItem{}, nil
	}

	storageItems := r.storage[userId]
	items := make([]model.CartItem, 0, len(storageItems))
	for _, item := range storageItems {
		items = append(items, item)
	}

	return items, nil
}
