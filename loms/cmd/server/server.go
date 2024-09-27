package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"homework/loms/core/reader"
	"homework/loms/core/utils"
	"homework/loms/internal/mw"
	"homework/loms/internal/repository/order"
	"homework/loms/internal/repository/stock"
	"homework/loms/internal/service/loms"
	desc "homework/loms/pkg/api/loms/v1"
	"log"
	"net"
)

const (
	grpcPort = ":50051"
	capacity = 1000
	filePath = "./loms/assets/stock-data.json"
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
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("loms", grpc_health_v1.HealthCheckResponse_SERVING)

	stocks, err := reader.ReadStocks(utils.GetEnv("DOCKER_PATH_ASSETS", filePath))
	if err != nil {
		fmt.Sprintf("Read stocks failed: %v", err.Error())
	}

	orderRepository := order.NewRepository(capacity)
	stockRepository := stock.NewRepository(capacity, stocks)
	controller := loms.NewService(orderRepository, stockRepository)

	desc.RegisterLomsServer(grpcServer, controller)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error server: %s", err.Error())
	}
}
