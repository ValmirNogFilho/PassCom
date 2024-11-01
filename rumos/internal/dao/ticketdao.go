package dao

import (
	"log"
	"vendepass/internal/models"
	"vendepass/internal/utils"
)

type DBTicketDAO struct{}

func (dao DBTicketDAO) New() {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)
	db.AutoMigrate(&models.Ticket{})
}

func (dao *DBTicketDAO) FindAll() []models.Ticket {
	var tickets []models.Ticket = make([]models.Ticket, 0)

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	db.Find(&tickets)

	return tickets
}

func (dao *DBTicketDAO) Insert(ticket models.Ticket) {

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	if err := db.Create(&ticket).Error; err != nil {
		log.Println("Error inserting ticket:", err)
	} else {
		log.Println("Ticket successfully inserted:", ticket)
	}
}
func (dao *DBTicketDAO) Update(a models.Ticket) error {

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var ticket models.Ticket
	if err := db.First(&ticket, "id = ?", a.ID).Error; err != nil {
		log.Println("Ticket not found:", err)
		return err
	}

	ticket = a
	if err := db.Save(&ticket).Error; err != nil {
		log.Println("Ticket not updated:", err)
		return err
	}
	log.Println("Ticket updated:", ticket)
	return nil
}

func (dao *DBTicketDAO) Delete(a models.Ticket) {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	if err := db.Delete(&models.Ticket{}, "id = ?", a.ID).Error; err != nil {
		log.Println("Error deleting ticket:", err)
	} else {
		log.Println("Ticket successfully deleted.")
	}
}

func (dao *DBTicketDAO) FindById(id uint) (*models.Ticket, error) {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var ticket models.Ticket
	if err := db.Preload("Flight").Take(&ticket, "id = ?", id).Error; err != nil {
		log.Println("Error searching ticket:", err)
		return nil, err
	}
	log.Println("Ticket found:", ticket)
	return &ticket, nil
}
