package loms

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	orderModel "homework/loms/internal/model/order"
	"homework/loms/pkg/api/loms/v1"
)

func (s Service) OrderInfo(ctx context.Context, request *loms.OrderInfoRequest) (*loms.OrderInfoResponse, error) {
	orderId := request.GetOrderId()
	orderItem, err := s.orderRepository.GetById(ctx, orderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return RepackOrderToProto(orderItem), nil
}

func RepackOrderToProto(orderItem *orderModel.Order) *loms.OrderInfoResponse {
	return &loms.OrderInfoResponse{
		Status: orderItem.Status,
		User:   orderItem.User,
		Items:  orderItem.Items,
	}
}
