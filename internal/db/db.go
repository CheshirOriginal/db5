package db

import (
	"database/sql"
	"db5/config"
	"db5/internal/types"
	"fmt"
	"sync"

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
	GetSupplierInfo() ([]types.SupplierInfoResponse, error)
	GetProductInfoBySupplier(supplierID int64) ([]types.ProductInfoBySupplierResponse, error)
	CreateNewSupplierOrder(supplierOrderInfo types.SupplierOrderInfoRequest) error
	GetFullProductInfo() ([]types.FullProductInfoResponse, error)
	GetFullReceiptInfo() ([]types.FullReceiptInfoResponse, error)
	GetFullSupplierOrderInfo() ([]types.FullSupplierOrderInfoResponse, error)
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

	db.db.SetMaxOpenConns(10)
	db.db.SetMaxIdleConns(5)
	db.db.SetConnMaxLifetime(time.Minute * 5)

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

func (db *DB) GetSupplierInfo() ([]types.SupplierInfoResponse, error) {
	var suppliers []types.SupplierInfoResponse
	rows, err := db.db.Query("select id, name from Supplier")
	if err != nil {
		return nil, fmt.Errorf("GetSupplierInfo: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var supplier types.SupplierInfoResponse
		if err := rows.Scan(&supplier.ID, &supplier.Name); err != nil {
			return nil, fmt.Errorf("GetSupplierInfo: %v", err)
		}
		suppliers = append(suppliers, supplier)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetSupplierInfo: %v", err)
	}
	return suppliers, nil
}

func (db *DB) GetProductInfoBySupplier(supplierID int64) ([]types.ProductInfoBySupplierResponse, error) {
	var query = `
	select DISTINCT
	p.name,
	p.id
	from Product as p
	join Supplier_Order_Items as soi on p.id = soi.product_id
	join Supplier_Order as so on so.id = soi.order_id
	where so.supplier_id = $1`

	var products []types.ProductInfoBySupplierResponse
	rows, err := db.db.Query(query, supplierID)
	if err != nil {
		return nil, fmt.Errorf("GetProductInfoBySupplier: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var product types.ProductInfoBySupplierResponse
		if err := rows.Scan(&product.Name, &product.ID); err != nil {
			return nil, fmt.Errorf("GetProductInfoBySupplier: %v", err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetProductInfoBySupplier: %v", err)
	}
	return products, nil
}

func (db *DB) CreateNewSupplierOrder(supplierOrderInfo types.SupplierOrderInfoRequest) error {
	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("CreateNewSupplierOrder: %v", err)
	}
	defer tx.Rollback()

	supplierOrderID, err := db.insertSupplierOrder(tx, supplierOrderInfo)
	if err != nil {
		return fmt.Errorf("CreateNewSupplierOrder: %v", err)
	}
	for _, item := range supplierOrderInfo.SupplierOrderItems {
		if err := db.insertSupplierOrderItem(tx, item, supplierOrderID); err != nil {
			return fmt.Errorf("CreateNewSupplierOrderItem: %v", err)
		}
	}
	return tx.Commit()
}

func (db *DB) GetFullProductInfo() ([]types.FullProductInfoResponse, error) {
	query := `
	select
	p.name,
	p.price,
	p.category,
	p.quantity_in_stock,
	d.name
	from Product as p
	join Department as d on p.department_id = d.id`
	var Products []types.FullProductInfoResponse
	rows, err := db.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetFullProductInfo: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var product types.FullProductInfoResponse
		if err := rows.Scan(&product.Name, &product.Price, &product.Category, &product.Quantity, &product.DepartmentName); err != nil {
			return nil, fmt.Errorf("GetFullProductInfo: %v", err)
		}
		Products = append(Products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetFullProductInfo: %v", err)
	}
	return Products, nil
}

func (db *DB) GetFullReceiptInfo() ([]types.FullReceiptInfoResponse, error) {
	query := `
	select
	e.first_name,
	e.last_name,
	e.middle_name,
	r.id,
	r.number,
	r.date_time,
	r.total_amount,
	lc.number
	from Receipt as r
	join Employee as e on e.id = r.employee_id
	join Loyalty_Card as lc on lc.id = r.loyalty_card_id
	order by r.id
	`
	var receipts []types.FullReceiptInfoResponse

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetFullReceiptInfo: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var receipt types.FullReceiptInfoResponse
		if err := rows.Scan(&receipt.TellerFirstName, &receipt.TellerLastName, &receipt.TellerMiddleName, &receipt.ID, &receipt.Number, &receipt.Date, &receipt.Total, &receipt.LoyaltyCardNumber); err != nil {
			return nil, fmt.Errorf("GetFullReceiptInfo: %v", err)
		}
		receipts = append(receipts, receipt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetFullReceiptInfo: %v", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(receipts))

	for i := range receipts {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			products, err := db.getReceiptProductByProductID(receipts[i].ID)
			if err != nil {
				errChan <- fmt.Errorf("GetFullReceiptInfo: %v", err)
				return
			}
			mu.Lock()
			receipts[i].Products = products
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return receipts, nil
}

func (db *DB) GetFullSupplierOrderInfo() ([]types.FullSupplierOrderInfoResponse, error) {
	query := `
	select
	so.order_date,
	so.date_of_receipt,
	so.id,
	so.total_amount,
	s.name
 	from Supplier_Order as so
 	join Supplier as s on so.supplier_id = s.id
 	order by so.id`
	var supplierOrders []types.FullSupplierOrderInfoResponse
	var nt sql.NullTime

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetFullSupplierOrderInfo: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var supplier types.FullSupplierOrderInfoResponse
		if err := rows.Scan(&supplier.OrderDate, &nt, &supplier.ID, &supplier.Total, &supplier.SupplierName); err != nil {
			return nil, fmt.Errorf("GetFullSupplierOrderInfo: %v", err)
		}
		if nt.Valid {
			supplier.DateOfReceipt = nt.Time
		}
		supplierOrders = append(supplierOrders, supplier)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetFullSupplierOrderInfo: %v", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(supplierOrders))

	for i := range supplierOrders {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			items, err := db.getSupplierOrderItemByOrderID(supplierOrders[i].ID)
			if err != nil {
				errChan <- fmt.Errorf("GetFullReceiptInfo: %v", err)
				return
			}
			mu.Lock()
			supplierOrders[i].SupplierOrderItems = items
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return supplierOrders, nil
}

func (db *DB) getSupplierOrderItemByOrderID(orderID int64) ([]types.SupplierOrderItemResponse, error) {
	var supplierOrderItems []types.SupplierOrderItemResponse
	query := `
		SELECT p.name, soi.quantity, soi.purchase_price
		FROM Supplier_Order_Items as soi
		JOIN Product as p ON soi.product_id = p.id
		WHERE soi.order_id = $1
	`

	rows, err := db.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var supplierOrderItem types.SupplierOrderItemResponse
		if err := rows.Scan(&supplierOrderItem.ProductName, &supplierOrderItem.Quantity, &supplierOrderItem.Price); err != nil {
			return nil, err
		}
		supplierOrderItem.Amount = supplierOrderItem.Price * float64(supplierOrderItem.Quantity)
		supplierOrderItems = append(supplierOrderItems, supplierOrderItem)
	}

	return supplierOrderItems, nil
}

func (db *DB) getReceiptProductByProductID(receiptID int64) ([]types.ReceiptProductResponse, error) {
	var products []types.ReceiptProductResponse
	query := `
		SELECT p.name, rp.quantity, rp.amount, rp.price_at_purchase
		FROM Receipt_Product as rp
		JOIN Product as p ON rp.product_id = p.id
		WHERE rp.receipt_id = $1
	`

	rows, err := db.db.Query(query, receiptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var receiptProduct types.ReceiptProductResponse
		if err := rows.Scan(&receiptProduct.Name, &receiptProduct.Quantity, &receiptProduct.Amount, &receiptProduct.Price); err != nil {
			return nil, err
		}
		products = append(products, receiptProduct)
	}

	return products, nil
}

func (db *DB) insertSupplierOrder(tx *sql.Tx, supplierOrderInfo types.SupplierOrderInfoRequest) (int64, error) {
	var supplierOrderID int64
	err := tx.QueryRow("insert into Supplier_Order (total_amount, supplier_id) values ($1, $2) returning id",
		0, supplierOrderInfo.SupplierID,
	).Scan(&supplierOrderID)
	if err != nil {
		return 0, fmt.Errorf("insertSupplierOrder: %v", err)
	}
	return supplierOrderID, nil
}

func (db *DB) insertSupplierOrderItem(tx *sql.Tx, supplierOrderItem types.SupplierOrderItemInfoRequest, orderID int64) error {
	_, err := tx.Exec("insert into Supplier_Order_Items (purchase_price, quantity, product_id, order_id) values ($1, $2, $3, $4)",
		supplierOrderItem.Price, supplierOrderItem.Quantity, supplierOrderItem.ProductID, orderID)
	return err
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
	var loyaltyCardId any

	if receipt.LoyaltyCardNumber == 0 {
		loyaltyCardId = nil
	} else {
		loyaltyCardId = receipt.LoyaltyCardNumber
	}

	err := tx.QueryRow(
		"insert into Receipt (total_amount, employee_id, loyalty_card_id) values ($1, $2, $3) returning id",
		0, receipt.TellerID, loyaltyCardId,
	).Scan(&receiptID)
	return receiptID, err
}

func (db *DB) insertReceiptProduct(tx *sql.Tx, receiptProduct types.ReceiptProductInfoRequest, receiptID int64) error {
	_, err := tx.Exec("insert into Receipt_Product (receipt_id, product_id, quantity, amount, price_at_purchase) values ($1, $2, $3, $4, $5)",
		receiptID, receiptProduct.ProductID, receiptProduct.Quantity, receiptProduct.Amount, receiptProduct.Price)
	return err
}
