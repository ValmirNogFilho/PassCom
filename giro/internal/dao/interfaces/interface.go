// Package dao implements the interfaces for the database of the server.
package interfaces

import (
	"giro/internal/models"

	"github.com/google/uuid"
)

type FlightDAO interface {
	FindAll() []models.Flight
	Insert(models.Flight)
	Update(models.Flight) error
	Delete(models.Flight) error
	FindById(uint) (*models.Flight, error)
	FindBySource(uint) ([]models.Flight, error)
	FindBySourceAndDest(uint, uint) ([]models.Flight, error)
	FindByCompany(string) ([]models.Flight, error)
	FindByUniqueId(string) (*models.Flight, error)
	FindPathBFS(uint, uint) ([]models.Flight, error)
	DeleteByUniqueId(string) error
	DeleteByCompany(string) error
	DeleteAll()
	New()
}

type ClientDAO interface {
	FindAll() []models.Client
	Insert(models.Client)
	Update(models.Client) error
	Delete(models.Client)
	FindById(uint) (*models.Client, error)
	FindByUsername(username string) (*models.Client, error)
	New()
}

type SessionDAO interface {
	FindAll() []*models.Session
	Insert(*models.Session)
	Update(*models.Session) error
	Delete(*models.Session)
	FindById(uuid.UUID) (*models.Session, error)
	DeleteAll()
	New()
}

type AirportDAO interface {
	FindAll() []models.Airport
	Insert(models.Airport)
	Update(models.Airport) error
	Delete(models.Airport)
	FindById(uint) (*models.Airport, error)
	New()
	FindByName(name string) *models.Airport
}

type TicketDAO interface {
	FindAll() []models.Ticket
	Insert(models.Ticket)
	Update(models.Ticket) error
	Delete(models.Ticket)
	FindById(uint) (*models.Ticket, error)
	FindByUniqueId(string) (*models.Ticket, error)
	DeleteByUniqueId(string) error
	New()
}

type MessageDAO interface {
	FindAll() []models.Message
	Insert(models.Message)
	Update(models.Message) error
	Delete(models.Message)
	FindById(uint) (*models.Message, error)
	New()
	FindByName(name string) *models.Message
}
