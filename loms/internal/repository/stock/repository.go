package stock

import (
	"context"
	"fmt"
	orderModel "homework/loms/internal/model/order"
	stockModel "homework/loms/internal/model/stock"
	"sync"
)

type (
	Storage = map[stockModel.SKU]stockModel.Stock

	Repository struct {
		storage Storage
		mx      sync.Mutex
	}
)

func NewRepository(capacity int, stocks []stockModel.Stock) *Repository {
	storage := make(Storage, capacity)
	for _, stock := range stocks {
		storage[stock.SKU] = stock
	}
	return &Repository{
		storage: storage,
	}
}

func (r *Repository) Reserve(_ context.Context, order orderModel.Order) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	for _, orderItem := range order.Items {
		sku := orderItem.Sku
		quantity := orderItem.Count

		stockItem, itemsFound := r.storage[sku]
		if !itemsFound {
			return fmt.Errorf("SKU %d not found in stock", sku)
		}

		available := stockItem.TotalCount - stockItem.Reserved
		if available < quantity {
			return fmt.Errorf("not enough items to reserve for SKU: %d", sku)
		}

		stockItem.TotalCount -= quantity
		stockItem.Reserved += quantity
		r.storage[sku] = stockItem
	}

	return nil
}

func (r *Repository) ReserveRemove(_ context.Context, order *orderModel.Order) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	for _, orderItem := range order.Items {
		sku := orderItem.Sku
		quantity := orderItem.Count

		stockItem, itemsFound := r.storage[sku]
		if !itemsFound {
			return fmt.Errorf("SKU %d not found in stock", sku)
		}

		if stockItem.Reserved < quantity {
			return fmt.Errorf("not enough reserved items for SKU: %d", sku)
		}

		stockItem.Reserved -= quantity
		r.storage[sku] = stockItem
	}

	return nil
}

func (r *Repository) ReserveCancel(_ context.Context, order *orderModel.Order) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	for _, orderItem := range order.Items {
		sku := orderItem.Sku
		quantity := orderItem.Count

		stockItem, itemsFound := r.storage[sku]
		if !itemsFound {
			return fmt.Errorf("SKU %d not found in stock", sku)
		}

		if stockItem.Reserved < quantity {
			return fmt.Errorf("not enough reserved items to cancel for SKU: %d", sku)
		}

		stockItem.Reserved -= quantity
		stockItem.TotalCount += quantity
		r.storage[sku] = stockItem
	}

	return nil
}

func (r *Repository) GetBySKU(_ context.Context, sku stockModel.SKU) (stockModel.TotalCount, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	stockItem, itemsFound := r.storage[sku]
	if !itemsFound {
		return 0, fmt.Errorf("SKU %d not found in stock", sku)
	}

	return stockItem.TotalCount, nil
}
