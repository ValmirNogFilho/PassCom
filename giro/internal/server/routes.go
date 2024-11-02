package server

import (
	"encoding/json"
	"fmt"
	"giro/internal/dao"
	"giro/internal/models"
	"net"
	"net/http"
)

func GetAirports(request models.Request) models.Response {
	_, exists := SessionIfExists(request.Auth)

	if !exists {
		return models.Response{
			Error: "not authorized",
		}
	}
	responseData := make([]map[string]interface{}, 0)

	airports := dao.GetAirportDAO().FindAll()

	for _, airport := range airports {

		airportresponse := make(map[string]interface{})
		airportresponse["Name"] = airport.Name
		airportresponse["City"] = airport.City
		responseData = append(responseData, airportresponse)
	}

	return models.Response{
		Data: map[string]interface{}{
			"Airports": responseData,
		},
		Status: http.StatusOK,
	}
}

// AllRoutes handles the retrieval of all available routes.
// It checks if the provided authentication token is valid and returns a list of all routes if authorized.
//
// Parameters:
//   - auth: A string representing the authentication token provided by the client.
//   - conn: A net.Conn object representing the connection to the client.
func AllRoutes(auth string, conn net.Conn) models.Response {

	_, exists := SessionIfExists(auth)

	if !exists {
		return models.Response{
			Error:  "not authorized",
			Status: http.StatusUnauthorized,
		}
	}

	dao := dao.GetFlightDAO()
	dao.New()

	return models.Response{
		Data: map[string]interface{}{
			"all-routes": dao.FindAll(),
		},
		Status: http.StatusOK,
	}

}

// Route handles the retrieval of a route between two cities.
// It checks if the provided authentication token is valid and returns a route if authorized.
// If the source or destination city is not found, it returns an error response.
// If no route is found between the source and destination cities, it returns an error response.
//
// Parameters:
//   - auth: A string representing the authentication token provided by the client.
//   - data: An interface containing the source and destination city names.
//   - conn: A net.Conn object representing the connection to the client.

// func Route(auth string, data interface{}, conn net.Conn) {
// 	_, exists := SessionIfExists(auth)

// 	if !exists {
// 		WriteNewResponse(models.Response{
// 			Error: "not authorized",
// 		}, conn)
// 		return
// 	}

// 	var routeRequest models.RouteRequest
// 	var response models.Response

// 	jsonData, _ := json.Marshal(data)
// 	json.Unmarshal(jsonData, &routeRequest)

// 	src := dao.GetAirportDAO().FindByName(routeRequest.Source)
// 	dest := dao.GetAirportDAO().FindByName(routeRequest.Dest)

// 	if src == nil || dest == nil {
// 		WriteNewResponse(models.Response{
// 			Error: "not valid city name",
// 		}, conn)
// 		return
// 	}

// 	// route, err := dao.GetFlightDAO().FindBySourceAndDest(src.Id, dest.Id)
// 	path, path_err := dao.GetFlightDAO().BreadthFirstSearch(src.ID, dest.ID)
// 	if path_err != nil {
// 		response.Error = "no route"
// 	} else {
// 		cities_path := make([]models.Route, len(path))
// 		for i, flight := range path {
// 			cities_path[i].Path = make([]models.City, 2)
// 			cities_path[i].FlightId = flight.Id
// 			srcById, _ := dao.GetAirportDAO().FindById(flight.SourceAirportId)
// 			cities_path[i].Path[0] = srcById.City
// 			destById, _ := dao.GetAirportDAO().FindById(flight.DestAirportId)
// 			cities_path[i].Path[1] = destById.City
// 		}

// 		response.Data = map[string]interface{}{
// 			"path": cities_path,
// 		}

// 	}

// 	WriteNewResponse(response, conn)
// }

// Flights handles the retrieval of flight details based on provided flight IDs.
// It checks if the provided authentication token is valid and returns flight details if authorized.
// If any of the provided flight IDs does not exist, it returns an error response.
//
// Parameters:
//   - auth: A string representing the authentication token provided by the client.
//   - data: An interface containing the flight IDs.
//   - conn: A net.Conn object representing the connection to the client.
//
// Return:
//   - This function does not return any value. It writes a response to the client's connection.
//   - The response contains flight details if authorized and valid flight IDs are provided.
//   - If not authorized, it returns an error response with the message "not authorized".
//   - If any of the provided flight IDs does not exist, it returns an error response.
func Flights(request models.Request) models.Response {
	_, exists := SessionIfExists(request.Auth)
	if !exists {
		return models.Response{
			Error: "not authorized",
		}
	}

	var flightsRequest models.FlightsRequest

	jsonData, _ := json.Marshal(request.Data)
	json.Unmarshal(jsonData, &flightsRequest)

	responseData, err := getRoute(flightsRequest.FlightIds)
	if err != nil {
		return models.Response{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		}
	}

	return models.Response{
		Data: map[string]interface{}{
			"Flights": responseData,
		},
		Status: http.StatusOK,
	}
}

// getRoute retrieves flight details for a given list of flight IDs.
// It fetches the flight details from the database and constructs a response containing the flight details.
//
// Parameters:
//   - flightIds: A slice of uuid.UUID representing the flight IDs for which the details need to be retrieved.
//
// Return:
//   - A slice of map[string]interface{} containing the flight details. Each map represents a flight and contains the following keys:
//   - "Seats": An integer representing the number of available seats on the flight.
//   - "Src": A string representing the source city of the flight.
//   - "Dest": A string representing the destination city of the flight.
//   - An error if any of the provided flight IDs does not exist in the database.
func getRoute(flightIds []uint) ([]map[string]interface{}, error) {
	responseData := make([]map[string]interface{}, len(flightIds))
	for i, id := range flightIds {
		flightresponse := make(map[string]interface{})
		flight, err := dao.GetFlightDAO().FindById(id)
		if err != nil {
			return nil, fmt.Errorf("some flight doesn't exist: %s", id)
		}
		fmt.Println(flight.OriginAirport)

		flightresponse["Seats"] = flight.Seats
		flightresponse["Src"] = flight.OriginAirport.City.CityName
		flightresponse["Dest"] = flight.DestinationAirport.City.CityName
		responseData[i] = flightresponse
	}
	return responseData, nil
}