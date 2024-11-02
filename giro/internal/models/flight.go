package models

import (
	"gorm.io/gorm"
)

type Flight struct {
	gorm.Model
	Company              string
	Price                uint
	OriginAirportID      uint    `gorm:"not null"`
	OriginAirport        Airport `gorm:"foreignKey:OriginAirportID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	DestinationAirportID uint    `gorm:"not null"`
	DestinationAirport   Airport `gorm:"foreignKey:DestinationAirportID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Seats                int
	Tickets              []Ticket `gorm:"foreignKey:FlightId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
