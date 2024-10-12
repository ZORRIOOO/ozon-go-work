package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"homework/loms/internal/infra/kafka/provider/producer"
	"homework/loms/internal/mw"
	"homework/loms/internal/repository/order"
	"homework/loms/internal/repository/stock"
	"homework/loms/internal/service/loms"
	desc "homework/loms/pkg/api/loms/v1"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	grpcPort   = ":50051"
	brokerAddr = "http://kafka-hw:9092"
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
		kafkaProducer, configuration = producer.NewKafkaProducer(brokerAddr)
		orderRepository              = order.NewRepository(dbConn)
		stockRepository              = stock.NewRepository(dbConn)
		controller                   = loms.NewService(orderRepository, stockRepository)
	)

	bytes, err := json.Marshal("EVENT!")
	msg := &sarama.ProducerMessage{
		Topic: configuration.Producer.Topic,
		Key:   sarama.StringEncoder(strconv.FormatInt(1, 10)),
		Value: sarama.ByteEncoder(bytes),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("app-name"),
				Value: []byte("route256-sync-prod"),
			},
		},
		Partition: 1,
		Timestamp: time.Now(),
	}

	partition, offset, err := kafkaProducer.SendMessage(msg)

	fmt.Println(partition, offset, err)

	desc.RegisterLomsServer(grpcServer, controller)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error server: %s", err.Error())
	}
}
