package models

import (
	"sync"

	"github.com/google/uuid"
)

type Airport struct {
	Id   uuid.UUID `json:"Id"`
	Name string    `json:"Name"`
	City City      `json:"City"`
	Mu   sync.RWMutex
}
