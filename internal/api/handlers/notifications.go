// AnhCao 2024
package handlers

import (
	"net/http"

	"go.uber.org/zap"

	title "github.com/AnhCaooo/electric-notifications/internal/constants"
	"github.com/AnhCaooo/electric-notifications/internal/models"
	"github.com/AnhCaooo/go-goods/encode"
)

// create token to Mongo or update time live if it exists
func (h Handler) CreateToken(w http.ResponseWriter, r *http.Request) {
	reqBody, err := encode.DecodeRequest[models.NotificationToken](r)
	if err != nil {
		h.logger.Error(title.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the token into the database
	err = h.mongo.InsertToken(reqBody)
	if err != nil {
		// You should have better error handling here
		h.logger.Error(title.Server, zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) SendNotifications(w http.ResponseWriter, r *http.Request) {
	reqBody, err := encode.DecodeRequest[models.NotificationMessage](r)
	if err != nil {
		h.logger.Error(title.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// retrieve all associated device tokens with given userId
	tokens, err := h.mongo.GetTokens(reqBody.UserId)
	if err != nil {
		h.logger.Error(title.Server, zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.firebase.SendToMultiTokens(tokens, reqBody.UserId, reqBody.Message)
	if err != nil {
		h.logger.Error(title.Server, zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
