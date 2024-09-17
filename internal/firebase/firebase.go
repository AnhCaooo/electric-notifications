package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var FcmClient *messaging.Client

func Init(ctx context.Context) error {
	opt := option.WithCredentialsFile("../config/firebaseKey.json")
	// Initialize Firebase SDK with Google Application Default credentials
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing connection with Firebase app: %s", err.Error())
	}
	// Get the FCM object
	FcmClient, err = app.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("failed to get FCM instance: %s", err.Error())
	}

	return nil
}
