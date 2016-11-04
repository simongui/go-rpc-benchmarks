package main

import (
	"log"
	"net"
	"testing"

	"github.com/msgpack-rpc/msgpack-rpc-go/rpc"
	"github.com/simongui/go-rpc-benchmarks/msgpack"
)

func BenchmarkMSGPack(b *testing.B) {
	var conn net.Conn
	var err error

	service := &msgpack.Service{}
	service.Listen("localhost:50000", "localhost:8002")
	go service.Serve()
	defer service.Close()

	conn, err = net.Dial("tcp", "127.0.0.1:50000")
	if err != nil {
		log.Fatal(err)
	}
	client := rpc.NewSession(conn, true)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = client.Send("echo", "world")
		if err != nil {
			log.Fatal(err)
		}
	}
}
