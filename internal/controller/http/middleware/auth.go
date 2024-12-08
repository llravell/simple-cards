package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
)

const (
	TokenCookieName = "user-token"
)

type contextKey string

var userUUIDContextKey contextKey = "userUUID"

type JWTParser interface {
	Parse(tokenString string) (*jwt.Token, error)
}

type authenticator struct {
	jwtParser JWTParser
	log       zerolog.Logger
}

func (auth *authenticator) parseUserUUIDFromRequest(r *http.Request) string {
	tokenCookie, err := r.Cookie(TokenCookieName)
	if err != nil {
		auth.log.Error().Err(err).Msg("jwt cookie finding failed")

		return ""
	}

	token, err := auth.jwtParser.Parse(tokenCookie.Value)
	if err != nil {
		auth.log.Error().Err(err).Msg("jwt parsing failed")

		return ""
	}

	if !token.Valid {
		auth.log.Error().Msg("got invalid jwt")

		return ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		auth.log.Error().Err(err).Msg("jwt claims type casting failed")

		return ""
	}

	userUUID, ok := claims["sub"].(string)
	if !ok {
		auth.log.Error().Err(err).Msg("userUUID type casting failed")

		return ""
	}

	return userUUID
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

func NewAuthMiddleware(jwtParser JWTParser, log zerolog.Logger) func(next http.Handler) http.Handler {
	auth := &authenticator{
		jwtParser: jwtParser,
		log:       log,
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
