package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
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
	filePath = "./assets/stock-data.json"
)

func main() {
	log.Println("Go loms service starting")

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("loms", grpc_health_v1.HealthCheckResponse_SERVING)

	orderRepository := order.NewRepository(capacity)
	stockRepository := stock.NewRepository(capacity, filePath)
	controller := loms.NewService(orderRepository, stockRepository)

	desc.RegisterLomsServer(grpcServer, controller)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error server: %s", err.Error())
	}
}
