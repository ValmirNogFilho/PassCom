package models

import "time"

type Status int
type LogType string

const (
	PENDING Status = iota
	COMMITED
	REJECTED
)

const (
	TRANSACTION LogType = "transaction"
	MESSAGE     LogType = "message"
)

type LogMessage struct {
	Timestamp time.Time   `json:"timestamp"`
	Type      LogType     `json:"type"`
	Status    Status      `json:"status"`
	Data      interface{} `json:"data"`
}
