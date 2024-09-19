package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework/loms/pkg/api/loms/v1"
)

var _ loms.LomsServer = (*Service)(nil)

type Service struct {
	loms.UnimplementedLomsServer
}

func NewService() *Service {
	return &Service{}
}

func (s Service) OrderCreate(ctx context.Context, request *loms.OrderCreateRequest) (*loms.OrderCreateResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) OrderInfo(ctx context.Context, request *loms.OrderInfoRequest) (*loms.OrderInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) OrderPay(ctx context.Context, request *loms.OrderPayRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) OrderCancel(ctx context.Context, request *loms.OrderCancelRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) StocksInfo(ctx context.Context, request *loms.StocksInfoRequest) (*loms.StocksInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) mustEmbedUnimplementedLomsServer() {
	//TODO implement me
	panic("implement me")
}
