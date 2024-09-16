// Created by Anh Cao on 27.08.2024.

package handlers

import (
	"context"
	"net/http"

	"github.com/AnhCaooo/electric-push-notifications/internal/db"
	"github.com/AnhCaooo/electric-push-notifications/internal/helpers"
	"github.com/AnhCaooo/electric-push-notifications/internal/logger"
	"github.com/AnhCaooo/electric-push-notifications/internal/models"
	"go.uber.org/zap"
)

// return response when request url is not found
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	logger.Logger.Info("undefined endpoint", zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
	w.Write([]byte("404 - Not found"))
}

// return response when request method is not allowed
func NotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	logger.Logger.Info("method not allowed", zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
	w.Write([]byte("405 - Method not allowed"))
}

// create token to database or update time live if it is existing
func CreateToken(w http.ResponseWriter, r *http.Request) {
	reqBody, err := helpers.DecodeRequest[models.NotificationToken](r)
	if err != nil {
		logger.Logger.Error("[request] failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the token into the database
	err = db.InsertToken(db.Collection, reqBody, context.TODO())
	if err != nil {
		// You should have better error handling here
		logger.Logger.Error("[server]", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
