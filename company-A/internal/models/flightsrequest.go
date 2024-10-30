package models

import "github.com/google/uuid"

type FlightsRequest struct {
	FlightIds []uuid.UUID
}
