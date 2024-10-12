// Created by Anh Cao on 27.08.2024.
// This file contains all business request functions
package api

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	title "github.com/AnhCaooo/electric-push-notifications/internal/constants"
	"github.com/AnhCaooo/electric-push-notifications/internal/db"
	"github.com/AnhCaooo/electric-push-notifications/internal/firebase"
	"github.com/AnhCaooo/electric-push-notifications/internal/helpers"
	"github.com/AnhCaooo/electric-push-notifications/internal/logger"
	"github.com/AnhCaooo/electric-push-notifications/internal/models"
	"github.com/AnhCaooo/electric-push-notifications/internal/notification"
)

// Ping the connection to the server
func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

// create token to database or update time live if it is existing
func CreateToken(w http.ResponseWriter, r *http.Request) {
	reqBody, err := helpers.DecodeRequest[models.NotificationToken](r)
	if err != nil {
		logger.Logger.Error(title.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the token into the database
	err = db.InsertToken(db.Collection, reqBody, context.TODO())
	if err != nil {
		// You should have better error handling here
		logger.Logger.Error(title.Server, zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SendNotifications(w http.ResponseWriter, r *http.Request) {
	reqBody, err := helpers.DecodeRequest[models.NotificationMessage](r)
	if err != nil {
		logger.Logger.Error(title.Client, zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// retrieve all associated device tokens with given userId
	tokens, err := db.GetTokens(db.Collection, context.TODO(), reqBody.UserId)
	if err != nil {
		logger.Logger.Error(title.Server, zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = notification.SendToMultiTokens(firebase.FcmClient, context.TODO(), tokens, reqBody.UserId, reqBody.Message)
	if err != nil {
		logger.Logger.Error(title.Server, zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
