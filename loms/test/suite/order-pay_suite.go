package suite

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework/loms/internal/infra/kafka/broker/producer"
	"homework/loms/internal/repository/order"
	"homework/loms/internal/repository/stock"
	lomsService "homework/loms/internal/service/loms"
	"homework/loms/pkg/api/loms/v1"
	"os"
)

type OrderPaySuite struct {
	suite.Suite
	service *lomsService.Service
}

func (s *OrderPaySuite) SetupSuite() {
	const (
		connection = "postgres://user:password@localhost:5432/homework"
		brokerAddr = "localhost:9092"
	)

	dbConn, err := pgx.Connect(context.Background(), connection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	var (
		kafkaProducer   = producer.NewKafkaProducer(brokerAddr)
		orderRepository = order.NewRepository(dbConn)
		stockRepository = stock.NewRepository(dbConn)
		controller      = lomsService.NewService(orderRepository, stockRepository, kafkaProducer)
	)
	s.service = controller
}

func (s *OrderPaySuite) TestOrderPay() {
	ctx := context.Background()
	sku := int64(773297411)
	count := int32(1)
	user := int64(1)
	items := make([]*loms.Item, 0)
	status := "payed"
	items = append(items, &loms.Item{
		Sku:   sku,
		Count: count,
	})

	orderCreateRequest := &loms.OrderCreateRequest{User: user, Items: items}
	orderCreateResponse, err := s.service.OrderCreate(ctx, orderCreateRequest)

	require.NoError(s.T(), err)

	orderPayRequest := &loms.OrderPayRequest{OrderId: orderCreateResponse.OrderId}
	orderPayResponse, err := s.service.OrderPay(ctx, orderPayRequest)

	require.NoError(s.T(), err)
	require.Equal(s.T(), &emptypb.Empty{}, orderPayResponse)

	orderInfoRequest := &loms.OrderInfoRequest{OrderId: orderCreateResponse.OrderId}
	orderInfoResponse, err := s.service.OrderInfo(ctx, orderInfoRequest)

	expectedResponse := &loms.OrderInfoResponse{Status: status, User: user, Items: items}
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedResponse, orderInfoResponse)
}
