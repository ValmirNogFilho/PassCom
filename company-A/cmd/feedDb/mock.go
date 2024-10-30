package main

import (
	"encoding/json"
	"log"
	"os"
	"vendepass/internal/dao"
	"vendepass/internal/models"
)

func mockAirports() {
	file, err := os.Open("internal/stubs/airports.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var airports []models.Airport

	err = json.NewDecoder(file).Decode(&airports)
	if err != nil {
		log.Fatal(err)
	}

	dao := dao.GetAirportDAO()
	for _, airport := range airports {
		dao.Insert(airport)
	}
}

func mockClients() {
	file, err := os.Open("internal/stubs/clients.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var clients []models.Client

	err = json.NewDecoder(file).Decode(&clients)
	if err != nil {
		log.Fatal(err)
	}

	dao := dao.GetClientDAO()
	for _, client := range clients {
		dao.Insert(client)
	}
}

func mockFlights() {
	file, err := os.Open("internal/stubs/flights.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var flights []models.Flight

	err = json.NewDecoder(file).Decode(&flights)
	if err != nil {
		log.Fatal(err)
	}

	flightdao := dao.GetFlightDAO()
	airportdao := dao.GetAirportDAO()
	for _, flight := range flights {
		src, _ := airportdao.FindById(flight.OriginAirportID)
		dest, _ := airportdao.FindById(flight.DestinationAirportID)
		flight.OriginAirport = *src
		flight.DestinationAirport = *dest
		flightdao.Insert(flight)
	}
}
