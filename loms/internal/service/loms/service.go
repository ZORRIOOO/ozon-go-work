package loms

import (
	"homework/loms/internal/repository/order"
	"homework/loms/internal/repository/stock"
	"homework/loms/pkg/api/loms/v1"
)

var _ loms.LomsServer = (*Service)(nil)

type Service struct {
	orderRepository *order.Repository
	stockRepository *stock.Repository
	loms.UnimplementedLomsServer
}

func NewService(orderRepository *order.Repository, stockRepository *stock.Repository) *Service {
	return &Service{
		orderRepository: orderRepository,
		stockRepository: stockRepository,
	}
}
