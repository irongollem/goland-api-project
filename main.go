package main

import (
	"github.com/irongollem/goland-api-project/internal/db"
	"github.com/irongollem/goland-api-project/internal/todo"
	"github.com/irongollem/goland-api-project/internal/transport"
	"log"
)

func main() {
	// fixme: Replace with credentials from ENV
	d, err := db.NewDB("postgres", "example", "postgres", "localhost", 5432)
	if err != nil {
		log.Fatal(err)
	}

	service := todo.NewService(d)
	server := transport.NewServer(service)

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
