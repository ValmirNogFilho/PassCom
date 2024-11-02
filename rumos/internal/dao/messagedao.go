package dao

import (
	"log"
	"rumos/internal/models"
	"rumos/internal/utils"
)

type DBMessageDAO struct {
}

func (dao *DBMessageDAO) New() {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)
	db.AutoMigrate(&models.Message{})
}

func (dao *DBMessageDAO) FindAll() []models.Message {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var messages []models.Message = make([]models.Message, 0)

	db.Find(&messages)

	return messages
}

func (dao *DBMessageDAO) Insert(message models.Message) {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	db.Create(&message)

}

func (dao *DBMessageDAO) Update(m models.Message) error {

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var message models.Message
	if err := db.First(&message, "id = ?", m.ID).Error; err != nil {
		log.Println("Message not found:", err)
		return err
	}

	message = m
	if err := db.Save(&message).Error; err != nil {
		log.Println("Message not updated:", err)
		return err
	}
	log.Println("Message updated:", message)
	return nil
}

func (dao *DBMessageDAO) Delete(a models.Flight) error {

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	if err := db.Delete(&models.Airport{}, "id = ?", a.ID).Error; err != nil {
		log.Println("Error deleting Flight:", err)
		return err
	}

	log.Println("Flight successfully deleted.")
	return nil

}

func (dao *DBMessageDAO) FindById(id uint) (*models.Message, error) {

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var message models.Message

	if err := db.First(&message, "id=?", id).Error; err != nil {
		log.Println("Error searching message:", err)
		return nil, err
	}
	log.Println("Message found:", message)
	return &message, nil
}
