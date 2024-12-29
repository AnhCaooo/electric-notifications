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
		h.logger.Error(fmt.Sprintf("[worker_%d] %s failed to decode request", h.workerID, constants.Client), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.UserId != "" && reqBody.UserId != userId {
		errMsg := fmt.Sprintf("[worker_%d] %s given `user_id` %s is different from `user_id` in `access_token`", h.workerID, constants.Client, reqBody.UserId)
		h.logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusForbidden)
		return
	}
	reqBody.UserId = userId

	// Insert the token into the database
	err = h.mongo.InsertToken(reqBody)
	if err != nil {
		h.logger.Error(fmt.Sprintf("[worker_%d] %s failed to insert token", h.workerID, constants.Server), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info(fmt.Sprintf("[worker_%d] insert token successfully", h.workerID))
}

func (h Handler) SendNotifications(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(constants.UserIdKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	reqBody, err := encode.DecodeRequest[models.NotificationMessage](r)
	if err != nil {
		h.logger.Error(fmt.Sprintf("[worker_%d] %s failed to decode request", h.workerID, constants.Client), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.UserId != "" && reqBody.UserId != userId {
		errMsg := fmt.Sprintf("[worker_%d] %s given `user_id` %s is different from `user_id` in `access_token`", h.workerID, constants.Client, reqBody.UserId)
		h.logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusForbidden)
		return
	}
	reqBody.UserId = userId

	// retrieve all associated device tokens with given userId
	tokens, err := h.mongo.GetTokens(reqBody.UserId)
	if err != nil {
		h.logger.Error(fmt.Sprintf("[worker_%d] %s failed to get tokens", h.workerID, constants.Server), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.firebase.SendToMultiTokens(tokens, reqBody.UserId, reqBody.Message)
	if err != nil {
		h.logger.Error(fmt.Sprintf("[worker_%d] %s failed to send multi tokens", h.workerID, constants.Server), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info(fmt.Sprintf("[worker_%d] send tokens successfully", h.workerID))
}
