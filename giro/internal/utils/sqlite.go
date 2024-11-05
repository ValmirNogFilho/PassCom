package utils

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// type supportedTypes interface {
// 	*models.Airport | *models.Client |
// 		map[uuid.UUID]*models.Flight | *models.Session
// }

func OpenDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CloseDb(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Erro ao obter a conexão do banco:", err)
	}

	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Fatal("Erro ao fechar a conexão com o banco:", err)
		}
	}()
}
