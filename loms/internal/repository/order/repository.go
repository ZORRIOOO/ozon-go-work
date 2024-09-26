package order

import (
	"context"
	"errors"
	model "homework/loms/internal/model/order"
	"sync"
)

type (
	Storage = map[model.Id]model.Order

	Repository struct {
		storage   Storage
		increment int64
		mx        sync.Mutex
	}
)

func NewRepository(capacity int) *Repository {
	return &Repository{
		storage:   make(Storage, capacity),
		increment: 0,
	}
}

func (r *Repository) Create(_ context.Context, order model.Order) (model.Id, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.increment++
	order.OrderId = r.increment

	r.storage[order.OrderId] = order
	return order.OrderId, nil
}

func (r *Repository) SetStatus(_ context.Context, id model.Id, status model.Status) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if order, exists := r.storage[id]; exists {
		order.Status = status
		r.storage[id] = order
		return nil
	}

	return errors.New("order not found")
}

func (r *Repository) GetById(_ context.Context, id model.Id) (*model.Order, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	if order, exists := r.storage[id]; exists {
		return &order, nil
	}

	return nil, errors.New("order not found")
}
