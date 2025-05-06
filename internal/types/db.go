package types

type Employee struct {
	ID         int64
	FirstName  string
	LastName   string
	MiddleName string
	Position   string
	Salary     float64
}

func (e *Employee) ToTellerInfoRequest() TellerInfoResponse {
	return TellerInfoResponse{
		ID:         e.ID,
		FirstName:  e.FirstName,
		LastName:   e.LastName,
		MiddleName: e.MiddleName,
	}
}

//type Product struct {
//	ID       int64
//	Name     string
//	Price    float64
//	Category string
//	Quantity int
//}
//type SupplierInfo struct {
//	ID   int64
//	Name string
//}
//
//type ReceiptProduct struct {
//	Name      string
//	ProductID int64
//	Quantity  int
//	Amount    float64
//	Price     float64
//}
//
//type Receipt struct {
//	ID              int64
//	EmployeeID      int64
//	Date            time.Time
//	Total           float64
//	ReceiptProducts []ReceiptProduct
//}
