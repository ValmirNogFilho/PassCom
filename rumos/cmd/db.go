package main

import (
	"log"
	"vendepass/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connectDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	db, err := connectDatabase()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	db.AutoMigrate(&models.Airport{})

	log.Println("Migração concluída e conexão estabelecida.")
}
