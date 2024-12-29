// AnhCao 2024

package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/AnhCaooo/electric-notifications/internal/api"
	"github.com/AnhCaooo/electric-notifications/internal/config"
	"github.com/AnhCaooo/electric-notifications/internal/constants"
	"github.com/AnhCaooo/electric-notifications/internal/db"
	"github.com/AnhCaooo/electric-notifications/internal/firebase"
	"github.com/AnhCaooo/electric-notifications/internal/models"
	"github.com/AnhCaooo/go-goods/log"
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

	// Initialize database connection
	mongo := db.NewMongo(ctx, &configuration.Database, logger)
	if err := mongo.EstablishConnection(); err != nil {
		logger.Error(constants.Server, zap.Error(err))
		os.Exit(1)
	}
	defer mongo.Client.Disconnect(ctx)

	// Initialize FCM connection
	firebase := firebase.NewFirebase(logger, ctx)
	if err = firebase.EstablishConnection(); err != nil {
		logger.Error(constants.Server, zap.Error(err))
		os.Exit(1)
	}
	// Start server
	run(ctx, logger, configuration, mongo, firebase)
}

// run initializes and starts the HTTP server, sets up signal handling for graceful shutdown,
// and manages the server lifecycle.
func run(ctx context.Context, logger *zap.Logger, config *models.Config, mongo *db.Mongo, firebase *firebase.Firebase) {
	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Error channel to listen for errors from goroutines
	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	// HTTP server
	httpServer := api.NewHTTPServer(ctx, logger, config, mongo, firebase)
	httpServer.Start(1, errChan, &wg)

	<-stop
	logger.Info("Termination signal received")
	httpServer.Stop()
	// Wait for all goroutines to finish
	wg.Wait()
	// Signal all errors to stop
	close(errChan)
	logger.Info("HTTP server exited gracefully")
}
