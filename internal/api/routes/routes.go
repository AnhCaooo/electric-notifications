// AnhCao 2024
package routes

import (
	"net/http"

	"github.com/AnhCaooo/electric-notifications/internal/api/handlers"
)

// Endpoint is the presentation of object which contains values for routing
type Endpoint struct {
	Path    string
	Handler http.HandlerFunc
	Method  string
}

func InitializeEndpoints(handler *handlers.Handler) []Endpoint {
	return []Endpoint{
		{
			Path:    "/v1/ping",
			Handler: handler.Ping,
			Method:  "GET",
		},
		{
			Path:    "/v1/token",
			Handler: handler.CreateToken,
			Method:  "POST",
		}, {
			Path:    "/v1/notifications",
			Handler: handler.SendNotifications,
			Method:  "POST",
		},
	}
}
