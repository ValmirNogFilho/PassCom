package models

const (
	TypePurchase = "purchase"
	TypeCancel   = "cancel"
)

type Transaction struct {
	Type     string
	FlightId string
}
