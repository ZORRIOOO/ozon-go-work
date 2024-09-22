package loms

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework/loms/internal/repository/order"
	"homework/loms/internal/repository/stock"
	"homework/loms/pkg/api/loms/v1"
)

var _ loms.LomsServer = (*Service)(nil)

type OrderService interface {
	OrderPay(ctx context.Context, request *loms.OrderPayRequest) (emptypb.Empty, error)
}

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
