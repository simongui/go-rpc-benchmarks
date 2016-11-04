package main

import (
	"log"
	"net"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	contract "github.com/simongui/go-rpc-benchmarks/grpc"
	"github.com/valyala/fasthttp"
)

// GRPCService Represents a GRPC service instance.
type GRPCService struct {
	grpcListener   net.Listener
	grpcServer     *grpc.Server
	grpcConnection *grpc.ClientConn
	grpcClient     contract.GreeterClient
	httpListener   net.Listener
}

// Listen Listens on the specified GRPC and HTTP addresses.
func (service *GRPCService) Listen(grpcAddress string, httpAddress string) {
	var err error
	service.grpcListener, err = net.Listen("tcp", grpcAddress)
	// listener, err := reuseport.Listen("tcp4", grpcAddress)
	if err != nil {
		log.Fatal(err)
	}

	service.grpcServer = grpc.NewServer()
	contract.RegisterGreeterServer(service.grpcServer, service)

	service.httpListener, err = net.Listen("tcp", httpAddress)
	if err != nil {
		log.Fatal(err)
	}

	service.grpcConnection, err = grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	service.grpcClient = contract.NewGreeterClient(service.grpcConnection)
}

// Serve Serves starts serving requests from clients.
func (service *GRPCService) Serve() error {
	go fasthttp.Serve(service.httpListener, service.grpcRequestHandler)
	return service.grpcServer.Serve(service.grpcListener)
}

// Close Closes all client connections to the GRPC service.
func (service *GRPCService) Close() error {
	var err error

	err = service.grpcConnection.Close()
	if err != nil {
		return err
	}

	err = service.httpListener.Close()
	if err != nil {
		return err
	}
	err = service.grpcListener.Close()
	if err != nil {
		return err
	}
	return nil
}

// SayHello implements grpc.GreeterService.
func (service *GRPCService) SayHello(ctx context.Context, in *contract.HelloRequest) (*contract.HelloResponse, error) {
	return &contract.HelloResponse{Message: "Hello " + in.Name}, nil
}

func (service *GRPCService) grpcRequestHandler(ctx *fasthttp.RequestCtx) {
	_, err := service.grpcClient.SayHello(context.Background(), &contract.HelloRequest{Name: "world"})
	if err != nil {
		ctx.SetStatusCode(500)
	}
}
