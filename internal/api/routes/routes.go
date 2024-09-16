// Created by Anh Cao on 27.08.2024.

package routes

import (
	"net/http"

	"github.com/AnhCaooo/electric-push-notifications/internal/api/handlers"
)

// Endpoint is the presentation of object which contains values for routing
type Endpoint struct {
	Path    string
	Handler http.HandlerFunc
	Method  string
}

var Endpoints = []Endpoint{
	{
		Path:    "/v1/tokens",
		Handler: handlers.CreateToken,
		Method:  "POST",
	}, {
		Path:    "/v1/notifications",
		Handler: nil,
		Method:  "POST",
	},
}
