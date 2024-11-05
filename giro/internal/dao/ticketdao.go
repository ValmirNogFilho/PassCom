package dao

import (
	"giro/internal/models"
	"giro/internal/utils"
	"log"

	"github.com/google/uuid"
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

	// Gera um UniqueId se não estiver presente
	if ticket.UniqueId == "" {
		uniqueId, err := uuid.NewV7()
		if err != nil {
			log.Fatal("Error generating unique ID:", err)
		}
		ticket.UniqueId = uniqueId.String()
	}

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

// FindByUniqueId busca um ticket pelo UniqueId.
func (dao *DBTicketDAO) FindByUniqueId(uniqueId string) (*models.Ticket, error) {
	db, err := utils.OpenDb()
	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var ticket models.Ticket
	if err := db.Preload("Flight").
		Where("unique_id = ?", uniqueId).
		First(&ticket).Error; err != nil {
		log.Println("Error finding ticket by unique ID:", err)
		return nil, err
	}

	log.Println("Ticket found by unique ID:", ticket)
	return &ticket, nil
}

// DeleteByUniqueId remove um ticket pelo UniqueId.
func (dao *DBTicketDAO) DeleteByUniqueId(uniqueId string) error {
	db, err := utils.OpenDb()
	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	if err := db.Where("unique_id = ?", uniqueId).Delete(&models.Ticket{}).Error; err != nil {
		log.Println("Error deleting ticket by unique ID:", err)
		return err
	}

	log.Printf("Ticket with unique ID %s successfully deleted.", uniqueId)
	return nil
}
