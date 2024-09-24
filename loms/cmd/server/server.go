package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"homework/loms/internal/mw"
	"homework/loms/internal/repository/order"
	"homework/loms/internal/repository/stock"
	"homework/loms/internal/service/loms"
	desc "homework/loms/pkg/api/loms/v1"
	"log"
	"net"
	"net/http"
)

const (
	grpcPort = ":50051"
	httpPort = ":8081"
	capacity = 1000
	filePath = "loms/assets/stock-data.json"
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

	orderRepository := order.NewRepository(capacity)
	stockRepository := stock.NewRepository(capacity, filePath)
	controller := loms.NewService(orderRepository, stockRepository)

	desc.RegisterLomsServer(grpcServer, controller)
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("Error server: %s", err.Error())
		}
	}()

	conn, err := grpc.NewClient(grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error connecting to server: %s", err.Error())
	}
	gwmux := runtime.NewServeMux()
	if err = desc.RegisterLomsHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err.Error())
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})
	gwServer := &http.Server{
		Addr:    httpPort,
		Handler: mw.HTTPLogging(c.Handler(gwmux)),
	}
	log.Printf("Serving gRPC-Gateway on PORT: %s\n", gwServer.Addr)
	log.Fatalln(gwServer.ListenAndServe())
}
