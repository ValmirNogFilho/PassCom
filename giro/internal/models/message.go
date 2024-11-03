package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Id          string         `json:"Id"`          // Serializado como string
	From        string         `json:"From"`        // Serializado como string
	To          string         `json:"To"`          // Serializado como string
	VectorClock map[string]int `json:"VectorClock"` // Mapeia como string para evitar problemas
	Body        interface{}    `json:"Body"`        // Pode ser qualquer tipo de dado serializável
}

func NewMessageId() (uuid.UUID, error) {
	return uuid.NewUUID()
}

func NewMessageIdString() (string, error) {
	id, err := NewMessageId()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
