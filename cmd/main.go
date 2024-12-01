// AnhCao 2024

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AnhCaooo/electric-notifications/internal/api/handlers"
	"github.com/AnhCaooo/electric-notifications/internal/api/middleware"
	"github.com/AnhCaooo/electric-notifications/internal/api/routes"
	"github.com/AnhCaooo/electric-notifications/internal/cache"
	"github.com/AnhCaooo/electric-notifications/internal/config"
	"github.com/AnhCaooo/electric-notifications/internal/constants"
	"github.com/AnhCaooo/electric-notifications/internal/db"
	"github.com/AnhCaooo/electric-notifications/internal/firebase"
	"github.com/AnhCaooo/electric-notifications/internal/logger"
	"github.com/AnhCaooo/electric-notifications/internal/models"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	// Initialize logger
	logger := logger.Init()
	defer logger.Sync()

	configuration := &models.Config{}
	// Read configuration file
	err := config.ReadFile(configuration)
	if err != nil {
		logger.Error(constants.Server, zap.Error(err))
		os.Exit(1)
	}

	// Initialize in-memory cache
	cache := cache.NewCache(logger)

	// Initialize database connection
	mongo := db.NewMongo(ctx, &configuration.Database, logger)
	mongoClient, err := mongo.EstablishConnection()
	if err != nil {
		logger.Error(constants.Server, zap.Error(err))
		os.Exit(1)
	}
	defer mongoClient.Disconnect(ctx)

	// Initialize FCM connection
	firebase := firebase.NewFirebase(logger, ctx)
	if err = firebase.EstablishConnection(); err != nil {
		logger.Error(constants.Server, zap.Error(err))
		os.Exit(1)
	}

	// Initialize Middleware
	middleware := middleware.NewMiddleware(logger, configuration)
	// Initialize Handler
	handler := handlers.NewHandler(logger, cache, mongo, firebase)
	// Initialize Endpoints pool
	endpoints := routes.InitializeEndpoints(handler)
	// Initial new router
	r := mux.NewRouter()

	// Apply middlewares
	middlewares := []func(http.Handler) http.Handler{
		middleware.Logger,
		middleware.Authenticate,
	}
	for _, mw := range middlewares {
		r.Use(mw)
	}

	for _, endpoint := range endpoints {
		r.HandleFunc(endpoint.Path, endpoint.Handler).Methods(endpoint.Method)
	}

	r.MethodNotAllowedHandler = http.HandlerFunc(handler.NotAllowed)
	r.NotFoundHandler = http.HandlerFunc(handler.NotFound)

	// Start server
	logger.Info("Server started on", zap.String("port", configuration.Server.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", configuration.Server.Port), r))
}
