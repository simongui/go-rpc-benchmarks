package msgpack

import (
	"fmt"
	"log"
	"net"
	"reflect"

	"github.com/msgpack-rpc/msgpack-rpc-go/rpc"
	"github.com/valyala/fasthttp"
)

// Service Represents an instance of the MSGPack RPC service.
type Service struct {
	msgpackListener net.Listener
	msgpackServer   *rpc.Server
	msgpackClient   *rpc.Session
	httpListener    net.Listener
}

// Listen Starts the TCP listener.
func (service *Service) Listen(msgpackAddress string, httpAddress string) error {
	var err error

	res := Resolver{
		"echo": reflect.ValueOf(echo),
	}

	service.msgpackServer = rpc.NewServer(res, true, nil)
	service.msgpackListener, err = net.Listen("tcp", msgpackAddress)
	if err != nil {
		log.Fatal(err)
	}
	service.msgpackServer.Listen(service.msgpackListener)

	service.httpListener, err = net.Listen("tcp", httpAddress)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.Dial("tcp", msgpackAddress)
	if err != nil {
		log.Fatal(err)
	}
	service.msgpackClient = rpc.NewSession(conn, true)
	return nil
}

// Serve Starts serving requests.
func (service *Service) Serve() error {
	go fasthttp.Serve(service.httpListener, service.msgpackRequestHandler)
	service.msgpackServer.Run()
	return nil
}

// Close Closes all client connections to the MSGPack service.
func (service *Service) Close() error {
	var err error

	err = service.httpListener.Close()
	if err != nil {
		return err
	}
	service.msgpackServer.Stop()
	err = service.msgpackListener.Close()
	if err != nil {
		return err
	}
	return nil
}

func echo(name string) (string, fmt.Stringer) {
	return "Hello, " + name, nil
}

func (service *Service) msgpackRequestHandler(ctx *fasthttp.RequestCtx) {
	// var retval reflect.Value
	var err error

	_, err = service.msgpackClient.Send("echo", "world")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		ctx.SetStatusCode(500)
	}
}
