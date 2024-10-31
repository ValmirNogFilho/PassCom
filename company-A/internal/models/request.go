package models

type Request struct {
	Auth string      `json:"Auth"` // Modifique para aceitar nulos
	Data interface{} `json:"Data"`
}
