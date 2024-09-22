package order

import (
	"context"
	"errors"
	model "homework/loms/internal/model/order"
	"sync"
)

type OrdersStorage = map[model.User]map[model.Id]model.Order

type Repository struct {
	storage   OrdersStorage
	increment int64
	mx        sync.Mutex
}

func NewRepository(capacity int) *Repository {
	return &Repository{
		storage:   make(OrdersStorage, capacity),
		increment: 0,
	}
}

func (r *Repository) Create(_ context.Context, order model.Order) (model.Id, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.increment++
	order.OrderId = r.increment
	if _, exists := r.storage[order.User]; !exists {
		r.storage[order.User] = make(map[model.Id]model.Order)
	}

	r.storage[order.User][order.OrderId] = order
	return order.OrderId, nil
}

func (r *Repository) SetStatus(_ context.Context, id model.Id, status model.Status) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	for user, orders := range r.storage {
		if order, exists := orders[id]; exists {
			order.Status = status
			r.storage[user][id] = order
			return nil
		}
	}

	return errors.New("order not found")
}

func (r *Repository) GetById(_ context.Context, id model.Id) (*model.Order, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	for _, orders := range r.storage {
		if order, exists := orders[id]; exists {
			return &order, nil
		}
	}

	return nil, errors.New("order not found")
}
