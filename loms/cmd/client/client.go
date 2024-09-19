package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	desc "homework/loms/pkg/api/loms/v1"
)

func main() {
	conn, err := grpc.Dial(":50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := desc.NewLomsClient(conn)
	ctx := context.Background()

	response, err := client.OrderCreate(ctx, &desc.OrderCreateRequest{
		User:  123,
		Items: make([]*desc.Item, 0),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(response)
}
