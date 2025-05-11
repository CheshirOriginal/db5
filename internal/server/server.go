package server

import (
	"db5/internal/db"
	"net/http"
	"time"

	"github.com/rs/cors"
)

func CreateNewServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler,
	}
}

func CreateNewServerMux(store db.Store) *http.Handler {
	mux := http.NewServeMux()

	employeeHandler := CreateEmployeeHandler(store)
	employeeTeller := CreateEmployeeTellerHandler(store)
	receiptHandler := CreateReceiptHandler(store)
	departmentInfoHandler := CreateDepartmentInfoHandler(store)
	productHandler := CreateProductHandler(store)
	productInfoHandler := CreateProductInfoHandler(store)
	supplierInfoHandler := CreateSupplierInfoHandler(store)
	supplierProductHandler := CreateSupplierProductHandler(store)
	orderHandler := CreateOrderHandler(store)

	mux.Handle("/employee", employeeHandler)
	mux.Handle("/employee/teller/info", employeeTeller)
	mux.Handle("/receipt", receiptHandler)
	mux.Handle("/department/info", departmentInfoHandler)
	mux.Handle("/product", productHandler)
	mux.Handle("/product/info", productInfoHandler)
	mux.Handle("/supplier/info", supplierInfoHandler)
	mux.Handle("/supplier/product/{id}", supplierProductHandler)
	mux.Handle("/order", orderHandler)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
	}).Handler(mux)

	return &handler
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
