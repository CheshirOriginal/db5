package main

import (
	"db5/config"
	"db5/internal/db"
	"fmt"
	"log"
)

func main() {
	conf := config.LoadConfig()

	var Database db.DB

	err := Database.Connect(conf)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	defer Database.Close()

	Cashiers, err := Database.GetAllReceipt()
	if err != nil {
		log.Fatalf("failed to get products: %v", err)
	}

	fmt.Println(Cashiers)
}
