package loms

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework/loms/pkg/api/loms/v1"
)

func (s Service) StocksInfo(ctx context.Context, request *loms.StocksInfoRequest) (*loms.StocksInfoResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	sku := request.GetSku()
	count, err := s.stockRepository.GetBySKU(ctx, sku)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &loms.StocksInfoResponse{Count: count}, err
}
