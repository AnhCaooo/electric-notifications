// AnhCao 2024
// This folder tends to contains basic handlers for such cases: NotFound, NotAllowed, etc.
package handlers

import (
	"fmt"
	"net/http"

	"github.com/AnhCaooo/electric-notifications/internal/cache"
	"github.com/AnhCaooo/electric-notifications/internal/db"
	"github.com/AnhCaooo/electric-notifications/internal/firebase"
	"go.uber.org/zap"
)

type Handler struct {
	logger   *zap.Logger
	cache    *cache.Cache
	mongo    *db.Mongo
	firebase *firebase.Firebase
	workerID int
}

// NewHandler returns a new Handler instance
func NewHandler(
	logger *zap.Logger,
	cache *cache.Cache,
	mongo *db.Mongo,
	firebase *firebase.Firebase,
	workerID int,
) *Handler {
	if mongo == nil {
		logger.Warn(fmt.Sprintf("[worker_%d] mongoDB client is nil, using mock or no-op database", workerID))
	}
	return &Handler{
		logger:   logger,
		cache:    cache,
		mongo:    mongo,
		firebase: firebase,
		workerID: workerID,
	}
}

// return response when request url is not found
func (h Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	h.logger.Info(fmt.Sprintf("[worker_%d] undefined endpoint", h.workerID), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
	w.Write([]byte("404 - Not found"))
}

// return response when request method is not allowed
func (h Handler) NotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	h.logger.Info(fmt.Sprintf("[worker_%d] method not allowed", h.workerID), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
	w.Write([]byte("405 - Method not allowed"))
}

// Ping the connection to the server
func (h Handler) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}
