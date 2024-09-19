package main

import (
	"google.golang.org/grpc"
	loms "homework/loms/internal/loms/service"
	desc "homework/loms/pkg/api/loms/v1"
	"net"
)

const grpcPort = ":50051"

func main() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	controller := loms.NewService()

	desc.RegisterLomsServer(grpcServer, controller)

	if err = grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
