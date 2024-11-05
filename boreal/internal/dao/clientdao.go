package dao

import (
	"boreal/internal/models"
	"boreal/internal/utils"
	"log"
)

// DBClientDAO is a data access object (DAO) for managing client data in memory.
// It provides methods for inserting, updating, deleting, and retrieving clients.
type DBClientDAO struct{}

// New initializes the DBClientDAO by loading client data from a JSON file.
// It sets up the data map with the client data from the JSON file.
func (dao *DBClientDAO) New() {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.Client{})
	defer utils.CloseDb(db)
}

// FindAll retrieves all clients from the memory data store.
//
// It iterates over the data map and appends each client to a slice.
// The function returns a slice of pointers to the clients.
//
// The returned slice is created with a capacity equal to the length of the data map.
// This ensures that the slice can accommodate all clients without resizing.
//
// If no clients are found, an empty slice is returned.
func (dao *DBClientDAO) FindAll() []models.Client {

	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var clients []models.Client = make([]models.Client, 0)

	db.Find(&clients)
	return clients
}

func (dao *DBClientDAO) Insert(client models.Client) {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	db.Create(&client)
}

// Update updates an existing client in the memory data store.
//
// The function checks if a client with the given Id exists in the data map.
// If the client is found, the function updates the client's data in the data map.
// If the client is not found, the function returns an error indicating that the client was not found.
//
// Parameters:
//   - t: A pointer to the client model to be updated. The client's Id field should be set to the desired client's UUID.
//
// Return:
//   - An error if the client was not found in the data map.
func (dao *DBClientDAO) Update(c models.Client) error {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var client models.Client
	if err := db.First(&client, "id = ?", c.ID).Error; err != nil {
		log.Println("Client not found:", err)
		return err
	}

	client = c
	if err := db.Save(&client).Error; err != nil {
		log.Println("Client not updated:", err)
		return err
	}
	log.Println("Client updated:", client)
	return nil

}

// Delete removes a client from the memory data store based on the provided client model.
//
// The function checks if a client with the given Id exists in the data map.
// If the client is found, the function deletes the client from the data map.
// If the client is not found, the function does nothing.
//
// Parameters:
//   - t: The client model to be deleted. The function uses the client's Id field to identify the client in the data map.
func (dao *DBClientDAO) Delete(client models.Client) {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	db.Delete(client)
}

// FindById retrieves a client from the memory data store based on the provided UUID.
//
// The function checks if a client with the given UUID exists in the data map.
// If the client is found, the function returns a pointer to the client and a nil error.
// If the client is not found, the function returns nil and an error indicating that the client was not found.
//
// Parameters:
//   - id: The UUID of the client to be retrieved.
//
// Return:
//   - A pointer to the client if found, nil otherwise.
//   - An error indicating that the client was not found, nil otherwise.
func (dao *DBClientDAO) FindById(id uint) (*models.Client, error) {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var client models.Client
	if err := db.
		Preload("ClientFlights").
		Preload("ClientFlights.Flight").
		Preload("ClientFlights.Flight.OriginAirport").
		Preload("ClientFlights.Flight.DestinationAirport").Take(&client, "id = ?", id).Error; err != nil {
		log.Println("Error searching client:", err)
		return nil, err
	}
	log.Println("Client found:", client)
	return &client, nil
}

func (dao *DBClientDAO) FindByUsername(username string) (*models.Client, error) {
	db, err := utils.OpenDb()

	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseDb(db)

	var client models.Client
	if err := db.Where(&models.Client{
		Username: username,
	}).Take(&client).Error; err != nil {
		log.Println("Error searching client:", err)
		return nil, err
	}
	log.Println("Client found:", client)
	return &client, nil
}
