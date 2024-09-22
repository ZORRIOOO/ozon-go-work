package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	orderModel "homework/loms/internal/model/order"
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

type OrderRepository interface {
	Create(order orderModel.Order) (orderModel.Id, error)
	GetById(id orderModel.Id) (orderModel.Order, error)
}

func NewService(orderRepository *order.Repository, stockRepository *stock.Repository) *Service {
	return &Service{
		orderRepository: orderRepository,
		stockRepository: stockRepository,
	}
}

func (s *Service) OrderCreate(ctx context.Context, request *loms.OrderCreateRequest) (*loms.OrderCreateResponse, error) {
	orderItem := RepackOrderCreate("new", request)
	orderId, createErr := s.orderRepository.Create(ctx, orderItem)
	if createErr != nil {
		return nil, status.Errorf(codes.Internal, createErr.Error())
	}

	orderItem.OrderId = orderId
	reserveErr := s.stockRepository.Reserve(orderItem)

	orderStatus := GetStatus(reserveErr)
	statusErr := s.orderRepository.SetStatus(orderId, orderStatus)
	if statusErr != nil {
		return nil, status.Errorf(codes.Internal, statusErr.Error())
	}

	return &loms.OrderCreateResponse{OrderId: orderId}, nil
}

func (s Service) OrderInfo(ctx context.Context, request *loms.OrderInfoRequest) (*loms.OrderInfoResponse, error) {
	orderId := request.GetOrderId()
	orderItem, err := s.orderRepository.GetById(ctx, orderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return RepackOrderToProto(orderItem), nil
}

func RepackOrderCreate(status orderModel.Status, in *loms.OrderCreateRequest) orderModel.Order {
	return orderModel.Order{
		Status: status,
		User:   in.User,
		Items:  in.Items,
	}
}

func RepackOrderToProto(orderItem *orderModel.Order) *loms.OrderInfoResponse {
	return &loms.OrderInfoResponse{
		Status: orderItem.Status,
		User:   orderItem.User,
		Items:  orderItem.Items,
	}
}

func GetStatus(err error) orderModel.Status {
	if err != nil {
		return "failed"
	} else {
		return "awaiting payment"
	}
}
