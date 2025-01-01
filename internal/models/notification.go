// AnhCao 2024
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// NotificationToken represents a token used for sending notifications to a specific device.
type NotificationToken struct {
	// Unique identifier for the notification token.
	ID bson.ObjectID `bson:"_id" json:"id" example:"1234567890"`
	// Identifier of the user associated with the notification token.
	UserId string `bson:"userId" json:"userId" example:"1234567890"`
	// Identifier of the device associated with the notification token.
	// todo: maybe this could be a slice instead of single deviceID. This way we can send notifications to multiple devices that user has.
	DeviceId string `bson:"deviceId" json:"deviceId" example:"1234567890"`
	// The time when the notification token was created.
	Timestamp time.Time `bson:"timestamp" json:"timestamp" example:"2025-01-02 14:00:00 +0200 EET"`
}

// NotificationMessage represents a message to be sent to a user.
type NotificationMessage struct {
	UserId  string `json:"userId" example:"1234567890"`
	Message string `json:"message" example:"Hello, World!"`
}
