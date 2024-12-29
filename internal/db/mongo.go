// AnhCao 2024
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/AnhCaooo/electric-notifications/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

// Mongo represents a MongoDB client with configuration, logger, context, and collection information.
type Mongo struct {
	config     *models.Database
	logger     *zap.Logger
	ctx        context.Context
	Client     *mongo.Client
	collection *mongo.Collection
}

// NewMongo initializes a new Mongo instance with the provided context, database configuration, and logger.
// It returns a pointer to the newly created Mongo instance.
func NewMongo(ctx context.Context, config *models.Database, logger *zap.Logger) *Mongo {
	return &Mongo{
		config:     config,
		logger:     logger,
		ctx:        ctx,
		collection: nil,
	}
}

// Function to connect to mongo database instance and create collection if it does not exist
func (db *Mongo) EstablishConnection() (err error) {
	clientOptions := options.Client().ApplyURI(db.getURI())
	db.Client, err = mongo.Connect(clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to database. Error: %s", err.Error())
	}

	err = db.Client.Ping(db.ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping database. Error: %s", err.Error())
	}

	db.collection = db.Client.Database(db.config.Name).Collection(db.config.Collection)
	if err = db.createIndex(db.collection); err != nil {
		return err
	}
	db.logger.Info("Successfully connected to database")
	return nil
}

func (db Mongo) getURI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/?timeoutMS=5000", db.config.Username, db.config.Password, "localhost", db.config.Port)
}

// createIndex creates an index on the specified MongoDB collection.
// The index is created on the "timestamp" field and is set to expire
// documents after a certain period of time. The expiration time is
// calculated as 3 weeks.
func (db Mongo) createIndex(collection *mongo.Collection) error {
	const weeklyHours = 24 * 7

	//create the index model with the field "timestamp"
	indexModel := mongo.IndexModel{
		Keys: bson.M{"timestamp": 1},
		Options: options.Index().SetExpireAfterSeconds(
			int32((time.Hour * 3 * weeklyHours).Seconds()),
		),
	}
	//Create the index on the token collection
	_, err := collection.Indexes().CreateOne(db.ctx, indexModel)
	if err != nil {
		return fmt.Errorf("mongo index error: %s", err.Error())
	}
	return nil
}

// Insert a notification token
func (db Mongo) InsertToken(token models.NotificationToken) error {
	// Check if the token already exists
	filter := bson.D{{Key: "deviceId", Value: token.DeviceId}}
	res := db.collection.FindOne(db.ctx, filter)

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			// If token does not exist then insert it
			token.ID = bson.NewObjectID()
			_, err := db.collection.InsertOne(db.ctx, token)
			return fmt.Errorf("failed to insert token: %s", err.Error())
		}
		return res.Err()
	}

	// If token exists update the timestamp to now
	_, err := db.collection.UpdateOne(db.ctx, filter, bson.M{"$set": bson.M{"timestamp": time.Now().UTC()}})
	if err != nil {
		return fmt.Errorf("failed to update existing token: %s", err.Error())
	}
	return nil
}

// Get all the tokens registered for a user
func (db Mongo) GetTokens(userId string) ([]string, error) {
	filter := bson.D{{Key: "userId", Value: userId}}
	tokenCursor, err := db.collection.Find(db.ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find tokens for user: %s", err.Error())
	}

	tokens := make([]string, 0)
	for tokenCursor.Next(db.ctx) {
		var token models.NotificationToken
		err = tokenCursor.Decode(&token)
		tokens = append(tokens, token.DeviceId)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to decode notification token: %s", err.Error())
	}
	return tokens, nil
}
