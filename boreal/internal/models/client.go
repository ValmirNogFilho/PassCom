package models

import (
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Name          string   `gorm:"size:100"`
	Username      string   `gorm:"size:30"`
	Password      string   `gorm:"size:30"`
	ClientFlights []Ticket `gorm:"foreignKey:ClientId"` // Relacionamento one-to-many com Ticket
}
