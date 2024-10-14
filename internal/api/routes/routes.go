// AnhCao 2024
package routes

import (
	"net/http"

	"github.com/AnhCaooo/electric-push-notifications/internal/api"
)

// Endpoint is the presentation of object which contains values for routing
type Endpoint struct {
	Path    string
	Handler http.HandlerFunc
	Method  string
}

var Endpoints = []Endpoint{
	{
		Path:    "/v1/ping",
		Handler: api.Ping,
		Method:  "GET",
	},
	{
		Path:    "/v1/tokens",
		Handler: api.CreateToken,
		Method:  "POST",
	}, {
		Path:    "/v1/notifications",
		Handler: api.SendNotifications,
		Method:  "POST",
	},
}
