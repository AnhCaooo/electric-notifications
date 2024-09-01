// Created by Anh Cao on 27.08.2024.

package routes

import (
	"net/http"
)

// Endpoint is the presentation of object which contains values for routing
type Endpoint struct {
	Path    string
	Handler http.HandlerFunc
	Method  string
}

var Endpoints = []Endpoint{
	{
		Path:    "/v1/get-request",
		Handler: nil,
		Method:  "GET",
	},
	{
		Path:    "/v1/post-request",
		Handler: nil,
		Method:  "POST",
	}, {
		Path:    "/v1/patch-request",
		Handler: nil,
		Method:  "PATCH",
	},
	{
		Path:    "/v1/put-request",
		Handler: nil,
		Method:  "PUT",
	},
	{
		Path:    "/v1/delete-request",
		Handler: nil,
		Method:  "DELETE",
	},
}
