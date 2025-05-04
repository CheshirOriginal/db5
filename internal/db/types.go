package db

import "time"

type Product struct {
	ID       int64
	Name     string
	Price    float64
	Category string
	Quantity int
}

type ProductInfo struct {
	ID   int64
	Name string
}

type Employee struct {
	ID         int64
	FirstName  string
	LastName   string
	MiddleName string
	Position   string
	Salary     float64
}

type SupplierInfo struct {
	ID   int64
	Name string
}

type ReceiptProduct struct {
	Name     string
	Quantity int
	Amount   float64
	Price    float64
}

type Receipt struct {
	ID              int64
	Number          int
	Date            time.Time
	Total           float64
	ReceiptProducts []ReceiptProduct
}
