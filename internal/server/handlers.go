package server

import (
	"db5/internal/db"
	"db5/internal/types"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

func CreateProductInfoHandler(store db.Store) *ProductInfoHandler {
	return &ProductInfoHandler{
		store: store,
	}
}

type ProductInfoHandler struct {
	store db.Store
}

func (pi *ProductInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		pi.GetProductInfo(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func (pi *ProductInfoHandler) GetProductInfo(w http.ResponseWriter, r *http.Request) {
	productInfo, err := pi.store.GetProductInfo()
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	jsonData, err := json.Marshal(productInfo)
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func CreateProductHandler(store db.Store) *ProductHandler {
	return &ProductHandler{
		store: store,
	}
}

type ProductHandler struct {
	store db.Store
}

func (p *ProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		p.GetProduct(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func (p *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	product, err := p.store.GetFullProductInfo()
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}
	jsonData, err := json.Marshal(product)
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func CreateEmployeeHandler(store db.Store) *EmployeeHandler {
	return &EmployeeHandler{
		store: store,
	}
}

type EmployeeHandler struct {
	store db.Store
}

func (e *EmployeeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		e.GetEmployee(w, r)
	case "POST":
		e.PostEmployee(w, r)
	case "DELETE":
		e.DeleteEmployee(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func (e *EmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	employee, err := e.store.GetEmployeeInfo()
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}
	jsonData, err := json.Marshal(employee)
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (e *EmployeeHandler) PostEmployee(w http.ResponseWriter, r *http.Request) {
	var employee types.EmployeeInfoCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	if err := e.store.CreateNewEmployee(employee); err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (e *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	var employee types.EmployeeInfoDeleteRequest

	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	if err := e.store.DeleteEmployee(employee); err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func CreateEmployeeTellerHandler(store db.Store) *EmployeeTellerHandler {
	return &EmployeeTellerHandler{
		store: store,
	}
}

type EmployeeTellerHandler struct {
	store db.Store
}

func (e *EmployeeTellerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		e.GetEmployeeTeller(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func (e *EmployeeTellerHandler) GetEmployeeTeller(w http.ResponseWriter, r *http.Request) {
	teller, err := e.store.GetTellerInfo()
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	jsonData, err := json.Marshal(teller)
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func CreateReceiptHandler(store db.Store) *ReceiptHandler {
	return &ReceiptHandler{
		store: store,
	}
}

type ReceiptHandler struct {
	store db.Store
}

func (rh *ReceiptHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rh.GetReceipt(w, r)
	case "POST":
		rh.PostReceipt(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func (rh *ReceiptHandler) GetReceipt(w http.ResponseWriter, r *http.Request) {
	receipt, err := rh.store.GetFullReceiptInfo()
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	jsonData, err := json.Marshal(receipt)
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (rh *ReceiptHandler) PostReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt types.ReceiptInfoRequest

	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	if err := rh.store.CreateNewReceipt(receipt); err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func CreateDepartmentInfoHandler(store db.Store) *DepartmentInfoHandler {
	return &DepartmentInfoHandler{
		store: store,
	}
}

type DepartmentInfoHandler struct {
	store db.Store
}

func (d *DepartmentInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		d.GetDepartmentInfo(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func (d *DepartmentInfoHandler) GetDepartmentInfo(w http.ResponseWriter, r *http.Request) {
	departmentInfo, err := d.store.GetDepartmentInfo()
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	jsonData, err := json.Marshal(departmentInfo)
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func CreateSupplierInfoHandler(store db.Store) *SupplierInfoHandler {
	return &SupplierInfoHandler{
		store: store,
	}
}

type SupplierInfoHandler struct {
	store db.Store
}

func (s *SupplierInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.GetSupplierInfo(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func (s *SupplierInfoHandler) GetSupplierInfo(w http.ResponseWriter, r *http.Request) {
	supplierInfo, err := s.store.GetSupplierInfo()
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	jsonData, err := json.Marshal(supplierInfo)
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func CreateSupplierProductHandler(store db.Store) *SupplierProductHandler {
	return &SupplierProductHandler{
		store: store,
	}
}

type SupplierProductHandler struct {
	store db.Store
}

func (sp *SupplierProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		sp.GetSupplierProductInfo(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func (sp *SupplierProductHandler) GetSupplierProductInfo(w http.ResponseWriter, r *http.Request) {
	strSupplierID := r.PathValue("id")
	supplierID, err := strconv.ParseInt(strSupplierID, 10, 64)
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	supplierProduct, err := sp.store.GetProductInfoBySupplier(supplierID)
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	jsonData, err := json.Marshal(supplierProduct)
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func CreateOrderHandler(store db.Store) *OrderHandler {
	return &OrderHandler{
		store: store,
	}
}

type OrderHandler struct {
	store db.Store
}

func (o *OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		o.GetOrderInfo(w, r)
	case "POST":
		o.PostOrderInfo(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func (o *OrderHandler) GetOrderInfo(w http.ResponseWriter, r *http.Request) {
	order, err := o.store.GetFullSupplierOrderInfo()
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (o *OrderHandler) PostOrderInfo(w http.ResponseWriter, r *http.Request) {
	var orderInfo types.SupplierOrderInfoRequest

	if err := json.NewDecoder(r.Body).Decode(&orderInfo); err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	err := o.store.CreateNewSupplierOrder(orderInfo)
	if err != nil {
		InternalServerErrorHandler(w, r)
		slog.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
