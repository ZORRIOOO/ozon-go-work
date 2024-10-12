package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"homework/loms/internal/infra/kafka/emitter"
	"homework/loms/internal/mw"
	"homework/loms/internal/repository/order"
	"homework/loms/internal/repository/stock"
	"homework/loms/internal/service/loms"
	desc "homework/loms/pkg/api/loms/v1"
	"log"
	"net"
	"os"
)

const (
	grpcPort   = ":50051"
	brokerAddr = "broker:29092"
	connection = "postgres://user:password@database:5432/homework"
)

func main() {
	log.Println("Go loms service starting")

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		mw.Validate,
		mw.Logging,
		mw.Panic,
	))
	reflection.Register(grpcServer)

	dbConn, err := pgx.Connect(context.Background(), connection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err.Error())
		os.Exit(1)
	}
	defer dbConn.Close(context.Background())

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("loms", grpc_health_v1.HealthCheckResponse_SERVING)

	var (
		kafkaEmitter    = emitter.NewEmitter(brokerAddr)
		orderRepository = order.NewRepository(dbConn)
		stockRepository = stock.NewRepository(dbConn)
		controller      = loms.NewService(orderRepository, stockRepository, kafkaEmitter)
	)

	desc.RegisterLomsServer(grpcServer, controller)
	log.Println("Go loms service ready")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error server: %s", err.Error())
	}
}
