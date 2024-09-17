// Created by Anh Cao on 27.08.2024.

package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/AnhCaooo/electric-push-notifications/internal/api/handlers"
	"github.com/AnhCaooo/electric-push-notifications/internal/api/middleware"
	"github.com/AnhCaooo/electric-push-notifications/internal/api/routes"
	"github.com/AnhCaooo/electric-push-notifications/internal/cache"
	title "github.com/AnhCaooo/electric-push-notifications/internal/constants"
	"github.com/AnhCaooo/electric-push-notifications/internal/db"
	"github.com/AnhCaooo/electric-push-notifications/internal/firebase"
	"github.com/AnhCaooo/electric-push-notifications/internal/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const uri string = "mongodb://<dummy_user>:<dummy_pass>@localhost:27017/?timeoutMS=5000"

func main() {
	ctx := context.Background()
	// Initialize logger
	logger.Init()
	// Initialize cache
	cache.NewCache()
	// Initialize database connection
	mongo, err := db.Init(ctx, uri)
	if err != nil {
		logger.Logger.Error(title.Server, zap.Error(err))
		os.Exit(1)
	}
	defer mongo.Disconnect(ctx)
	// Initialize FCM connection
	if err = firebase.Init(ctx); err != nil {
		logger.Logger.Error(title.Server, zap.Error(err))
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
	logger.Logger.Info("Server started on :8002")
	log.Fatal(http.ListenAndServe(":8002", r))
}
