// AnhCao 2024
package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/AnhCaooo/electric-notifications/internal/config"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type Firebase struct {
	logger       *zap.Logger
	cloudMessage *messaging.Client
	ctx          context.Context
}

// Initialize new a new Firebase instance
func NewFirebase(logger *zap.Logger, ctx context.Context) *Firebase {
	return &Firebase{
		logger:       logger,
		cloudMessage: nil,
		ctx:          ctx,
	}
}

// Establish a connection to the Firebase Cloud Message service
func (fb *Firebase) EstablishConnection() error {
	fullPath, err := config.DecryptFirebaseKeyFile()
	if err != nil {
		return err
	}
	opt := option.WithCredentialsFile(fullPath)
	// Initialize Firebase SDK with Google Application Default credentials
	app, err := firebase.NewApp(fb.ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing connection with Firebase app: %s", err.Error())
	}
	// Get the FCM object
	fb.cloudMessage, err = app.Messaging(fb.ctx)
	if err != nil {
		return fmt.Errorf("error getting Messaging client: %s", err.Error())
	}

	fb.logger.Info("Successfully connected to Firebase Cloud Message platform")
	return nil
}

// Send notification based on a device token
func (fb Firebase) SendToSingleToken(
	token, userId, message string,
) error {
	// define message payload
	payload := &messaging.Message{
		Token: token,
		Data: map[string]string{
			message: message,
		},
	}
	// send a message to the device based on given token
	_, err := fb.cloudMessage.Send(fb.ctx, payload)
	if err != nil {
		return fmt.Errorf("error sending notification to single device: %s", err.Error())
	}
	return nil
}

// Send notification based on multi device tokens
func (fb Firebase) SendToMultiTokens(
	tokens []string,
	userId, message string,
) error {
	payload := &messaging.MulticastMessage{
		Data: map[string]string{
			message: message,
		},
		Tokens: tokens,
	}
	//Send to Multiple Tokens
	batchResponse, err := fb.cloudMessage.SendEachForMulticast(fb.ctx, payload)
	if err != nil {
		return fmt.Errorf("error sending notifications to multi devices: %s", err.Error())
	}

	// check which tokens resulted in errors
	if batchResponse.FailureCount > 0 {
		var failedTokens []string
		for idx, resp := range batchResponse.Responses {
			if !resp.Success {
				// The order of responses corresponds to the order of the registration tokens.
				failedTokens = append(failedTokens, tokens[idx])
			}
		}
		fb.logger.Error("List of tokens that cause failures", zap.Any("tokens", failedTokens))
	}
	return nil
}
