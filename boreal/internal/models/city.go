package models

type City struct {
	Name      string  `json:"Name"`
	State     string  `json:"State"`
	Country   string  `json:"Country"`
	Latitude  float32 `json:"Latitude"`
	Longitude float32 `json:"Longitude"`
}
