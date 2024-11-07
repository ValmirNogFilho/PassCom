package models

import (
	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	FlightSrcCity string
	FlightDestCity string
	ClientId uint   `gorm:"not null;constraint:OnDelete:CASCADE"` // Chave estrangeira para Client
	FlightId uint   `gorm:"not null;constraint:OnDelete:CASCADE"` // Chave estrangeira para Flight
	UniqueId string `gorm:"unique_id;unique"`

	Client Client `gorm:"foreignKey:ClientId;references:ID"` // Relacionamento many-to-one com Client
	Flight Flight `gorm:"foreignKey:FlightId;references:ID"` // Relacionamento many-to-one com Flight
}
