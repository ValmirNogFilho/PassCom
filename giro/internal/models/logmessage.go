package models

import "time"

type Status int

const (
	PENDING Status = iota
	COMMITED
	REJECTED
)

type LogMessage struct {
	Timestamp time.Time   `json:"timestamp"`
	Status    Status      `json:"status"`
	Data      interface{} `json:"data"`
}
