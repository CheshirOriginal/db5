package types

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
