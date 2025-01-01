// AnhCao 2024

package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/AnhCaooo/electric-notifications/docs"
	"github.com/AnhCaooo/electric-notifications/internal/api"
	"github.com/AnhCaooo/electric-notifications/internal/config"
	"github.com/AnhCaooo/electric-notifications/internal/constants"
	"github.com/AnhCaooo/electric-notifications/internal/db"
	"github.com/AnhCaooo/electric-notifications/internal/firebase"
	"github.com/AnhCaooo/electric-notifications/internal/models"
	"github.com/AnhCaooo/electric-notifications/internal/rabbitmq"
	"github.com/AnhCaooo/go-goods/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//	@title			Notifications API
//	@version		1.0.0
//	@description	Push notifications service for electric application

//	@contact.name	Anh Cao
//	@contact.email	anhcao4922@gmail.com

// @host		localhost:5003
// @BasePath	/
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

	var wg sync.WaitGroup
	// Error channel to listen for errors from goroutines
	errChan := make(chan error, 2)
	// StopChan to listen for stop signal
	stopChan := make(chan struct{})

	// HTTP server
	httpServer := api.NewHTTPServer(ctx, logger, config, mongo, firebase)
	httpServer.Start(1, errChan, &wg)
	// RabbitMQ consumer
	rabbitMQ := rabbitmq.NewRabbit(ctx, &config.MessageBroker, logger, mongo, firebase)
	if err := rabbitMQ.EstablishConnection(); err != nil {
		logger.Fatal("failed to establish connection with RabbitMQ", zap.Error(err))
	}
	rabbitMQ.StartConsumers(&wg, errChan, stopChan)

	// Monitor all errors from errChan and log them
	go func() {
		for err := range errChan {
			logger.Error("error occurred", zap.Error(err))
		}
	}()

	// Wait for termination signal
	<-stop
	logger.Info("Termination signal received")
	close(stopChan)
	httpServer.Stop()
	// Wait for all goroutines to finish
	wg.Wait()
	// Signal all errors to stop
	close(errChan)
	logger.Info("HTTP server and RabbitMQ exited gracefully")
}
