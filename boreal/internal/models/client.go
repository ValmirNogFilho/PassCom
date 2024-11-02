package models

import (
	"github.com/google/uuid"
)

type Client struct {
	Id             uuid.UUID `json:"Id"`
	Name           string    `json:"Name"`
	Username       string    `json:"Username"`
	Password       string    `json:"Password"`
	Client_flights []*Ticket `json:"Client_flights"`
}
