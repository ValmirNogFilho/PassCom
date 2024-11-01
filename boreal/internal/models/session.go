package models

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID                 uuid.UUID
	ClientID           uuid.UUID
	LastTimeActive     time.Time
	Reservations       map[uuid.UUID]Reservation
	FailedReservations chan string
	Mu                 sync.RWMutex
}
