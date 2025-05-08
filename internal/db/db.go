package db

import (
	"database/sql"
	"db5/config"
	"db5/internal/types"
	"fmt"
	_ "sync"

	_ "github.com/lib/pq"
)

// Store подумать над интерфейсом для настроек
type Store interface {
	Connect(c config.Config) error
	Close()
	GetProductInfo() ([]types.ProductInfoResponse, error)
	GetTellerInfo() ([]types.TellerInfoResponse, error)
	CreateNewReceipt(receiptInfo types.ReceiptInfoRequest) error
	GetDepartmentInfo() ([]types.DepartmentInfoResponse, error)
	CreateNewEmployee(employeeInfo types.EmployeeInfoCreateRequest) error
	GetEmployeeInfo() ([]types.EmployeeInfoResponse, error)
	DeleteEmployee(employeeInfo types.EmployeeInfoDeleteRequest) error
}

type DB struct {
	db *sql.DB
}

func (db *DB) Connect(c config.Config) error {
	database, err := sql.Open("postgres", c.GetDSN())
	if err != nil {
		return fmt.Errorf("sql open error: %w", err)
	}

	pingErr := database.Ping()
	if pingErr != nil {
		return fmt.Errorf("ping error: %w", err)
	}

	db.db = database
	fmt.Println("Connected!")

	return nil
}

func (db *DB) Close() {
	if db.db != nil {
		_ = db.db.Close()
	}
}

// GetProductInfo возможно стоит сделать как в GetTellerInfo
func (db *DB) GetProductInfo() ([]types.ProductInfoResponse, error) {
	var products []types.ProductInfoResponse

	rows, err := db.db.Query("select id, name, price, quantity_in_stock  from Product")
	if err != nil {
		return nil, fmt.Errorf("GetProductInfo: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product types.ProductInfoResponse
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity); err != nil {
			return nil, fmt.Errorf("GetProductInfo: %v", err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetProductInfo: %v", err)
	}
	return products, nil
}

func (db *DB) GetTellerInfo() ([]types.TellerInfoResponse, error) {
	result, err := db.getEmployeeByPosition("Кассир")
	if err != nil {
		return nil, fmt.Errorf("GetTellerInfo: %v", err)
	}

	tellers := make([]types.TellerInfoResponse, len(result))

	for i := range result {
		tellers[i] = result[i].TellerInfoResponse()
	}

	return tellers, nil
}

// CreateNewReceipt добавить работу с номером карты
func (db *DB) CreateNewReceipt(receiptInfo types.ReceiptInfoRequest) error {
	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("CreateNewReceipt: %v", err)
	}
	defer tx.Rollback()

	receiptID, err := db.insertReceipt(tx, receiptInfo)
	if err != nil {
		return fmt.Errorf("CreateNewReceipt: %v", err)
	}
	for _, item := range receiptInfo.Products {
		if err := db.insertReceiptProduct(tx, item, receiptID); err != nil {
			return fmt.Errorf("CreateNewReceiptProduct: %v", err)
		}
	}
	return tx.Commit()
}

func (db *DB) GetDepartmentInfo() ([]types.DepartmentInfoResponse, error) {
	result, err := db.getDepartment()
	if err != nil {
		return nil, fmt.Errorf("GetDepartmentInfo: %v", err)
	}
	departments := make([]types.DepartmentInfoResponse, len(result))
	for i := range result {
		departments[i] = result[i].ToDepartmentInfoResponse()
	}
	return departments, nil
}

func (db *DB) CreateNewEmployee(employeeInfo types.EmployeeInfoCreateRequest) error {
	_, err := db.db.Exec("insert into Employee (first_name, last_name, middle_name, position, salary, department_id) values ($1, $2, $3, $4, $5, $6)",
		employeeInfo.FirstName, employeeInfo.LastName, employeeInfo.MiddleName, employeeInfo.Position, employeeInfo.Salary, employeeInfo.DepartmentID)
	if err != nil {
		return fmt.Errorf("CreateNewEmployee: %v", err)
	}
	return nil
}

func (db *DB) GetEmployeeInfo() ([]types.EmployeeInfoResponse, error) {
	query := `
	select 
	    e.id,
	    e.first_name,
	    e.last_name,
	    e.middle_name,
	    e.position,
	    e.salary,
	    d.name
    from Employee as e
    left join Department as d on e.department_id = d.id`
	var employees []types.EmployeeInfoResponse
	rows, err := db.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetEmployeeInfo: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var employee types.EmployeeInfoResponse
		if err := rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.MiddleName, &employee.Position, &employee.Salary, &employee.Department); err != nil {
			return nil, fmt.Errorf("GetEmployeeInfo: %v", err)
		}
		employees = append(employees, employee)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetEmployeeInfo: %v", err)
	}
	return employees, nil
}

