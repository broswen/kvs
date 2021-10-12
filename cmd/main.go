package main

import (
	"flag"
	"log"

	"github.com/broswen/kvs/internal/server"
)

var grpcFlag bool

func main() {

	flag.BoolVar(&grpcFlag, "grpc", false, "set to enable gRPC")
	flag.Parse()

	if grpcFlag {
		server, err := server.NewGRPC()
		if err != nil {
			log.Fatalf("Init server: %v\n", err)
		}

		if err := server.Start(); err != nil {
			log.Fatalf("Start server: %v\n", err)
		}
	} else {
		server, err := server.New()
		if err != nil {
			log.Fatalf("Init server: %v\n", err)
		}

		if err := server.Start(); err != nil {
			log.Fatalf("Start server: %v\n", err)
		}

	}

}
