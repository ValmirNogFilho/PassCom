package server

import (
	"fmt"
	"rumos/internal/dao"
	"rumos/internal/models"
	"time"

	"github.com/google/uuid"
)

// CleanupSessions periodically checks for inactive sessions and reservations, and cleans them up.
// It runs every minute and checks each session and its reservations against the provided timeout.
// If a session or a reservation is inactive (i.e., its last activity time is older than the timeout),
// it is deleted from the system.
//
// Parameters:
//   - timeout: The duration after which a session or a reservation is considered inactive.
func CleanupSessions(timeout time.Duration) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		for _, session := range dao.GetSessionDAO().FindAll() {
			if time.Since(session.LastTimeActive) > timeout {
				fmt.Printf("Encerrando sessão %s por inatividade\n", session.ID)
				dao.GetSessionDAO().Delete(session)
			}
		}
	}
}

// SessionIfExists checks if a session exists for the given token.
// If a session is found, it updates the session's last activity time and returns the session along with a boolean value of true.
// If no session is found or an error occurs during the process, it returns nil and false.
//
// Parameters:
//   - token: A string representing the session token to be checked.
//
// Return:
//   - *models.Session: A pointer to the found session if it exists, or nil if no session is found or an error occurs.
//   - bool: A boolean value indicating whether a session was found (true) or not (false).
func SessionIfExists(token string) (*models.Session, bool) {
	uuid, err := uuid.Parse(token)
	if err != nil {
		return nil, false
	}
	session, err := dao.GetSessionDAO().FindById(uuid)
	if err != nil {
		return nil, false
	}
	session.LastTimeActive = time.Now()
	dao.GetSessionDAO().Update(session)
	return session, true
}
