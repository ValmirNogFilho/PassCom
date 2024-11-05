package dao

import (
	"boreal/internal/dao/interfaces"
	"boreal/internal/models"
	"sync"

	"github.com/google/uuid"
)

var airportDao interfaces.AirportDAO
var flightDao interfaces.FlightDAO
var clientDao interfaces.ClientDAO
var sessionDao interfaces.SessionDAO
var ticketDao interfaces.TicketDAO

func GetFlightDAO() interfaces.FlightDAO {
	if flightDao == nil {
		flightDao = &DBFlightDAO{}
		flightDao.New()
	}

	return flightDao
}

func GetClientDAO() interfaces.ClientDAO {
	if clientDao == nil {
		clientDao = &DBClientDAO{}
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
		airportDao = &DBAirportDAO{}
		airportDao.New()
	}

	return airportDao
}

func GetTicketDAO() interfaces.TicketDAO {
	if ticketDao == nil {
		ticketDao = &DBTicketDAO{}
		ticketDao.New()
	}

	return ticketDao
}
