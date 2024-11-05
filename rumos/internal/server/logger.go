package server

import (
	"rumos/internal/models"
	"time"
)

func (s *System) AddMessageToLog(timestamp time.Time, message models.Message, status models.Status) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	log := models.LogMessage{
		Timestamp: timestamp,
		Status:    status,
		Type:      models.MESSAGE,
		Data:      message,
	}

	s.Log = append(s.Log, log)
	if len(s.Log) > 1000 {
		s.Log = s.Log[1:]
	}
}

func (s *System) AddTransactionToLog(timestamp time.Time, transaction models.Transaction, status models.Status) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	log := models.LogMessage{
		Timestamp: timestamp,
		Status:    status,
		Type:      models.TRANSACTION,
		Data:      transaction,
	}

	s.Log = append(s.Log, log)
	if len(s.Log) > 1000 {
		s.Log = s.Log[1:]
	}
}
