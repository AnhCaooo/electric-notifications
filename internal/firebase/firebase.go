// AnhCao 2024
package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/AnhCaooo/electric-notifications/internal/config"
	"github.com/AnhCaooo/electric-notifications/internal/logger"
	"google.golang.org/api/option"
)

var FcmClient *messaging.Client

func Init(ctx context.Context) error {
	fullPath, err := config.DecryptFirebaseKeyFile()
	if err != nil {
		return err
	}
	opt := option.WithCredentialsFile(fullPath)
	// Initialize Firebase SDK with Google Application Default credentials
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing connection with Firebase app: %s", err.Error())
	}
	// Get the FCM object
	FcmClient, err = app.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("error getting Messaging client: %s", err.Error())
	}
	logger.Logger.Info("Successfully connected to Firebase Cloud Message platform")
	return nil
}
