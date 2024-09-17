package notification

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/messaging"
	"github.com/AnhCaooo/electric-push-notifications/internal/logger"
	"go.uber.org/zap"
)

// Send notification based on a device token
func SendToSingleToken(
	fcmClient *messaging.Client,
	ctx context.Context,
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
	_, err := fcmClient.Send(ctx, payload)
	if err != nil {
		return fmt.Errorf("error sending notification to single device: %s", err.Error())
	}
	logger.Logger.Info("Successfully sent notification to single device")
	return nil
}

// Send notification based on multi device tokens
func SendToMultiTokens(
	fcmClient *messaging.Client,
	ctx context.Context,
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
	batchResponse, err := fcmClient.SendEachForMulticast(ctx, payload)
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
		logger.Logger.Info("List of tokens that cause failures", zap.Any("tokens", failedTokens))
	} else {
		logger.Logger.Info("Successfully sent notifications to all  devices")
	}
	return nil
}
