package data

import (
	"boreal/internal/models"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

// MemoryFlightDAO is a data access object (DAO) for managing flight data in memory.
// It provides methods for inserting, updating, deleting, and retrieving flights.
// It also includes a breadth-first search algorithm for finding the shortest path between airports.
type MemoryFlightDAO struct {
	data map[uuid.UUID]map[uuid.UUID]*models.Flight
	mu   sync.RWMutex
}

// New initializes the MemoryFlightDAO by reading flight data from a JSON file and populating the internal data structure.
// It also creates a new session queue for each flight.
func (dao *MemoryFlightDAO) New() {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	jsonPath := filepath.Join(baseDir, "internal", "assets", "airports.json")

	b, _ := os.ReadFile(jsonPath)

	var flights []models.FlightJSON

	json.Unmarshal(b, &flights)

	for _, f := range flights {
		flight := models.Flight{
			Id:              f.Id,
			SourceAirportId: f.SourceAirportId,
			DestAirportId:   f.DestAirportId,
			Passengers:      f.Passengers,
			Seats:           f.Seats,
			Queue:           make(chan *models.Session, 10),
		}

		if dao.data[flight.SourceAirportId] == nil {
			dao.data[flight.SourceAirportId] = make(map[uuid.UUID]*models.Flight)
		}

		dao.data[flight.SourceAirportId][flight.DestAirportId] = &flight
	}
}

// FindAll retrieves all flights from the memory data structure.
// It iterates through the map of flights and appends each flight to a slice.
//
// Return:
//   - A slice of pointers to models.Flight, representing all flights in the data structure.
//   - If no flights are found, an empty slice is returned.
func (dao *MemoryFlightDAO) FindAll() []*models.Flight {
	dao.mu.RLock()
	defer dao.mu.RUnlock()
	var v []*models.Flight
	for _, array := range dao.data {
		for _, flight := range array {
			v = append(v, flight)
		}
	}

	return v
}

// Insert adds a new flight to the memory data structure.
// It generates a new UUID for the flight, sets the flight's ID, and creates a new session queue.
// If the source airport ID does not exist in the data structure, a new map is created for that airport.
// The flight is then added to the data structure using the source airport ID and destination airport ID as keys.
//
// Parameters:
//   - t *models.Flight: A pointer to the flight to be inserted. The flight's ID, source airport ID, destination airport ID,
//   - passengers, and seats should be set before calling this function.
func (dao *MemoryFlightDAO) Insert(t *models.Flight) {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	id := uuid.New()

	t.Id = id

	if dao.data[t.SourceAirportId] == nil {
		dao.data[t.SourceAirportId] = make(map[uuid.UUID]*models.Flight)
	}

	dao.data[t.SourceAirportId][t.DestAirportId] = t
	t.Queue = make(chan *models.Session)
}

// Update updates an existing flight in the memory data structure.
// It checks if the flight exists in the data structure based on the source airport ID and destination airport ID.
// If the flight is found, it updates the flight's details in the data structure.
// If the flight is not found, it returns an error.
//
// Parameters:
//   - t *models.Flight: A pointer to the flight to be updated. The flight's ID, source airport ID, destination airport ID,
//     passengers, and seats should be set before calling this function.
//
// Return:
//   - An error if the flight is not found in the data structure.
//   - nil if the flight is successfully updated.
func (dao *MemoryFlightDAO) Update(t *models.Flight) error {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	_, exists := dao.data[t.SourceAirportId][t.DestAirportId]

	if !exists {
		return errors.New("not found")
	}

	dao.data[t.SourceAirportId][t.DestAirportId] = t

	return nil
}

// Delete removes a flight from the memory data structure based on the provided flight object.
// It deletes the flight from the map using the source airport ID and destination airport ID as keys.
// After deletion, it checks if the flight still exists in the data structure.
// If the flight still exists, it returns an error indicating that the deletion was unsuccessful.
// If the flight is successfully deleted, it returns nil.
//
// Parameters:
//   - t *models.Flight: A pointer to the flight to be deleted. The flight's source airport ID and destination airport ID
//     should be set before calling this function.
//
// Return:
//   - An error if the flight is not found in the data structure after deletion.
//   - nil if the flight is successfully deleted.
func (dao *MemoryFlightDAO) Delete(t *models.Flight) error {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	delete(dao.data[t.SourceAirportId], t.DestAirportId)

	_, exists := dao.data[t.SourceAirportId][t.DestAirportId]

	if exists {
		return errors.New("delete was unsuccessful")
	}

	return nil
}

// FindById retrieves a flight from the memory data structure based on its unique ID.
// It iterates through the map of flights and checks if the flight's ID matches the provided ID.
// If a matching flight is found, it is returned along with a nil error.
// If no matching flight is found, nil is returned along with an error indicating that the flight was not found.
//
// Parameters:
//   - id uuid.UUID: The unique ID of the flight to be retrieved.
//
// Return:
//   - *models.Flight: A pointer to the flight with the matching ID, or nil if no matching flight is found.
//   - error: An error indicating that the flight was not found, or nil if the flight is successfully retrieved.
func (dao *MemoryFlightDAO) FindById(id uuid.UUID) (*models.Flight, error) {
	dao.mu.RLock()
	defer dao.mu.RUnlock()
	for _, array := range dao.data {
		for _, flight := range array {
			if flight.Id == id {
				return flight, nil
			}
		}
	}

	return nil, errors.New("flight not found")
}

// FindBySource retrieves all flights departing from a specific airport.
// It iterates through the map of flights and checks if the source airport ID matches the provided ID.
// If a matching flight is found, it is added to a slice.
// If no matching flights are found, nil is returned along with an error indicating that the airport was not found.
//
// Parameters:
//   - id uuid.UUID: The unique ID of the airport from which flights are to be retrieved.
//
// Return:
//   - []*models.Flight: A slice of pointers to flights departing from the specified airport.
//   - error: An error indicating that the airport was not found, or nil if flights are successfully retrieved.
func (dao *MemoryFlightDAO) FindBySource(id uuid.UUID) ([]*models.Flight, error) {
	dao.mu.RLock()
	defer dao.mu.RUnlock()
	t, exists := dao.data[id]

	if !exists {
		return nil, errors.New("airport not found")
	}

	flights := make([]*models.Flight, 0, len(t))

	for _, flight := range t {
		flights = append(flights, flight)
	}

	return flights, nil
}

// FindBySourceAndDest retrieves a specific flight from the memory data structure based on its source and destination airport IDs.
// It checks if a flight exists in the data structure with the provided source and destination airport IDs.
// If a matching flight is found, it is returned along with a nil error.
// If no matching flight is found, nil is returned along with an error indicating that the flight was not found.
//
// Parameters:
//   - source uuid.UUID: The unique ID of the source airport.
//   - dest uuid.UUID: The unique ID of the destination airport.
//
// Return:
//   - *models.Flight: A pointer to the flight with the matching source and destination airport IDs, or nil if no matching flight is found.
//   - error: An error indicating that the flight was not found, or nil if the flight is successfully retrieved.
func (dao *MemoryFlightDAO) FindBySourceAndDest(source uuid.UUID, dest uuid.UUID) (*models.Flight, error) {
	dao.mu.RLock()
	defer dao.mu.RUnlock()
	t, exists := dao.data[source][dest]

	if !exists {
		return nil, errors.New("flight not found")
	}

	return t, nil
}

// BreadthFirstSearch performs a breadth-first search on the flight data structure to find the shortest path between two airports.
// It uses the source airport ID and destination airport ID as input parameters.
// The function returns a slice of pointers to models.Flight representing the shortest path between the source and destination airports.
// If no route is available, it returns an error indicating that no route was found.
//
// Parameters:
//   - source uuid.UUID: The unique ID of the source airport.
//   - dest uuid.UUID: The unique ID of the destination airport.
//
// Return:
//   - []*models.Flight: A slice of pointers to flights representing the shortest path between the source and destination airports.
//   - error: An error indicating that no route was found, or nil if a route is successfully retrieved.
func (dao *MemoryFlightDAO) BreadthFirstSearch(source uuid.UUID, dest uuid.UUID) ([]*models.Flight, error) {
	dao.mu.RLock()
	defer dao.mu.RUnlock()
	visited := make(map[uuid.UUID]bool, len(dao.data))
	queue := []uuid.UUID{source}
	visited[source] = true
	parent := make(map[uuid.UUID]uuid.UUID, len(dao.data))
	parent[source] = source

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current == dest {
			break
		}

		for neighbor, flight := range dao.data[current] {
			if !visited[neighbor] && flight.Seats > 0 {
				visited[neighbor] = true
				queue = append(queue, neighbor)
				parent[neighbor] = current
			}
		}
	}

	path := []*models.Flight{}
	current := dest

	if !visited[dest] {
		return nil, errors.New("no route available")
	}

	for current != source {
		prev := parent[current]
		flight := dao.data[prev][current]
		path = append([]*models.Flight{flight}, path...)
		current = prev
	}

	return path, nil
}

// DeleteAll removes all flights from the memory data structure.
// It resets the internal map of flights to an empty map, effectively deleting all flights.
// This function is useful for testing or resetting the data structure to its initial state.
func (dao *MemoryFlightDAO) DeleteAll() {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	dao.data = make(map[uuid.UUID]map[uuid.UUID]*models.Flight)
}
