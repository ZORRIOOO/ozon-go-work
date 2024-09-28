package stock

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	orderModel "homework/loms/internal/model/order"
	stockModel "homework/loms/internal/model/stock"
)

type (
	Repository struct{ connection *pgx.Conn }
)

func NewRepository(dbConn *pgx.Conn) *Repository {
	return &Repository{connection: dbConn}
}

const (
	selectStockBySKU = `
		SELECT total_count, reserved
		FROM stock
		WHERE sku = $1;
	`

	updateStockReserve = `
		UPDATE stock
		SET reserved = reserved + $1, total_count = total_count - $1
		WHERE sku = $2;
	`

	selectReservedBySKU = `
		SELECT reserved
		FROM stock
		WHERE sku = $1;
	`

	updateStockReserveRemove = `
		UPDATE stock
		SET reserved = reserved - $1, total_count = total_count - $1
		WHERE sku = $2;
	`

	updateStockReserveCancel = `
		UPDATE stock
		SET reserved = reserved - $1, total_count = total_count + $1
		WHERE sku = $2;
	`

	selectTotalCountBySKU = `
		SELECT total_count
		FROM stock
		WHERE sku = $1;
	`
)

func (r *Repository) Reserve(ctx context.Context, order orderModel.Order) error {
	for _, orderItem := range order.Items {
		var (
			sku        = orderItem.Sku
			quantity   = int(orderItem.Count)
			totalCount int
			reserved   int
			available  int
		)
		err := r.connection.QueryRow(ctx, selectStockBySKU, sku).Scan(&totalCount, &reserved)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("SKU %d not found in stock", sku)
			}
			return fmt.Errorf("select stock for SKU %d: %w", sku, err)
		}

		available = totalCount - reserved
		if available < quantity {
			return fmt.Errorf("not enough items to reserve for SKU %d", sku)
		}

		_, err = r.connection.Exec(ctx, updateStockReserve, quantity, sku)
		if err != nil {
			return fmt.Errorf("update stock reserve for SKU %d: %w", sku, err)
		}
	}

	return nil
}

func (r *Repository) ReserveRemove(ctx context.Context, order *orderModel.Order) error {
	for _, orderItem := range order.Items {
		sku := orderItem.Sku
		quantity := int(orderItem.Count)

		var reserved int
		err := r.connection.QueryRow(ctx, selectReservedBySKU, sku).Scan(&reserved)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("SKU %d not found in stock", sku)
			}
			return fmt.Errorf("select reserved stock for SKU %d: %w", sku, err)
		}

		if reserved < quantity {
			return fmt.Errorf("not enough reserved items for SKU %d", sku)
		}

		_, err = r.connection.Exec(ctx, updateStockReserveRemove, quantity, sku)
		if err != nil {
			return fmt.Errorf("update stock reserve removal for SKU %d: %w", sku, err)
		}
	}

	return nil
}

func (r *Repository) ReserveCancel(ctx context.Context, order *orderModel.Order) error {
	for _, orderItem := range order.Items {
		sku := orderItem.Sku
		quantity := int(orderItem.Count)

		var reserved int
		err := r.connection.QueryRow(ctx, selectReservedBySKU, sku).Scan(&reserved)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("SKU %d not found in stock", sku)
			}
			return fmt.Errorf("select reserved stock for SKU %d: %w", sku, err)
		}

		if reserved < quantity {
			return fmt.Errorf("not enough reserved items to cancel for SKU %d", sku)
		}

		_, err = r.connection.Exec(ctx, updateStockReserveCancel, quantity, sku)
		if err != nil {
			return fmt.Errorf("update stock reserve cancel for SKU %d: %w", sku, err)
		}
	}

	return nil
}

func (r *Repository) GetBySKU(ctx context.Context, sku stockModel.SKU) (stockModel.TotalCount, error) {
	var totalCount stockModel.TotalCount
	err := r.connection.QueryRow(ctx, selectTotalCountBySKU, sku).Scan(&totalCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("SKU %d not found in stock", sku)
		}
		return 0, fmt.Errorf("select total count for SKU %d: %w", sku, err)
	}

	return totalCount, nil
}
