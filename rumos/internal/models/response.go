package models

type Response struct {
	Error  string `json:"error"`
	Data   map[string]interface{}
	Status int `json:"status"`
}
