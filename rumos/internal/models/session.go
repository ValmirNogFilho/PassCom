package models

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID             uuid.UUID
	ClientID       uint
	LastTimeActive time.Time
	Mu             sync.RWMutex
}
