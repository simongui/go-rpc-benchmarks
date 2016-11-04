package main

import (
	"log"
	"testing"

	contract "github.com/simongui/go-rpc-benchmarks/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func BenchmarkGRPC(b *testing.B) {
	grpcService := &GRPCService{}
	grpcService.Listen("localhost:5555", "localhost:8001")
	go grpcService.Serve()
	defer grpcService.Close()

	connection, err := grpc.Dial("localhost:5555", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	client := contract.NewGreeterClient(connection)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client.SayHello(context.Background(), &contract.HelloRequest{Name: "world"})
	}
}
