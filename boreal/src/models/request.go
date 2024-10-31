package models

type Request struct {
	Action string      `json:"Action"`
	Auth   *string     `json:"Auth"`
	Data   interface{} `json:"Data"`
}
