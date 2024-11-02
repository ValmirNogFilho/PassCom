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

// MemoryClientDAO is a data access object (DAO) for managing client data in memory.
// It provides methods for inserting, updating, deleting, and retrieving clients.
type MemoryClientDAO struct {
	data map[uuid.UUID]*models.Client
	mu   sync.RWMutex
}

// New initializes the MemoryClientDAO by loading client data from a JSON file.
// It sets up the data map with the client data from the JSON file.
func (dao *MemoryClientDAO) New() {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	var clients []models.Client

	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	jsonPath := filepath.Join(baseDir, "internal", "assets", "airports.json")

	b, _ := os.ReadFile(jsonPath)

	json.Unmarshal(b, &clients)

	for _, client := range clients {
		dao.data[client.Id] = &client
	}
}

// FindAll retrieves all clients from the memory data store.
//
// It iterates over the data map and appends each client to a slice.
// The function returns a slice of pointers to the clients.
//
// The returned slice is created with a capacity equal to the length of the data map.
// This ensures that the slice can accommodate all clients without resizing.
//
// If no clients are found, an empty slice is returned.
func (dao *MemoryClientDAO) FindAll() []*models.Client {
	dao.mu.RLock()
	defer dao.mu.RUnlock()
	v := make([]*models.Client, 0, len(dao.data))

	for _, value := range dao.data {
		v = append(v, value)
	}

	return v
}

// Insert adds a new client to the memory data store.
//
// The function generates a new UUID for the client and assigns it to the client's Id field.
// Then, it inserts the client into the data map using the generated UUID as the key.
//
// Parameters:
//   - t: A pointer to the client model to be inserted. The client's Id field will be updated with a new UUID.
func (dao *MemoryClientDAO) Insert(t *models.Client) {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	id := uuid.New()

	t.Id = id

	dao.data[id] = t
}

// Update updates an existing client in the memory data store.
//
// The function checks if a client with the given Id exists in the data map.
// If the client is found, the function updates the client's data in the data map.
// If the client is not found, the function returns an error indicating that the client was not found.
//
// Parameters:
//   - t: A pointer to the client model to be updated. The client's Id field should be set to the desired client's UUID.
//
// Return:
//   - An error if the client was not found in the data map.
func (dao *MemoryClientDAO) Update(t *models.Client) error {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	lastClient, exists := dao.data[t.Id]

	if !exists {
		return errors.New("not found")
	}

	dao.data[t.Id] = lastClient

	return nil
}

// Delete removes a client from the memory data store based on the provided client model.
//
// The function checks if a client with the given Id exists in the data map.
// If the client is found, the function deletes the client from the data map.
// If the client is not found, the function does nothing.
//
// Parameters:
//   - t: The client model to be deleted. The function uses the client's Id field to identify the client in the data map.
func (dao *MemoryClientDAO) Delete(t models.Client) {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	delete(dao.data, t.Id)
}

// FindById retrieves a client from the memory data store based on the provided UUID.
//
// The function checks if a client with the given UUID exists in the data map.
// If the client is found, the function returns a pointer to the client and a nil error.
// If the client is not found, the function returns nil and an error indicating that the client was not found.
//
// Parameters:
//   - id: The UUID of the client to be retrieved.
//
// Return:
//   - A pointer to the client if found, nil otherwise.
//   - An error indicating that the client was not found, nil otherwise.
func (dao *MemoryClientDAO) FindById(id uuid.UUID) (*models.Client, error) {
	dao.mu.RLock()
	defer dao.mu.RUnlock()
	client, exists := dao.data[id]

	if !exists {
		return nil, errors.New("not found")
	}

	return client, nil
}
