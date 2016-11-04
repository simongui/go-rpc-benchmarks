package main

import "github.com/simongui/go-rpc-benchmarks/msgpack"

func main() {
	grpcService := &GRPCService{}
	grpcService.Listen("localhost:5555", "localhost:8001")
	go grpcService.Serve()
	defer grpcService.Close()

	service := &msgpack.Service{}
	service.Listen("localhost:50000", "localhost:8002")
	go service.Serve()
	defer service.Close()
	select {}
}
