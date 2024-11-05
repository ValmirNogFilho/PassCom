package server

import (
	"encoding/json"
	"giro/internal/dao"
	"giro/internal/models"
	"log"
	"net/http"
	"strconv"
)

func handleWishlist(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	switch r.Method {
	case http.MethodGet:
		handleGetWishlist(w, r)
	case http.MethodPost:
		handleAddToWishlist(w, r)
	case http.MethodDelete:
		handleRemoveFromWishlist(w, r)
	default:
		http.Error(w, "only GET, POST, DELETE allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleRemoveFromWishlist(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	token := r.Header.Get("Authorization")

	queryParams := r.URL.Query()

	id := queryParams.Get("id")

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		returnResponse(w, r, models.Response{
			Status: http.StatusBadRequest,
		})
		return
	}

	response := DeleteFromWishlist(uint(idUint),
		models.Request{
			Auth: token,
		},
	)

	returnResponse(w, r, response)
}

func handleAddToWishlist(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	token := r.Header.Get("Authorization")

	var addWish models.WishlistOperation

	err := json.NewDecoder(r.Body).Decode(&addWish)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := AddToWishlist(
		models.Request{
			Auth: token,
			Data: addWish,
		},
	)
	returnResponse(w, r, response)
}

func handleGetWishlist(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	response := GetWishlist(
		models.Request{
			Auth: token,
		},
	)
	returnResponse(w, r, response)
}

func GetWishlist(req models.Request) models.Response {
	session, exists := SessionIfExists(req.Auth)
	if !exists {
		return models.Response{
			Error:  "not authorized",
			Status: http.StatusUnauthorized,
		}
	}

	return models.Response{
		Data: map[string]interface{}{
			"Wishes": session.Wishlist,
		},
		Status: http.StatusOK,
	}
}

func DeleteFromWishlist(id uint, req models.Request) models.Response {
	session, exists := SessionIfExists(req.Auth)
	if !exists {
		return models.Response{
			Error:  "not authorized",
			Status: http.StatusUnauthorized,
		}
	}

	for i, w := range session.Wishlist {
		if w.ID == id {
			// Remover preservando a ordem
			session.Wishlist = append(session.Wishlist[:i], session.Wishlist[i+1:]...)
			break
		}
	}

	if err := dao.GetSessionDAO().Update(session); err != nil {
		log.Printf("Failed to update session: %v", err)
		return models.Response{
			Error:  "failed to update session",
			Status: http.StatusInternalServerError,
		}
	}

	return models.Response{
		Data: map[string]interface{}{
			"msg": "wish deleted",
		},
		Status: http.StatusOK,
	}
}

func AddToWishlist(req models.Request) models.Response {
	session, exists := SessionIfExists(req.Auth)
	if !exists {
		return models.Response{
			Error:  "not authorized",
			Status: http.StatusUnauthorized,
		}
	}

	var addWish models.WishlistOperation

	jsonData, _ := json.Marshal(req.Data)
	json.Unmarshal(jsonData, &addWish)

	flight, _ := dao.GetFlightDAO().FindById(addWish.FlightId)

	session.Wishlist = append(session.Wishlist, *flight)

	dao.GetSessionDAO().Update(session)

	return models.Response{
		Data: map[string]interface{}{
			"msg": "wish added",
		},
		Status: http.StatusOK,
	}
}
