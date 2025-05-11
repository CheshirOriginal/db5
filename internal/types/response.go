package types

import "time"

type ProductInfoResponse struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type TellerInfoResponse struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

type DepartmentInfoResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type EmployeeInfoResponse struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Position   string `json:"position"`
	Salary     string `json:"salary"`
	Department string `json:"department"`
}

type SupplierInfoResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ProductInfoBySupplierResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type FullProductInfoResponse struct {
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	Quantity       int     `json:"quantity"`
	Category       string  `json:"category"`
	DepartmentName string  `json:"department_name"`
}

type FullReceiptInfoResponse struct {
	ID                int64
	TellerFirstName   string                   `json:"teller_first_name"`
	TellerLastName    string                   `json:"teller_last_name"`
	TellerMiddleName  string                   `json:"teller_middle_name"`
	Number            int                      `json:"number"`
	Date              time.Time                `json:"date"`
	Total             float64                  `json:"total"`
	LoyaltyCardNumber int                      `json:"loyalty_card_number"`
	Products          []ReceiptProductResponse `json:"products"`
}

type ReceiptProductResponse struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Amount   float64 `json:"amount"`
}

type FullSupplierOrderInfoResponse struct {
	ID                 int64                       `json:"id"`
	OrderDate          time.Time                   `json:"order_date"`
	DateOfReceipt      time.Time                   `json:"date_of_receipt"`
	Total              float64                     `json:"total"`
	SupplierName       string                      `json:"supplier_name"`
	SupplierOrderItems []SupplierOrderItemResponse `json:"supplier_order_items"`
}

type SupplierOrderItemResponse struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Amount      float64 `json:"amount"`
}
