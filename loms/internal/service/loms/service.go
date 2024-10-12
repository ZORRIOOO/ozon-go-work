package loms

import (
	"context"
	"homework/loms/internal/infra/kafka/types"
	orderModel "homework/loms/internal/model/order"
	stockModel "homework/loms/internal/model/stock"
	"homework/loms/pkg/api/loms/v1"
)

var _ loms.LomsServer = (*Service)(nil)

type (
	OrderRepository interface {
		Create(context context.Context, order orderModel.Order) (orderModel.Id, error)
		SetStatus(context context.Context, id orderModel.Id, status orderModel.Status) error
		GetById(context context.Context, id orderModel.Id) (*orderModel.Order, error)
	}

	StocksRepository interface {
		Reserve(context context.Context, order orderModel.Order) error
		ReserveRemove(context context.Context, order *orderModel.Order) error
		ReserveCancel(context context.Context, order *orderModel.Order) error
		GetBySKU(context context.Context, sku stockModel.SKU) (stockModel.TotalCount, error)
	}

	KafkaEmitter interface {
		SendMessage(payload types.MessagePayload) error
	}

	Service struct {
		orderRepository OrderRepository
		stockRepository StocksRepository
		kafkaEmitter    KafkaEmitter
		loms.UnimplementedLomsServer
	}
)

func NewService(orderRepository OrderRepository, stockRepository StocksRepository, kafkaEmitter KafkaEmitter) *Service {
	return &Service{
		orderRepository: orderRepository,
		stockRepository: stockRepository,
		kafkaEmitter:    kafkaEmitter,
	}
}
