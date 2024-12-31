// AnhCao 2024
package helpers

import (
	"fmt"

	"github.com/AnhCaooo/electric-notifications/internal/models"
)

// todo: allow user to customize the message that will be pushed to user
// todo: use AI to generate useful message for user
// GenerateNotificationMessageForSpotPrice generates a notification message for the spot price.
// It takes a PricesMessage struct as input and returns a formatted string containing the price for tomorrow.
func GenerateNotificationMessageForSpotPrice(data *models.PricesMessage) string {
	return fmt.Sprintf("Tomorrow price is %f", data.Data.Tomorrow.Prices.Data[0].Price)
}
