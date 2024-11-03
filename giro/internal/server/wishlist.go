package server

import (
	"encoding/json"
	"giro/internal/dao"
	"giro/internal/models"
	"log"
	"net/http"
)

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
