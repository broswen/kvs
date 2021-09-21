package main

import (
	"log"

	"github.com/broswen/kvs/internal/server"
)

func main() {

	server, err := server.New()
	if err != nil {
		log.Fatalf("Init server: %v\n", err)
	}

	server.SetRoutes()

	if err := server.Start(); err != nil {
		log.Fatalf("Start server: %v\n", err)
	}

}
