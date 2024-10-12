package loms

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework/loms/internal/infra/kafka/utils"
	orderModel "homework/loms/internal/model/order"
	"homework/loms/pkg/api/loms/v1"
)

func (s *Service) OrderCreate(ctx context.Context, request *loms.OrderCreateRequest) (*loms.OrderCreateResponse, error) {
	orderStatus := "new"
	orderItem := RepackOrderCreate(orderStatus, request)
	orderId, createErr := s.orderRepository.Create(ctx, orderItem)
	if createErr != nil {
		return nil, status.Errorf(codes.Internal, createErr.Error())
	}

	payload := utils.RepackPayload(orderId, orderStatus)
	err := s.kafkaEmitter.SendMessage(payload) // используем emitter как зависимость
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	orderItem.OrderId = orderId
	reserveErr := s.stockRepository.Reserve(ctx, orderItem)
	orderStatus = GetStatus(reserveErr)

	payload = utils.RepackPayload(orderId, orderStatus)
	err = s.kafkaEmitter.SendMessage(payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &loms.OrderCreateResponse{OrderId: orderId}, nil
}

func RepackOrderCreate(status orderModel.Status, in *loms.OrderCreateRequest) orderModel.Order {
	items := make([]orderModel.Item, len(in.Items))
	for i, item := range in.Items {
		items[i] = orderModel.Item{
			Sku:   item.Sku,
			Count: item.Count,
		}
	}
	return orderModel.Order{
		Status: status,
		User:   in.User,
		Items:  items,
	}
}

func GetStatus(err error) orderModel.Status {
	if err != nil {
		return "failed"
	} else {
		return "awaiting payment"
	}
}
