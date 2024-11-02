// Package dao implements the Data Access Object for the database of the server.
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

// MemoryAirportDAO is a data access object (DAO) for managing airports in memory.
// It provides methods for retrieving, inserting, updating, and deleting airports.
type MemoryAirportDAO struct {
	data map[uuid.UUID]*models.Airport
	mu   sync.RWMutex
}

// New initializes the MemoryAirportDAO by loading airports data from a JSON file and storing them in a map.
// If the map is not initialized, it creates a new one.
func (dao *MemoryAirportDAO) New() {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	var airports []*models.Airport

	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	jsonPath := filepath.Join(baseDir, "internal", "assets", "airports.json")

	b, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &airports)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}

	if dao.data == nil {
		dao.data = make(map[uuid.UUID]*models.Airport)
	}

	for _, airport := range airports {
		dao.data[airport.Id] = airport
	}
}

// FindAll retrieves all airports from the memory data store.
//
// The function iterates over the map of airports and appends each value to a new slice.
// It returns the slice containing all airports.
//
// Return:
// - []models.Airport: A slice of all airports in the memory data store.
func (dao *MemoryAirportDAO) FindAll() []*models.Airport {
	dao.mu.RLock()
	defer dao.mu.RUnlock()
	v := make([]*models.Airport, 0, len(dao.data))

	for _, value := range dao.data {
		v = append(v, value)
	}

	return v
}

// Insert adds a new airport to the memory data store.
//
// The function generates a new UUID for the airport, assigns it to the airport's Id field,
// and then inserts the airport into the map using the generated UUID as the key.
//
// Parameters:
// - t: A pointer to the Airport struct representing the airport to be inserted.
//
// Return:
// - This function does not return any value.
func (dao *MemoryAirportDAO) Insert(t *models.Airport) {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	id := uuid.New()

	t.Id = id

	dao.data[id] = t
}

// Update updates an existing airport in the memory data store.
//
// The function checks if an airport with the given ID exists in the map.
// If the airport exists, it updates the airport's data with the provided airport struct.
// If the airport does not exist, it returns an error indicating that the airport was not found.
//
// Parameters:
//   - t: A pointer to the Airport struct representing the updated airport data.
//     The Id field of the airport struct must be set to the ID of the airport to be updated.
//
// Return:
//   - error: An error indicating whether the update was successful or not.
//     If the airport was found and updated, the function returns nil.
//     If the airport was not found, the function returns an error with the message "not found".
func (dao *MemoryAirportDAO) Update(t *models.Airport) error {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	_, exists := dao.data[t.Id]

	if !exists {
		return errors.New("not found")
	}

	dao.data[t.Id] = t

	return nil
}

// Delete removes an airport from the memory data store based on the provided airport struct.
//
// The function checks if an airport with the given ID exists in the map.
// If the airport exists, it deletes the airport from the map using the airport's ID as the key.
// If the airport does not exist, the function does nothing.
//
// Parameters:
//   - t: A models.Airport struct representing the airport to be deleted.
//     The Id field of the airport struct must be set to the ID of the airport to be deleted.
//
// Return:
// - This function does not return any value.
func (dao *MemoryAirportDAO) Delete(t *models.Airport) {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	delete(dao.data, t.Id)
}

// FindById retrieves an airport from the memory data store based on the provided UUID.
//
// The function searches for an airport in the map using the provided UUID as the key.
// If the airport is found, it returns a pointer to the airport and a nil error.
// If the airport is not found, it returns nil and an error indicating that the airport was not found.
//
// Parameters:
// - id: A uuid.UUID representing the unique identifier of the airport to be retrieved.
//
// Return:
//   - *models.Airport: A pointer to the airport if found, or nil if not found.
//   - error: An error indicating whether the airport was found or not.
//     If the airport was found, the function returns nil.
//     If the airport was not found, the function returns an error with the message "not found".
func (dao *MemoryAirportDAO) FindById(id uuid.UUID) (*models.Airport, error) {
	dao.mu.RLock()
	defer dao.mu.RUnlock()
	airport, exists := dao.data[id]

	if !exists {
		return nil, errors.New("not found")
	}

	return airport, nil
}

// FindByName retrieves an airport from the memory data store based on the provided city name.
//
// The function iterates over the map of airports and checks if the city name of each airport matches the provided name.
// If a match is found, the function returns a pointer to the airport.
// If no match is found, the function returns nil.
//
// Parameters:
//   - name: A string representing the city name of the airport to be retrieved.
//
// Return:
//   - *models.Airport: A pointer to the airport if found, or nil if not found.
//     If the airport is found, the function returns a pointer to the airport.
//     If the airport is not found, the function returns nil.
func (dao *MemoryAirportDAO) FindByName(name string) *models.Airport {
	dao.mu.RLock()
	defer dao.mu.RUnlock()
	for _, value := range dao.data {
		if value.City.Name == name {
			return value
		}
	}
	return nil
}
