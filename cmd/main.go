// AnhCao 2024

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AnhCaooo/electric-notifications/internal/api/handlers"
	"github.com/AnhCaooo/electric-notifications/internal/api/middleware"
	"github.com/AnhCaooo/electric-notifications/internal/api/routes"
	"github.com/AnhCaooo/electric-notifications/internal/cache"
	"github.com/AnhCaooo/electric-notifications/internal/config"
	"github.com/AnhCaooo/electric-notifications/internal/constants"
	"github.com/AnhCaooo/electric-notifications/internal/db"
	"github.com/AnhCaooo/electric-notifications/internal/firebase"
	"github.com/AnhCaooo/electric-notifications/internal/models"
	"github.com/AnhCaooo/go-goods/log"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx := context.Background()
	// Initialize logger
	logger := log.InitLogger(zapcore.InfoLevel)
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

	handler := handlers.NewHandler(logger, cache, mongo, firebase)
	middleware := middleware.NewMiddleware(logger, configuration)
	endpoints := routes.InitializeEndpoints(handler)

	// Initial new Mux router
	r := initializeRouter(handler, middleware, endpoints)
	// Start server
	run(ctx, logger, configuration, r)
}

func initializeRouter(handler *handlers.Handler, middleware *middleware.Middleware, endpoints []routes.Endpoint) *mux.Router {
	r := mux.NewRouter()
	// Apply middlewares
	middlewares := []func(http.Handler) http.Handler{
		middleware.Logger,
		middleware.Authenticate,
	}
	for _, mw := range middlewares {
		r.Use(mw)
	}

	// Apply endpoint handlers
	for _, endpoint := range endpoints {
		r.HandleFunc(endpoint.Path, endpoint.Handler).Methods(endpoint.Method)
	}

	r.MethodNotAllowedHandler = http.HandlerFunc(handler.NotAllowed)
	r.NotFoundHandler = http.HandlerFunc(handler.NotFound)
	return r
}

func run(ctx context.Context, logger *zap.Logger, config *models.Config, r *mux.Router) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Server.Port),
		Handler: r,
	}

	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Run the server in a separate goroutine
	go func() {
		logger.Info("Server starting", zap.String("port", config.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server error", zap.Error(err))
		}
	}()

	// Wait for termination signal
	select {
	case <-ctx.Done(): // Context cancellation
		logger.Warn("Context canceled")
	case <-stop: // OS signal received
		logger.Info("Termination signal received")
	}

	// Create a new context with timeout for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("Shutting down server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}
