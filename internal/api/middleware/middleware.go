// AnhCao 2024
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/AnhCaooo/electric-notifications/internal/constants"
	"github.com/AnhCaooo/electric-notifications/internal/models"
	"github.com/AnhCaooo/go-goods/auth"
	"go.uber.org/zap"
)

type Middleware struct {
	logger *zap.Logger
	config *models.Config
}

func NewMiddleware(logger *zap.Logger, config *models.Config) *Middleware {
	return &Middleware{
		logger: logger,
		config: config,
	}
}

// log the coming request to the server
func (m Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info("request received", zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}

// read the token from request and do verify the access token
func (m Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusForbidden)
			m.logger.Info("permission Denied: No token provided")
			w.Write([]byte("403 - Forbidden"))
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := auth.VerifyToken(tokenString, m.config.Supabase.Auth.JwtSecret)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			m.logger.Error("unauthorized request", zap.Error(err))
			w.Write([]byte("401 - Unauthorized"))
			return
		}

		// due to 'Supabase' authentication, it stores userId via "sub" field
		userID, err := auth.ExtractValueFromTokenClaim(token, "sub")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			m.logger.Error("unauthorized request", zap.Error(err))
			w.Write([]byte("401 - Unauthorized"))
			return
		}

		// Add userID to the context
		ctx := context.WithValue(r.Context(), constants.UserIdKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
