// Created by Anh Cao on 27.08.2024.

package db

import (
	"context"
	"fmt"

	"github.com/AnhCaooo/electric-push-notifications/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	databaseName   = "electricApp"
	collectionName = "pushNotifications"
)
var Collection *mongo.Collection

// Function to connect to mongo database instance and create collection if it does not exist
func Init(ctx context.Context, URI string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database. Error: %s", err.Error())
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database. Error: %s", err.Error())
	}

	Collection = client.Database(databaseName).Collection(collectionName)
	logger.Logger.Info("Successfully connected to database")
	return client, nil
}
