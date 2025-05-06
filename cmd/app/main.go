package main

import (
	"db5/config"
	"db5/internal/db"
	"db5/internal/types"
	"log"
)

func main() {
	conf := config.LoadConfig()

	var store db.Store
	var Database db.DB
	store = &Database

	err := store.Connect(conf)
	if err != nil {

		log.Fatalf("failed to connect to DB: %v", err)
	}

	defer store.Close()

	var receipt types.ReceiptInfoRequest
	receipt.TellerID = 1
	receipt.Products = append(receipt.Products, types.ReceiptProductInfo{ProductID: 1, Quantity: 2, Price: 700, Amount: 1400})
	err = store.CreateNewReceipt(receipt)
	if err != nil {
		log.Fatalf("failed to create new receipt: %v", err)
	}

	//result, err := store.GetTellerInfo()
	//if err != nil {
	//	log.Fatalf("failed to get product info: %v", err)
	//}
	//fmt.Println(result)

}
