package models

import (
	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	ClientId uint `gorm:"not null"` // Chave estrangeira para Client
	FlightId uint `gorm:"not null"` // Chave estrangeira para Flight

	Client Client `gorm:"foreignKey:ClientId;references:ID"` // Relacionamento many-to-one com Client
	Flight Flight `gorm:"foreignKey:FlightId;references:ID"` // Relacionamento many-to-one com Flight
}
