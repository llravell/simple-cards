package middleware

import (
	"context"
	"net/http"

	"github.com/llravell/simple-cards/internal/entity"
	"github.com/rs/zerolog"
)

const (
	TokenCookieName = "user-token"
)

type contextKey string

var userUUIDContextKey contextKey = "userUUID"

type authenticator struct {
	secret []byte
	log    zerolog.Logger
}

func (auth *authenticator) parseUserUUIDFromRequest(r *http.Request) string {
	tokenCookie, err := r.Cookie(TokenCookieName)
	if err != nil {
		auth.log.Error().Err(err).Msg("jwt cookie finding failed")

		return ""
	}

	token, claims, err := entity.ParseJWTString(tokenCookie.Value, auth.secret)
	if err != nil {
		auth.log.Error().Err(err).Msg("jwt parsing failed")

		return ""
	}

	if !token.Valid {
		auth.log.Error().Msg("got invalid jwt")

		return ""
	}

	return claims.UserUUID
}

func (auth *authenticator) provideUserUUIDToRequestContext(r *http.Request, userUUID string) *http.Request {
	ctx := context.WithValue(r.Context(), userUUIDContextKey, userUUID)

	return r.WithContext(ctx)
}

func (auth *authenticator) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userUUID := auth.parseUserUUIDFromRequest(r)

		if userUUID == "" {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, auth.provideUserUUIDToRequestContext(r, userUUID))
	})
}

func NewAuthMiddleware(secretKey string, log zerolog.Logger) func(next http.Handler) http.Handler {
	auth := &authenticator{
		secret: []byte(secretKey),
		log:    log,
	}

	return auth.Handler
}

func GetUserUUIDFromRequest(r *http.Request) string {
	v := r.Context().Value(userUUIDContextKey)
	userUUID, ok := v.(string)

	if !ok {
		return ""
	}

	return userUUID
}
