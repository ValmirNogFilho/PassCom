package dao

import (
	"errors"
	"log"
	"rumos/internal/models"
	"rumos/internal/utils"
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

func (dao *DBFlightDAO) FindPathBFS(source uint, dest uint) ([]models.Flight, error) {
	db, err := utils.OpenDb()
	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var flights []models.Flight
	if err := db.Preload("OriginAirport").
		Preload("DestinationAirport").
		Find(&flights).Error; err != nil {
		log.Println("Error loading flights:", err)
		return nil, err
	}

	graph := make(map[uint][]models.Flight)
	for _, flight := range flights {
		graph[flight.OriginAirportID] = append(graph[flight.OriginAirportID], flight)
	}

	type Node struct {
		AirportID uint
		Path      []models.Flight
	}

	queue := []Node{{AirportID: source, Path: []models.Flight{}}}
	visited := make(map[uint]bool)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		visited[current.AirportID] = true

		if current.AirportID == dest {
			return current.Path, nil
		}

		for _, flight := range graph[current.AirportID] {
			if visited[flight.DestinationAirportID] || flight.Seats <= 0 {
				continue
			}

			newPath := append([]models.Flight{}, current.Path...)
			newPath = append(newPath, flight)

			queue = append(queue, Node{
				AirportID: flight.DestinationAirportID,
				Path:      newPath,
			})
		}
	}

	log.Println("No path found from source to destination")
	return nil, errors.New("no path found from source to destination")
}

func (dao *DBFlightDAO) FindByCompany(company string) ([]models.Flight, error) {
	db, err := utils.OpenDb()
	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var flights []models.Flight
	if err := db.Preload("OriginAirport").
		Preload("DestinationAirport").
		Preload("Tickets").
		Where("company = ?", company).
		Find(&flights).Error; err != nil {
		log.Println("Error finding flights by company:", err)
		return nil, err
	}

	log.Printf("Flights found for company %s: %v", company, flights)
	return flights, nil
}

func (dao *DBFlightDAO) FindByUniqueId(uniqueId string) (*models.Flight, error) {
	db, err := utils.OpenDb()
	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var flight models.Flight

	// Busca o voo usando o campo `UniqueId`
	if err := db.Preload("OriginAirport").
		Preload("DestinationAirport").
		Preload("Tickets").
		Where("unique_id = ?", uniqueId).
		First(&flight).Error; err != nil {
		log.Println("Error searching flight by unique ID:", err)
		return nil, err
	}

	log.Println("Flight found by unique ID:", flight)
	return &flight, nil
}

func (dao *DBFlightDAO) DeleteByUniqueId(uniqueId string) error {
	db, err := utils.OpenDb()
	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	// Exclui o voo com o `UniqueId` especificado
	if err := db.Where("unique_id = ?", uniqueId).Delete(&models.Flight{}).Error; err != nil {
		log.Println("Error deleting flight by unique ID:", err)
		return err
	}

	log.Printf("Flight with unique ID %s successfully deleted.", uniqueId)
	return nil
}

func (dao *DBFlightDAO) DeleteByCompany(company string) error {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	if err := db.Where("company =?", company).Delete(&models.Flight{}).Error; err != nil {
		log.Println("Error deleting flights by company:", err)
		return err
	}

	log.Printf("Flights deleted for company %s", company)
	return nil
}

func (dao *DBFlightDAO) DeleteAll() {

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	db.Unscoped().Where("1=1").Delete(&models.Flight{})
}
