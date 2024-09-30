package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework/loms/internal/mw"
	desc "homework/loms/pkg/api/loms/v1"
	"log"
	"net/http"
)

const (
	grpcPort = ":50051"
	httpPort = ":8081"
)

func main() {
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
