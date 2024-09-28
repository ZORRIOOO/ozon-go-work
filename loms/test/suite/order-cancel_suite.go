package suite

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework/loms/internal/repository/order"
	"homework/loms/internal/repository/stock"
	lomsService "homework/loms/internal/service/loms"
	"homework/loms/pkg/api/loms/v1"
	"os"
)

type OrderCancelSuite struct {
	suite.Suite
	service *lomsService.Service
}

func (s *OrderCancelSuite) SetupSuite() {
	const connection = "postgres://user:password@localhost:5432/homework"

	dbConn, err := pgx.Connect(context.Background(), connection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	var (
		orderRepository = order.NewRepository(dbConn)
		stockRepository = stock.NewRepository(dbConn)
		controller      = lomsService.NewService(orderRepository, stockRepository)
	)
	s.service = controller
}

func (s *OrderCancelSuite) TestOrderCancel() {
	ctx := context.Background()
	sku := int64(773297411)
	count := int32(1)
	user := int64(1)
	items := make([]*loms.Item, 0)
	status := "cancelled"
	items = append(items, &loms.Item{
		Sku:   sku,
		Count: count,
	})

	orderCreateRequest := &loms.OrderCreateRequest{User: user, Items: items}
	orderCreateResponse, err := s.service.OrderCreate(ctx, orderCreateRequest)

	require.NoError(s.T(), err)

	orderCancelRequest := &loms.OrderCancelRequest{OrderId: orderCreateResponse.OrderId}
	orderCancelResponse, err := s.service.OrderCancel(ctx, orderCancelRequest)

	require.NoError(s.T(), err)
	require.Equal(s.T(), &emptypb.Empty{}, orderCancelResponse)

	orderInfoRequest := &loms.OrderInfoRequest{OrderId: orderCreateResponse.OrderId}
	orderInfoResponse, err := s.service.OrderInfo(ctx, orderInfoRequest)

	expectedResponse := &loms.OrderInfoResponse{Status: status, User: user, Items: items}
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedResponse, orderInfoResponse)
}
