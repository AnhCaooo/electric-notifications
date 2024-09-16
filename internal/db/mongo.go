// Created by Anh Cao on 27.08.2024.

package db

import (
	"context"
	"fmt"
	"time"

	"github.com/AnhCaooo/electric-push-notifications/internal/logger"
	"github.com/AnhCaooo/electric-push-notifications/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Insert a notification token
func InsertToken(
	collection *mongo.Collection,
	token models.NotificationToken,
	ctx context.Context,
) error {
	// Check if the token already exists
	filter := bson.D{{Key: "deviceId", Value: token.DeviceId}}
	res := collection.FindOne(ctx, filter)

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			// If token does not exist then insert it
			token.ID = primitive.NewObjectID()
			_, err := collection.InsertOne(ctx, token)
			return fmt.Errorf("failed to insert token: %s", err.Error())
		}
		return res.Err()
	}

	// If token exists update the timestamp to now
	_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"timestamp": time.Now().UTC()}})
	return fmt.Errorf("failed to update existing token: %s", err.Error())
}

// Get all the tokens registered for a user
func GetTokens(
	collection *mongo.Collection,
	ctx context.Context,
	userId string,
) ([]string, error) {
	filter := bson.D{{Key: "userId", Value: userId}}
	tokenCursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find tokens for user: %s", err.Error())
	}

	tokens := make([]string, 0)
	for tokenCursor.Next(ctx) {
		var token models.NotificationToken
		err = tokenCursor.Decode(&token)
		tokens = append(tokens, token.DeviceId)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to decode notification token: %s", err.Error())
	}
	return tokens, nil
}
