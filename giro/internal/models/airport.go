package models

import (
	"gorm.io/gorm"
)

type Airport struct {
	gorm.Model
	Name string `gorm:"size:80"`
	City City   `gorm:"embedded;embeddedPrefix:city_"`
}
