package types

type ReceiptInfoRequest struct {
	LoyaltyCardNumber int64                       `json:"loyalty_card_number"`
	TellerID          int64                       `json:"teller_id"`
	Products          []ReceiptProductInfoRequest `json:"products"`
}

type ReceiptProductInfoRequest struct {
	ProductID int64   `json:"product_id"`
	Quantity  int64   `json:"quantity"`
	Price     float64 `json:"price"`
	Amount    float64 `json:"amount"`
}

type EmployeeInfoCreateRequest struct {
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	MiddleName   string  `json:"middle_name"`
	Position     string  `json:"position"`
	Salary       float64 `json:"salary"`
	DepartmentID int64   `json:"department_id"`
}

type EmployeeInfoDeleteRequest struct {
	ID int64 `json:"employee_id"`
}

type SupplierOrderInfoRequest struct {
	SupplierID         int64                          `json:"supplier_id"`
	SupplierOrderItems []SupplierOrderItemInfoRequest `json:"supplier_order_items"`
}

type SupplierOrderItemInfoRequest struct {
	Price     float64 `json:"price"`
	ProductID int64   `json:"product_id"`
	Quantity  int64   `json:"quantity"`
}
