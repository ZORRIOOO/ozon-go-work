package loms

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework/loms/core/utils"
	kafkaUtils "homework/loms/internal/infra/kafka/utils"
	"homework/loms/pkg/api/loms/v1"
)

func (s Service) OrderCancel(ctx context.Context, request *loms.OrderCancelRequest) (*emptypb.Empty, error) {
	orderId := request.GetOrderId()
	orderItem, err := s.orderRepository.GetById(ctx, orderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	if utils.Contains([]string{"payed", "cancelled"}, orderItem.Status) {
		return &emptypb.Empty{}, status.Error(codes.FailedPrecondition, fmt.Sprintf("order is unavailable"))
	}

	err = s.stockRepository.ReserveCancel(ctx, orderItem)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, err.Error())
	}

	orderStatus := "cancelled"
	err = s.orderRepository.SetStatus(ctx, orderId, orderStatus)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, err.Error())
	}

	payload := kafkaUtils.RepackPayload(orderId, orderStatus)
	err = s.kafkaEmitter.SendMessage(payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
