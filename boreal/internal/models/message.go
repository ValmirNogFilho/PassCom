package models

type Message struct {
	Id          string         `json:"Id"`
	From        string         `json:"From"`
	To          string         `json:"To"`
	VectorClock map[string]int `json:"VectorClock"`
	Request
}
