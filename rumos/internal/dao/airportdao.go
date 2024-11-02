// Package dao implements the Data Access Object for the database of the server.
package dao

import (
	"log"
	"rumos/internal/models"
	"rumos/internal/utils"
)

type DBAirportDAO struct{}

func (dao DBAirportDAO) New() {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.Airport{})
	utils.CloseDb(db)
}

func (dao *DBAirportDAO) FindAll() []models.Airport {
	var airports []models.Airport = make([]models.Airport, 0)

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	db.Find(&airports)

	return airports
}

func (dao *DBAirportDAO) Insert(airport models.Airport) {

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	if err := db.Create(&airport).Error; err != nil {
		log.Println("Error inserting airport:", err)
	} else {
		log.Println("Airport successfully inserted:", airport)
	}
}
func (dao *DBAirportDAO) Update(a models.Airport) error {

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var airport models.Airport
	if err := db.First(&airport, "id = ?", a.ID).Error; err != nil {
		log.Println("Airport not found:", err)
		return err
	}

	airport = a
	if err := db.Save(&airport).Error; err != nil {
		log.Println("Aiport not updated:", err)
		return err
	}
	log.Println("Airport updated:", airport)
	return nil
}

func (dao *DBAirportDAO) Delete(a models.Airport) {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	if err := db.Delete(&models.Airport{}, "id = ?", a.ID).Error; err != nil {
		log.Println("Error deleting airport:", err)
	} else {
		log.Println("Airport successfully deleted.")
	}
}

func (dao *DBAirportDAO) FindById(id uint) (*models.Airport, error) {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var airport models.Airport
	if err := db.Take(&airport, "id = ?", id).Error; err != nil {
		log.Println("Error searching airport:", err)
		return nil, err
	}
	log.Println("Airport found:", airport)
	return &airport, nil
}

func (dao *DBAirportDAO) FindByName(name string) *models.Airport {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var airport models.Airport
	if err := db.First(&airport, "name = ?", name).Error; err != nil {
		log.Println("Error searching airport:", err)
		return nil
	}
	log.Println("Airport found:", airport)
	return &airport
}
