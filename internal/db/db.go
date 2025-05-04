package db

import (
	"database/sql"
	"db5/config"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

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

func (db *DB) GetAllProduct() ([]Product, error) {
	var products []Product

	rows, err := db.db.Query("select id, name, price, category, quantity_in_stock from Product")
	if err != nil {
		return nil, fmt.Errorf("GetAllProduct: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Category, &product.Quantity); err != nil {
			return nil, fmt.Errorf("GetAllProduct: %v", err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllProduct: %v", err)
	}
	return products, nil
}

func (db *DB) GetAllProductInfo() ([]ProductInfo, error) {
	var products []ProductInfo

	rows, err := db.db.Query("select id, name from Product")
	if err != nil {
		return nil, fmt.Errorf("GetAllProductInfo: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product ProductInfo
		if err := rows.Scan(&product.ID, &product.Name); err != nil {
			return nil, fmt.Errorf("GetAllProductInfo: %v", err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllProductInfo: %v", err)
	}
	return products, nil
}

func (db *DB) GetEmployeeByPosition(position string) ([]Employee, error) {
	var employees []Employee

	rows, err := db.db.Query("select id, first_name, last_name, middle_name, salary from Employee where position = $1", position)
	if err != nil {
		return nil, fmt.Errorf("GetEmployeeByPosition: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var employee Employee
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

func (db *DB) GetAllSupplierInfo() ([]SupplierInfo, error) {
	var suppliersInfo []SupplierInfo

	rows, err := db.db.Query("select id, name from Supplier")
	if err != nil {
		return nil, fmt.Errorf("GetAllSupplierInfo: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var supplierInfo SupplierInfo
		if err := rows.Scan(&supplierInfo.ID, &supplierInfo.Name); err != nil {
			return nil, fmt.Errorf("GetAllSupplierInfo: %v", err)
		}
		suppliersInfo = append(suppliersInfo, supplierInfo)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllSupplierInfo: %v", err)
	}
	return suppliersInfo, nil
}

func (db *DB) GetProductBySupplierID(supplierID int64) ([]Product, error) {
	var products []Product

	query := `
    SELECT p.id, p.name, p.price, p.category, p.quantity_in_stock
    FROM Supplier s
    JOIN Supplier_Order so ON s.id = so.supplier_id
    JOIN Supplier_Order_Items soi ON so.id = soi.order_id
    JOIN Product p ON soi.product_id = p.id
    WHERE s.id = $1;
    `

	rows, err := db.db.Query(query, supplierID)
	if err != nil {
		return nil, fmt.Errorf("GetProductBySupplierID: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Category, &product.Quantity); err != nil {
			return nil, fmt.Errorf("GetProductBySupplierID: %v", err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetProductBySupplierID: %v", err)
	}
	return products, nil
}

func (db *DB) GetAllEmployee() ([]Employee, error) {
	var employees []Employee

	rows, err := db.db.Query("select id, first_name, last_name, middle_name, position, salary from Employee")
	if err != nil {
		return nil, fmt.Errorf("GetAllEmployee: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var employee Employee
		if err := rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.MiddleName, &employee.Position, &employee.Salary); err != nil {
			return nil, fmt.Errorf("GetAllEmployee: %v", err)
		}
		employees = append(employees, employee)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllEmployee: %v", err)
	}
	return employees, nil
}

func (db *DB) GetAllReceipt() ([]Receipt, error) {
	var receipts []Receipt

	rows, err := db.db.Query("select id, number, date_time, total_amount from Receipt")
	if err != nil {
		return nil, fmt.Errorf("GetAllReceipt: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var receipt Receipt
		if err := rows.Scan(&receipt.ID, &receipt.Number, &receipt.Date, &receipt.Total); err != nil {
			return nil, fmt.Errorf("GetAllReceipt: %v", err)
		}

		receipts = append(receipts, receipt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllReceipt: %v", err)
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
				errChan <- fmt.Errorf("GetAllReceipt: %v", err)
				return
			}

			mu.Lock()
			receipts[i].ReceiptProducts = products
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

func (db *DB) getReceiptProductByProductID(receiptID int64) ([]ReceiptProduct, error) {
	var products []ReceiptProduct
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
		var rp ReceiptProduct
		if err := rows.Scan(&rp.Name, &rp.Quantity, &rp.Amount, &rp.Price); err != nil {
			return nil, err
		}
		products = append(products, rp)
	}

	return products, nil
}
