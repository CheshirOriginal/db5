package main

import (
	"db5/config"
	"db5/internal/db"
	"db5/internal/server"
	"log"
)

func main() {
	conf := config.LoadConfig()

	var Database db.DB

	err := Database.Connect(conf)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	mux := server.CreateNewServerMux(&Database)

	s := server.CreateNewServer(*mux)

	err = s.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
