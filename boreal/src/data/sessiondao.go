package data

import (
	"errors"
	"passcom/boreal/src/models"
	"sync"

	"github.com/google/uuid"
)

type MemorySessionDAO struct {
	data map[uuid.UUID]*models.Session
	mu   sync.RWMutex
}

// New initializes the MemorySessionDAO by creating a new map to store sessions.
// It locks the mutex to ensure thread safety while creating the map.
func (dao *MemorySessionDAO) New() {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	dao.data = make(map[uuid.UUID]*models.Session)
}

// FindAll retrieves all sessions from the memory data store.
// It locks the read mutex to ensure thread safety while accessing the data.
//
// Returns:
//   - A slice of pointers to Session structs, representing all sessions in the data store.
//   - If no sessions are found, an empty slice is returned.
func (dao *MemorySessionDAO) FindAll() []*models.Session {
	dao.mu.RLock()
	defer dao.mu.RUnlock()

	v := make([]*models.Session, 0, len(dao.data))

	for _, value := range dao.data {
		v = append(v, value)
	}

	return v
}

// Insert adds a new session to the memory data store.
// It generates a new UUID for the session, sets it as the ID, initializes the reservations map,
// and then stores the session in the data map.
//
// Parameters:
//   - t: A pointer to a Session struct representing the session to be added.
func (dao *MemorySessionDAO) Insert(t *models.Session) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	id := uuid.New()
	t.ID = id
	t.Reservations = make(map[uuid.UUID]models.Reservation)
	t.Mu = sync.RWMutex{}
	t.FailedReservations = make(chan string)
	dao.data[id] = t
}

// Update updates an existing session in the memory data store.
// It locks the mutex to ensure thread safety while accessing the data.
//
// Parameters:
//   - t: A pointer to a Session struct representing the session to be updated.
//     The ID field of the session is used to identify the session to be updated.
//
// Returns:
//   - An error if the session with the given ID does not exist in the data store.
//   - nil if the session is successfully updated.
func (dao *MemorySessionDAO) Update(t *models.Session) error {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	_, exists := dao.data[t.ID]
	if !exists {
		return errors.New("not found")
	}

	dao.data[t.ID] = t
	return nil
}

// Delete removes a session from the memory data store based on the provided session object.
// It locks the mutex to ensure thread safety while accessing the data.
//
// Parameters:
//   - t: A pointer to a Session struct representing the session to be deleted.
//     The ID field of the session is used to identify the session to be deleted.
func (dao *MemorySessionDAO) Delete(t *models.Session) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	delete(dao.data, t.ID)
}

// FindById retrieves a session from the memory data store based on the provided ID.
// It locks the read mutex to ensure thread safety while accessing the data.
//
// Parameters:
//   - id: A UUID representing the ID of the session to be retrieved.
//
// Returns:
//   - A pointer to a Session struct representing the session with the given ID.
//     If the session is found, the function returns the session and nil as the error.
//   - If the session is not found, the function returns nil and an error with the message "not found".
func (dao *MemorySessionDAO) FindById(id uuid.UUID) (*models.Session, error) {
	dao.mu.RLock()
	defer dao.mu.RUnlock()

	session, exists := dao.data[id]
	if !exists {
		return nil, errors.New("not found")
	}

	return session, nil
}

// DeleteAll removes all sessions from the memory data store.
// It locks the mutex to ensure thread safety while accessing the data.
// After deleting all sessions, it initializes the data map with a new empty map.
func (dao *MemorySessionDAO) DeleteAll() {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	dao.data = make(map[uuid.UUID]*models.Session)
}
