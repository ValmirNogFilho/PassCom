package models

type Request struct {
	Action string      `json:"Action"`
	Auth   string      `json:"Auth"` // Modifique para aceitar nulos
	Data   interface{} `json:"Data"`
}
