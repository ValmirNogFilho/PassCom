package models

import (
	"github.com/google/uuid"
)

type Route struct {
	Path     []City
	FlightId uuid.UUID
}
