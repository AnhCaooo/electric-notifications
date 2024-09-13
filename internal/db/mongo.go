// Created by Anh Cao on 27.08.2024.

package db

import (
	"context"
	"fmt"

	"github.com/AnhCaooo/electric-push-notifications/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Function to connect to mongo database instance
func Init(ctx context.Context, URI string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("[server] failed to connect to database. Error: %s", err.Error())
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Logger.Fatal(err.Error())
		return nil, fmt.Errorf("[server] failed to ping database. Error: %s", err.Error())
	}

	logger.Logger.Info("Successfully connected to database")
	return client, nil
}
