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

func (dao *DBFlightDAO) FindBySourceAndDest(source uint, dest uint) (*models.Flight, error) {

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var flight models.Flight
	if err := db.Preload("OriginAirport").
		Preload("DestinationAirport").
		Preload("Tickets").
		Where(&models.Flight{
			OriginAirportID:      source,
			DestinationAirportID: dest,
		}).Find(&flight).Error; err != nil {
		log.Println("Error searching flight:", err)
		return nil, err
	}
	log.Println("Flights found:", flight)
	return &flight, nil
}

// func (dao *DBFlightDAO) BreadthFirstSearch(source uuid.UUID, dest uuid.UUID) ([]*models.Flight, error) {
// 	dao.mu.RLock()
// 	defer dao.mu.RUnlock()
// 	visited := make(map[uuid.UUID]bool, len(dao.data))
// 	queue := []uuid.UUID{source}
// 	visited[source] = true
// 	parent := make(map[uuid.UUID]uuid.UUID, len(dao.data))
// 	parent[source] = source

// 	for len(queue) > 0 {
// 		current := queue[0]
// 		queue = queue[1:]
// 		if current == dest {
// 			break
// 		}

// 		for neighbor, flight := range dao.data[current] {
// 			if !visited[neighbor] && flight.Seats > 0 {
// 				visited[neighbor] = true
// 				queue = append(queue, neighbor)
// 				parent[neighbor] = current
// 			}
// 		}
// 	}

// 	path := []*models.Flight{}
// 	current := dest

// 	if !visited[dest] {
// 		return nil, errors.New("no route available")
// 	}

// 	for current != source {
// 		prev := parent[current]
// 		flight := dao.data[prev][current]
// 		path = append([]*models.Flight{flight}, path...)
// 		current = prev
// 	}

// 	return path, nil
// }

func (dao *DBFlightDAO) DeleteAll() {

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	db.Unscoped().Where("1=1").Delete(&models.Flight{})
}
