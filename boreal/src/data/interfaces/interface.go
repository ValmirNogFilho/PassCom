// Package dao implements the interfaces for the database of the server.
package interfaces

import (
	"passcom/boreal/src/models"

	"github.com/google/uuid"
)

type FlightDAO interface {
	FindAll() []*models.Flight
	Insert(*models.Flight)
	Update(*models.Flight) error
	Delete(*models.Flight) error
	FindById(uuid.UUID) (*models.Flight, error)
	FindBySource(uuid.UUID) ([]*models.Flight, error)
	FindBySourceAndDest(uuid.UUID, uuid.UUID) (*models.Flight, error)
	BreadthFirstSearch(source uuid.UUID, dest uuid.UUID) ([]*models.Flight, error)
	DeleteAll()
	New()
}

type ClientDAO interface {
	FindAll() []*models.Client
	Insert(*models.Client)
	Update(*models.Client) error
	Delete(models.Client)
	FindById(uuid.UUID) (*models.Client, error)
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
	FindAll() []*models.Airport
	Insert(*models.Airport)
	Update(*models.Airport) error
	Delete(*models.Airport)
	FindById(uuid.UUID) (*models.Airport, error)
	New()
	FindByName(name string) *models.Airport
}
