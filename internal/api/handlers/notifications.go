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

// CreateToken handles the creation of a notification token for a user.
//
//	@Summary		Create a notification token that contains the user ID and device tokens
//	@Description	It extracts the user ID from the request context and decodes the request body to get the notification token details. If the user ID in the request body does not match the user ID in the context, it returns a forbidden error.
//	@Tags			notifications
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.NotificationToken	true	"represents a token used for sending notifications to  one or more specific device.""
//	@Success		200	{object}	string "If the token is successfully inserted into the database.""
//	@Failure		400	{string}	string "Invalid request"
//	@Failure		403	{string}	string "If the user ID in the request body does not match the user ID in the context."
//	@Failure		401	{string}	string "Unauthenticated/Unauthorized"
//	@Failure		500	{string}	string "If there is an error inserting the token into the database."
//	@Router			/v1/token [post]
//
// todo: any response value?
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

// SendNotifications sends notifications to user devices.
//
//	@Summary		Sends notifications to user devices
//	@Description	It retrieves the user ID from the request context and decodes the request body to get the notification message.
//	@Description	Then validates the user ID and retrieves the associated device tokens from the database. Finally, it sends the notification message to the retrieved device tokens using Firebase.
//
//	@Tags			notifications
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.NotificationMessage	true	"represents a message to be sent to all devices that user has."
//	@Success		200	{object}	string "If the token is successfully inserted into the database.""
//	@Failure		400	{string}	string "Invalid request"
//	@Failure		403	{string}	string "If the user ID in the request body does not match the user ID in the context."
//	@Failure		401	{string}	string "Unauthenticated/Unauthorized"
//	@Failure		500	{string}	string "If there is an error retrieving the device tokens or sending the notifications, it responds with an internal server error."
//	@Router			/v1/notifications [post]
//
// todo: any response value?
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
