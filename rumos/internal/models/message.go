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

// CreateMessage creates a new Message instance with the provided parameters.
//
// The function generates a new UUID v7 for the message ID and creates a new Message instance.
// The 'from' and 'to' fields are initially set to empty strings.
// The provided vectorClock and body are assigned to the respective fields of the Message instance.
//
// Parameters:
// - from: The sender of the message.
// - to: The recipient of the message.
// - vectorClock: A map representing the vector clock of the message.
// - body: The content of the message. Can be of any serializable type.
//
// Returns:
// - A pointer to the newly created Message instance.
// - An error if the UUID generation fails.
func CreateMessage(from string, to string, vectorClock map[string]int, body interface{}) (*Message, error) {
	id, err := uuid.NewV7()

	if err != nil {
		return nil, err
	}

	message := &Message{
		Id:          id.String(),
		From:        from,
		To:          to,
		VectorClock: vectorClock,
		Body:        body,
	}

	return message, nil
}
