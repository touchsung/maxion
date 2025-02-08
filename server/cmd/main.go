package main

import (
	"log"

	"github.com/touchsung/maxion-server/internal/config"
	"github.com/touchsung/maxion-server/internal/server"
)

func main() {
	db, err := config.GetDatabaseConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	srv := server.NewServer(db)
	log.Fatal(srv.Start(":3000"))
}
