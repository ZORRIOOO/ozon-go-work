package loms

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	orderModel "homework/loms/internal/model/order"
	"homework/loms/pkg/api/loms/v1"
)

func (s *Service) OrderCreate(ctx context.Context, request *loms.OrderCreateRequest) (*loms.OrderCreateResponse, error) {
	orderItem := RepackOrderCreate("new", request)
	orderId, createErr := s.orderRepository.Create(ctx, orderItem)
	if createErr != nil {
		return nil, status.Errorf(codes.Internal, createErr.Error())
	}

	orderItem.OrderId = orderId
	reserveErr := s.stockRepository.Reserve(ctx, orderItem)

	orderStatus := GetStatus(reserveErr)
	statusErr := s.orderRepository.SetStatus(ctx, orderId, orderStatus)
	if statusErr != nil {
		return nil, status.Errorf(codes.Internal, statusErr.Error())
	}

	return &loms.OrderCreateResponse{OrderId: orderId}, nil
}

func RepackOrderCreate(status orderModel.Status, in *loms.OrderCreateRequest) orderModel.Order {
	return orderModel.Order{
		Status: status,
		User:   in.User,
		Items:  in.Items,
	}
}

func GetStatus(err error) orderModel.Status {
	if err != nil {
		return "failed"
	} else {
		return "awaiting payment"
	}
}
