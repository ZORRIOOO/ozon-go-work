package main

import (
	"google.golang.org/grpc"
	"homework/loms/internal/repository/order"
	"homework/loms/internal/repository/stock"
	loms "homework/loms/internal/service"
	desc "homework/loms/pkg/api/loms/v1"
	"net"
)

const (
	grpcPort = ":50051"
	capacity = 1000
	filePath = "./loms/assets/stock-data.json"
)

func main() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	orderRepository := order.NewRepository(capacity)
	stockRepository := stock.NewRepository(capacity, filePath)
	controller := loms.NewService(orderRepository, stockRepository)

	desc.RegisterLomsServer(grpcServer, controller)
	if err = grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
