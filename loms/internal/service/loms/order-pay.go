package loms

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework/loms/pkg/api/loms/v1"
)

func (s Service) OrderPay(ctx context.Context, request *loms.OrderPayRequest) (*emptypb.Empty, error) {
	err := request.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	orderId := request.GetOrderId()
	orderItem, err := s.orderRepository.GetById(ctx, orderId)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.NotFound, err.Error())
	}

	err = s.stockRepository.ReserveRemove(ctx, orderItem)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, err.Error())
	}

	orderStatus := "payed"
	err = s.orderRepository.SetStatus(ctx, orderId, orderStatus)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
