package tests

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"homework/loms/internal/service/loms"
	"homework/loms/internal/service/loms/mocks"
	lomsApi "homework/loms/pkg/api/loms/v1"
	"testing"
)

func Test_StockInfo(t *testing.T) {
	ctx := context.Background()
	ctrl := minimock.NewController(t)
	orderRepositoryMock := mocks.NewOrderRepositoryMock(ctrl)
	stocksRepositoryMock := mocks.NewStocksRepositoryMock(ctrl)
	lomsService := loms.NewService(orderRepositoryMock, stocksRepositoryMock)

	sku := int64(773297411)
	count := int32(10)
	request := &lomsApi.StocksInfoRequest{Sku: sku}

	stocksRepositoryMock.GetBySKUMock.Expect(ctx, sku).Return(count, nil)

	actualResponse, err := lomsService.StocksInfo(ctx, request)
	expectedResponse := &lomsApi.StocksInfoResponse{Count: count}

	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)
}
