package models

import "github.com/google/uuid"

type CancelReservationRequest struct {
	ReservationId uuid.UUID
}
