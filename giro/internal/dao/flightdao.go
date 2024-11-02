package dao

import (
	"giro/internal/models"
	"giro/internal/utils"
	"log"
)

type DBFlightDAO struct {
}

func (dao *DBFlightDAO) New() {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)
	db.AutoMigrate(&models.Ticket{})
}

func (dao *DBFlightDAO) FindAll() []models.Flight {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var flights []models.Flight = make([]models.Flight, 0)

	db.Find(&flights)

	return flights
}

func (dao *DBFlightDAO) Insert(flight models.Flight) {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	db.Create(&flight)

}

func (dao *DBFlightDAO) Update(f models.Flight) error {

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var flight models.Flight
	if err := db.First(&flight, "id = ?", f.ID).Error; err != nil {
		log.Println("Flight not found:", err)
		return err
	}

	flight = f
	if err := db.Save(&flight).Error; err != nil {
		log.Println("Flight not updated:", err)
		return err
	}
	log.Println("Flight updated:", flight)
	return nil
}

func (dao *DBFlightDAO) Delete(a models.Flight) error {

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	if err := db.Delete(&models.Flight{}, "id = ?", a.ID).Error; err != nil {
		log.Println("Error deleting Flight:", err)
		return err
	}

	log.Println("Flight successfully deleted.")
	return nil

}

func (dao *DBFlightDAO) FindById(id uint) (*models.Flight, error) {

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var flight models.Flight

	if err := db.Preload("OriginAirport").
		Preload("DestinationAirport").
		Preload("Tickets").First(&flight, "id=?", id).Error; err != nil {
		log.Println("Error searching flight:", err)
		return nil, err
	}
	log.Println("Flight found:", flight)
	return &flight, nil
}

func (dao *DBFlightDAO) FindBySource(id uint) ([]models.Flight, error) {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var flights []models.Flight = make([]models.Flight, 0)
	if err := db.Preload("OriginAirport").
		Preload("DestinationAirport").
		Preload("Tickets").Where(&models.Flight{
		OriginAirportID: id,
	}).Find(&flights).Error; err != nil {
		log.Println("Error searching flights:", err)
		return nil, err
	}
	log.Println("Flights found:", flights)
	return flights, nil
}

func (dao *DBFlightDAO) FindBySourceAndDest(source uint, dest uint) ([]models.Flight, error) {

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var flights []models.Flight = make([]models.Flight, 0)
	if err := db.Preload("OriginAirport").
		Preload("DestinationAirport").
		Preload("Tickets").
		Where(&models.Flight{
			OriginAirportID:      source,
			DestinationAirportID: dest,
		}).Find(&flights).Error; err != nil {
		log.Println("Error searching flight:", err)
		return nil, err
	}
	log.Println("Flights found:", flights)
	return flights, nil
}

func (dao *DBFlightDAO) DeleteAll() {

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	db.Unscoped().Where("1=1").Delete(&models.Flight{})
}
