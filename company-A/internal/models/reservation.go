package models

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	Id        uuid.UUID
	CreatedAt time.Time
	*Ticket
}
