package data

import (
	"boreal/src/data/interfaces"
	"boreal/src/models"
	"sync"

	"github.com/google/uuid"
)

var airportDao interfaces.AirportDAO
var flightDao interfaces.FlightDAO
var clientDao interfaces.ClientDAO
var sessionDao interfaces.SessionDAO

// GetFlightDAO returns a singleton instance of FlightDAO.
// If the instance does not exist, it creates a new one and initializes it.
//
// The FlightDAO interface provides methods for managing flight data.
// This implementation uses a MemoryFlightDAO, which stores flight data in memory.
//
// Parameters:
// None
//
// Return:
// flightDao (interfaces.FlightDAO) - A singleton instance of FlightDAO.
func GetFlightDAO() interfaces.FlightDAO {
	if flightDao == nil {
		flightDao = &MemoryFlightDAO{data: make(map[uuid.UUID]map[uuid.UUID]*models.Flight),
			mu: sync.RWMutex{}}
		flightDao.New()
	}

	return flightDao
}

func GetClientDAO() interfaces.ClientDAO {
	if clientDao == nil {
		clientDao = &MemoryClientDAO{data: make(map[uuid.UUID]*models.Client),
			mu: sync.RWMutex{}}
		clientDao.New()
	}

	return clientDao
}

func GetSessionDAO() interfaces.SessionDAO {
	if sessionDao == nil {
		sessionDao = &MemorySessionDAO{
			data: make(map[uuid.UUID]*models.Session),
			mu:   sync.RWMutex{}}
		sessionDao.New()
	}

	return sessionDao
}

func GetAirportDAO() interfaces.AirportDAO {
	if airportDao == nil {
		airportDao = &MemoryAirportDAO{data: make(map[uuid.UUID]*models.Airport),
			mu: sync.RWMutex{}}
		airportDao.New()
	}

	return airportDao
}
