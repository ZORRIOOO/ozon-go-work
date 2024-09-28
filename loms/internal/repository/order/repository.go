package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	model "homework/loms/internal/model/order"
	"sync"
)

type (
	Storage = map[model.Id]model.Order

	Repository struct {
		conn      *pgx.Conn
		storage   Storage
		increment int64
		mx        sync.Mutex
	}
)

func NewRepository(conn *pgx.Conn) *Repository {
	return &Repository{
		conn:      conn,
		storage:   make(Storage, 0),
		increment: 0,
	}
}

const (
	insertOrder = `
		INSERT INTO "order" (status, "user")
		VALUES ($1, $2)
		RETURNING id;
	`

	insertItem = `
		INSERT INTO item (sku, "count", order_id)
		VALUES ($1, $2, $3);
	`

	updateOrderStatus = `
		UPDATE "order"
		SET status = $1
		WHERE id = $2;
	`

	selectOrderById = `
		SELECT id, status, "user"
		FROM "order"
		WHERE id = $1;
	`

	selectItemsByOrderId = `
		SELECT sku, "count"
		FROM item
		WHERE order_id = $1;
	`
)

func (r *Repository) Create(ctx context.Context, order model.Order) (orderId model.Id, err error) {
	err = r.conn.QueryRow(ctx, insertOrder, order.Status, order.User).Scan(&orderId)
	if err != nil {
		return 0, fmt.Errorf("insert order: %w", err)
	}
	for _, item := range order.Items {
		if _, err = r.conn.Exec(ctx, insertItem, item.Sku, item.Count, orderId); err != nil {
			return 0, fmt.Errorf("insert item: %w", err)
		}
	}
	return orderId, nil
}

func (r *Repository) SetStatus(ctx context.Context, id model.Id, status model.Status) error {
	result, err := r.conn.Exec(ctx, updateOrderStatus, status, id)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("order not found")
	}

	return nil
}

func (r *Repository) GetById(ctx context.Context, id model.Id) (*model.Order, error) {
	var (
		orderId model.Id
		status  model.Status
		user    model.User
		items   []model.Item
	)
	err := r.conn.QueryRow(ctx, selectOrderById, id).Scan(&orderId, &status, &user)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("select order: %w", err)
	}

	rows, err := r.conn.Query(ctx, selectItemsByOrderId, id)
	if err != nil {
		return nil, fmt.Errorf("select items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.Sku, &item.Count); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &model.Order{
		OrderId: orderId,
		Status:  status,
		User:    user,
		Items:   items,
	}, nil
}
