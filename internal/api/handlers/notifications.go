// AnhCao 2024
package handlers

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/AnhCaooo/electric-notifications/internal/constants"
	"github.com/AnhCaooo/electric-notifications/internal/models"
	"github.com/AnhCaooo/go-goods/encode"
)

// create token to Mongo or update time live if it exists
func (h Handler) CreateToken(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(constants.UserIdKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	reqBody, err := encode.DecodeRequest[models.NotificationToken](r)
	if err != nil {
		h.logger.Error(constants.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.UserId != "" && reqBody.UserId != userId {
		err = fmt.Errorf("given `user_id` %s is different from `user_id` in `access_token`", reqBody.UserId)
		h.logger.Error(constants.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	reqBody.UserId = userId

	// Insert the token into the database
	err = h.mongo.InsertToken(reqBody)
	if err != nil {
		// You should have better error handling here
		h.logger.Error(constants.Server, zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) SendNotifications(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(constants.UserIdKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	reqBody, err := encode.DecodeRequest[models.NotificationMessage](r)
	if err != nil {
		h.logger.Error(constants.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.UserId != "" && reqBody.UserId != userId {
		err = fmt.Errorf("given `user_id` %s is different from `user_id` in `access_token`", reqBody.UserId)
		h.logger.Error(constants.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	reqBody.UserId = userId

	// retrieve all associated device tokens with given userId
	tokens, err := h.mongo.GetTokens(reqBody.UserId)
	if err != nil {
		h.logger.Error(constants.Server, zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.firebase.SendToMultiTokens(tokens, reqBody.UserId, reqBody.Message)
	if err != nil {
		h.logger.Error(constants.Server, zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
