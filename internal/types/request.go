package types

type ReceiptInfoRequest struct {
	LoyaltyCardNumber int64                `json:"loyalty_card_number"`
	TellerID          int64                `json:"teller_id"`
	Products          []ReceiptProductInfo `json:"products"`
}

type ReceiptProductInfo struct {
	ProductID int64   `json:"product_id"`
	Quantity  int64   `json:"quantity"`
	Price     float64 `json:"price"`
	Amount    float64 `json:"amount"`
}
