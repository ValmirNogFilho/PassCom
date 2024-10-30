package models

import "github.com/google/uuid"

type CancelBuyRequest struct {
	TicketId uuid.UUID
}
