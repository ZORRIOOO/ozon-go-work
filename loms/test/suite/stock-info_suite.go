package suite

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"homework/loms/internal/infra/kafka/emitter"
	"homework/loms/internal/repository/order"
	"homework/loms/internal/repository/stock"
	lomsService "homework/loms/internal/service/loms"
	"homework/loms/pkg/api/loms/v1"
	"os"
)

type StockInfoSuite struct {
	suite.Suite
	service *lomsService.Service
}

func (s *StockInfoSuite) SetupSuite() {
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
		kafkaEmitter    = emitter.NewEmitter(brokerAddr)
		orderRepository = order.NewRepository(dbConn)
		stockRepository = stock.NewRepository(dbConn)
		controller      = lomsService.NewService(orderRepository, stockRepository, kafkaEmitter)
	)
	s.service = controller
}

func (s *StockInfoSuite) TestStockInfo() {
	ctx := context.Background()
	sku := int64(773297411)

	stockInfoRequest := &loms.StocksInfoRequest{Sku: sku}
	stockInfoResponse, err := s.service.StocksInfo(ctx, stockInfoRequest)

	expectedResponse := &loms.StocksInfoResponse{Count: stockInfoResponse.Count}
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedResponse, stockInfoResponse)
}
