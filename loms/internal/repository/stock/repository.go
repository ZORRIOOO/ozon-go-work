package stock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	orderModel "homework/loms/internal/model/order"
	model "homework/loms/internal/model/stock"
	"io"
	"log"
	"os"
	"sync"
)

type Repository struct {
	storage []model.Stock
	mx      sync.Mutex
}

func NewRepository(capacity int, filePath string) *Repository {
	repository := &Repository{
		storage: make([]model.Stock, 0, capacity),
	}

	err := repository.InitStock(filePath)
	if err != nil {
		log.Fatalf("Error loading data from JSON: %v", err.Error())
	}

	return repository
}

func (r *Repository) Reserve(_ context.Context, order orderModel.Order) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	for _, orderItem := range order.Items {
		sku := orderItem.Sku
		quantity := orderItem.Count
		itemsFound := false

		for i, stockItem := range r.storage {
			if stockItem.SKU == sku {
				itemsFound = true
				available := stockItem.TotalCount - stockItem.Reserved

				if available < quantity {
					return fmt.Errorf("not enough items to reserve for SKU: %d", sku)
				}

				r.storage[i].TotalCount -= quantity
				r.storage[i].Reserved += quantity
				break
			}
		}

		if !itemsFound {
			return fmt.Errorf("SKU %d not found in stock", sku)
		}
	}

	return nil
}

func (r *Repository) ReserveRemove(_ context.Context, order *orderModel.Order) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	for _, orderItem := range order.Items {
		sku := orderItem.Sku
		quantity := orderItem.Count
		itemsFound := false

		for i, stockItem := range r.storage {
			if stockItem.SKU == sku {
				itemsFound = true

				if stockItem.Reserved < quantity {
					return fmt.Errorf("not enough reserved items for SKU: %d", sku)
				}

				r.storage[i].Reserved -= quantity
				break
			}
		}

		if !itemsFound {
			return fmt.Errorf("SKU %d not found in stock", sku)
		}
	}

	return nil
}

func (r *Repository) ReserveCancel(_ context.Context, order *orderModel.Order) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	for _, orderItem := range order.Items {
		sku := orderItem.Sku
		quantity := orderItem.Count
		itemsFound := false

		for i, stockItem := range r.storage {
			if stockItem.SKU == sku {
				itemsFound = true

				if stockItem.Reserved < quantity {
					return fmt.Errorf("not enough reserved items to cancel for SKU: %d", sku)
				}

				r.storage[i].Reserved -= quantity
				r.storage[i].TotalCount += quantity
				break
			}
		}

		if !itemsFound {
			return fmt.Errorf("SKU %d not found in stock", sku)
		}
	}

	return nil
}

func (r *Repository) InitStock(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("File open error: %v", err))
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return errors.New(fmt.Sprintf("File read error: %v", err))
	}

	var stocks []model.Stock
	err = json.Unmarshal(byteValue, &stocks)
	if err != nil {
		return errors.New(fmt.Sprintf("JSON parsing error: %v", err))
	}

	r.storage = stocks
	return nil
}
