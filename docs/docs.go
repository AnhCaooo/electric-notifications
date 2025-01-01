// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Anh Cao",
            "email": "anhcao4922@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/notifications": {
            "post": {
                "description": "It retrieves the user ID from the request context and decodes the request body to get the notification message.\nThen validates the user ID and retrieves the associated device tokens from the database. Finally, it sends the notification message to the retrieved device tokens using Firebase.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notifications"
                ],
                "summary": "Sends notifications to user devices",
                "parameters": [
                    {
                        "description": "represents a message to be sent to a user.",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.NotificationMessage"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "If the token is successfully inserted into the database.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthenticated/Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "If the user ID in the request body does not match the user ID in the context.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "If there is an error retrieving the device tokens or sending the notifications, it responds with an internal server error.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v1/token": {
            "post": {
                "description": "It extracts the user ID from the request context and decodes the request body to get the notification token details. If the user ID in the request body does not match the user ID in the context, it returns a forbidden error.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notifications"
                ],
                "summary": "Create a notification token that contains the user ID and device token",
                "parameters": [
                    {
                        "description": "represents a token used for sending notifications to  one or more specific device.",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.NotificationToken"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "If the token is successfully inserted into the database.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthenticated/Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "If the user ID in the request body does not match the user ID in the context.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "If there is an error inserting the token into the database.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.NotificationMessage": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Hello, World!"
                },
                "userId": {
                    "type": "string",
                    "example": "1234567890"
                }
            }
        },
        "models.NotificationToken": {
            "type": "object",
            "properties": {
                "deviceId": {
                    "description": "Identifier of the device associated with the notification token.\ntodo: maybe this could be a slice instead of single deviceID. This way we can send notifications to multiple devices that user has.",
                    "type": "string",
                    "example": "1234567890"
                },
                "id": {
                    "description": "Unique identifier for the notification token.",
                    "type": "string",
                    "example": "1234567890"
                },
                "timestamp": {
                    "description": "The time when the notification token was created.",
                    "type": "string",
                    "example": "2025-01-02 14:00:00 +0200 EET"
                },
                "userId": {
                    "description": "Identifier of the user associated with the notification token.",
                    "type": "string",
                    "example": "1234567890"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost:5003",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Notifications API",
	Description:      "Push notifications service for electric application",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
