// AnhCao 2024
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Notification token schema
type NotificationToken struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	UserId    string             `bson:"userId" json:"userId"`
	DeviceId  string             `bson:"deviceId" json:"deviceId"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}

type NotificationMessage struct {
	UserId  string `json:"userId"`
	Message string `json:"message"`
}
