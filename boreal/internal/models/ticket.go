package models

import (
	"github.com/google/uuid"
)

type Ticket struct {
	Id       uuid.UUID
	ClientId uuid.UUID
	FlightId uuid.UUID
}