func (db *DB) DeleteEmployee(employeeInfo types.EmployeeInfoDeleteRequest) error {
	_, err := db.db.Exec("delete from Employee where id = $1", employeeInfo.ID)
	if err != nil {
		return fmt.Errorf("DeleteEmployee: %v", err)
	}
	return nil
}

func (db *DB) getDepartment() ([]types.Department, error) {
	var departments []types.Department

	rows, err := db.db.Query("select id, name, location, employee_count from Department")
	if err != nil {
		return nil, fmt.Errorf("getDepartment: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var department types.Department
		if err := rows.Scan(&department.ID, &department.Name, &department.Location, &department.EmployeeCount); err != nil {
			return nil, fmt.Errorf("getDepartment: %v", err)
		}
		departments = append(departments, department)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getDepartment: %v", err)
	}
	return departments, nil
}

func (db *DB) getEmployeeByPosition(position string) ([]types.Employee, error) {
	var employees []types.Employee

	rows, err := db.db.Query("select id, first_name, last_name, middle_name, salary from Employee where position = $1", position)
	if err != nil {
		return nil, fmt.Errorf("GetEmployeeByPosition: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var employee types.Employee
		if err := rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.MiddleName, &employee.Salary); err != nil {
			return nil, fmt.Errorf("GetEmployeeByPosition: %v", err)
		}
		employees = append(employees, employee)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetEmployeeByPosition: %v", err)
	}
	return employees, nil
}

func (db *DB) insertReceipt(tx *sql.Tx, receipt types.ReceiptInfoRequest) (int64, error) {
	var receiptID int64
	err := tx.QueryRow(
		"insert into Receipt (total_amount, employee_id) values ($1, $2) returning id",
		0, receipt.TellerID,
	).Scan(&receiptID)
	return receiptID, err
}

func (db *DB) insertReceiptProduct(tx *sql.Tx, receiptProduct types.ReceiptProductInfoRequest, receiptID int64) error {
	_, err := tx.Exec("insert into Receipt_Product (receipt_id, product_id, quantity, amount, price_at_purchase) values ($1, $2, $3, $4, $5)",
		receiptID, receiptProduct.ProductID, receiptProduct.Quantity, receiptProduct.Amount, receiptProduct.Price)
	return err
}

// GetAllProduct переработать: подстроить под интерфейс, возвращать будет джоин данные
//func (db *DB) GetAllProduct() ([]types.Product, error) {
//	var products []types.Product
//
//	rows, err := db.db.Query("select id, name, price, category, quantity_in_stock from Product")
//	if err != nil {
//		return nil, fmt.Errorf("GetAllProduct: %v", err)
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var product types.Product
//		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Category, &product.Quantity); err != nil {
//			return nil, fmt.Errorf("GetAllProduct: %v", err)
//		}
//		products = append(products, product)
//	}
//	if err := rows.Err(); err != nil {
//		return nil, fmt.Errorf("GetAllProduct: %v", err)
//	}
//	return products, nil
//}
//
//func (db *DB) GetAllSupplierInfo() ([]SupplierInfo, error) {
//	var suppliersInfo []SupplierInfo
//
//	rows, err := db.db.Query("select id, name from Supplier")
//	if err != nil {
//		return nil, fmt.Errorf("GetAllSupplierInfo: %v", err)
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var supplierInfo SupplierInfo
//		if err := rows.Scan(&supplierInfo.ID, &supplierInfo.Name); err != nil {
//			return nil, fmt.Errorf("GetAllSupplierInfo: %v", err)
//		}
//		suppliersInfo = append(suppliersInfo, supplierInfo)
//	}
//	if err := rows.Err(); err != nil {
//		return nil, fmt.Errorf("GetAllSupplierInfo: %v", err)
//	}
//	return suppliersInfo, nil
//}
//
//func (db *DB) GetProductBySupplierID(supplierID int64) ([]Product, error) {
//	var products []Product
//
//	query := `
//    SELECT p.id, p.name, p.price, p.category, p.quantity_in_stock
//    FROM Supplier s
//    JOIN Supplier_Order so ON s.id = so.supplier_id
//    JOIN Supplier_Order_Items soi ON so.id = soi.order_id
//    JOIN Product p ON soi.product_id = p.id
//    WHERE s.id = $1;
//    `
//
//	rows, err := db.db.Query(query, supplierID)
//	if err != nil {
//		return nil, fmt.Errorf("GetProductBySupplierID: %v", err)
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var product Product
//		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Category, &product.Quantity); err != nil {
//			return nil, fmt.Errorf("GetProductBySupplierID: %v", err)
//		}
//		products = append(products, product)
//	}
//	if err := rows.Err(); err != nil {
//		return nil, fmt.Errorf("GetProductBySupplierID: %v", err)
//	}
//	return products, nil
//}
//
//func (db *DB) GetAllEmployee() ([]Employee, error) {
//	var employees []Employee
//
//	rows, err := db.db.Query("select id, first_name, last_name, middle_name, position, salary from Employee")
//	if err != nil {
//		return nil, fmt.Errorf("GetAllEmployee: %v", err)
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var employee Employee
//		if err := rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.MiddleName, &employee.Position, &employee.Salary); err != nil {
//			return nil, fmt.Errorf("GetAllEmployee: %v", err)
//		}
//		employees = append(employees, employee)
//	}
//	if err := rows.Err(); err != nil {
//		return nil, fmt.Errorf("GetAllEmployee: %v", err)
//	}
//	return employees, nil
//}
//
//func (db *DB) GetAllReceipt() ([]Receipt, error) {
//	var receipts []Receipt
//
//	rows, err := db.db.Query("select id, date_time, total_amount from Receipt")
//	if err != nil {
//		return nil, fmt.Errorf("GetAllReceipt: %v", err)
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var receipt Receipt
//		if err := rows.Scan(&receipt.ID, &receipt.Date, &receipt.Total); err != nil {
//			return nil, fmt.Errorf("GetAllReceipt: %v", err)
//		}
//
//		receipts = append(receipts, receipt)
//	}
//	if err := rows.Err(); err != nil {
//		return nil, fmt.Errorf("GetAllReceipt: %v", err)
//	}
//
//	var wg sync.WaitGroup
//	var mu sync.Mutex
//	errChan := make(chan error, len(receipts))
//
//	for i := range receipts {
//		wg.Add(1)
//
//		go func(i int) {
//			defer wg.Done()
//
//			products, err := db.getReceiptProductByProductID(receipts[i].ID)
//			if err != nil {
//				errChan <- fmt.Errorf("GetAllReceipt: %v", err)
//				return
//			}
//
//			mu.Lock()
//			receipts[i].ReceiptProducts = products
//			mu.Unlock()
//		}(i)
//	}
//
//	wg.Wait()
//	close(errChan)
//
//	if len(errChan) > 0 {
//		return nil, <-errChan
//	}
//
//	return receipts, nil
//}
//
//func (db *DB) getReceiptProductByProductID(receiptID int64) ([]ReceiptProduct, error) {
//	var products []ReceiptProduct
//	query := `
//		SELECT p.name, rp.quantity, rp.amount, rp.price_at_purchase
//		FROM Receipt_Product as rp
//		JOIN Product as p ON rp.product_id = p.id
//		WHERE rp.receipt_id = $1
//	`
//
//	rows, err := db.db.Query(query, receiptID)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var rp ReceiptProduct
//		if err := rows.Scan(&rp.Name, &rp.Quantity, &rp.Amount, &rp.Price); err != nil {
//			return nil, err
//		}
//		products = append(products, rp)
//	}
//
//	return products, nil
//}
