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
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	// Initialize logger
	logger.Init()
	defer logger.Logger.Sync()

	// Read configuration file
	err := config.ReadFile(&config.Config)
	if err != nil {
		logger.Logger.Error(constants.Server, zap.Error(err))
		os.Exit(1)
	}

	// Initialize in-memory cache
	cache.NewCache()

	// Initialize database connection
	mongo, err := db.Init(ctx, config.Config.Database)
	if err != nil {
		logger.Logger.Error(constants.Server, zap.Error(err))
		os.Exit(1)
	}
	defer mongo.Disconnect(ctx)

	// Initialize FCM connection
	if err = firebase.Init(ctx); err != nil {
		logger.Logger.Error(constants.Server, zap.Error(err))
		os.Exit(1)
	}

	// Initial new router
	r := mux.NewRouter()
	for _, endpoint := range routes.Endpoints {
		r.HandleFunc(endpoint.Path, endpoint.Handler).Methods(endpoint.Method)
	}
	r.MethodNotAllowedHandler = http.HandlerFunc(handlers.NotAllowed)
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFound)

	// Middleware
	r.Use(middleware.Logger)

	// Start server
	logger.Logger.Info("Server started on", zap.String("port", config.Config.Server.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Config.Server.Port), r))
}
