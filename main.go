package main

import (
	"api-project/internal/db"
	"api-project/internal/todo"
	"api-project/internal/transport"
	"log"
)

func main() {
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
