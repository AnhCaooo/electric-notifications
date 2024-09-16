// Created by Anh Cao on 27.08.2024.

package db

import (
	"context"
	"fmt"
	"time"

	"github.com/AnhCaooo/electric-push-notifications/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	databaseName   = "electricApp"
	collectionName = "notificationTokens"
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
	if err = createIndex(Collection, ctx); err != nil {
		return nil, err
	}
	logger.Logger.Info("Successfully connected to database")
	return client, nil
}

func createIndex(collection *mongo.Collection, ctx context.Context) error {
	const weeklyHours = 24 * 7

	//create the index model with the field "timestamp"
	indexModel := mongo.IndexModel{
		Keys: bson.M{"timestamp": 1},
		Options: options.Index().SetExpireAfterSeconds(
			int32((time.Hour * 3 * weeklyHours).Seconds()),
		),
	}
	//Create the index on the token collection
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("mongo index error: %s", err.Error())
	}
	return nil
}
