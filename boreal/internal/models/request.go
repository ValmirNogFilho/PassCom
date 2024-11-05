package models

type Request struct {
	Auth string      `json:"Auth"`
	Data interface{} `json:"Data"`
}
 