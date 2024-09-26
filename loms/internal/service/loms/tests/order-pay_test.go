package tests

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	orderModel "homework/loms/internal/model/order"
	"homework/loms/internal/service/loms"
	"homework/loms/internal/service/loms/mocks"
	lomsApi "homework/loms/pkg/api/loms/v1"
	"testing"
)

func Test_OrderPay(t *testing.T) {
	ctx := context.Background()
	ctrl := minimock.NewController(t)
	orderRepositoryMock := mocks.NewOrderRepositoryMock(ctrl)
	stocksRepositoryMock := mocks.NewStocksRepositoryMock(ctrl)
	lomsService := loms.NewService(orderRepositoryMock, stocksRepositoryMock)

	orderId := int64(1)
	user := int64(1)
	sku := int64(773297411)
	status := "awaiting payment"
	statusPayed := "payed"
	count := int32(10)
	request := &lomsApi.OrderPayRequest{
		OrderId: orderId,
	}
	items := []*lomsApi.Item{{
		Sku:   sku,
		Count: count,
	}}
	order := &orderModel.Order{
		OrderId: orderId,
		Status:  status,
		User:    user,
		Items:   items,
	}

	orderRepositoryMock.GetByIdMock.Expect(ctx, orderId).Return(order, nil)
	stocksRepositoryMock.ReserveRemoveMock.Expect(ctx, order).Return(nil)
	orderRepositoryMock.SetStatusMock.Expect(ctx, orderId, statusPayed).Return(nil)

	actualResponse, err := lomsService.OrderPay(ctx, request)
	expectedResponse := &emptypb.Empty{}

	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)
}
