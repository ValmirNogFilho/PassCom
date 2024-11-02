package models

type City struct {
	CityName  string `gorm:"size:80"`
	State     string `gorm:"size:40"`
	Country   string `gorm:"size:40"`
	Latitude  float32
	Longitude float32
}