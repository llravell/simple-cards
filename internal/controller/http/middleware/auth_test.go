package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	testutils "github.com/llravell/simple-cards/internal"
	"github.com/llravell/simple-cards/internal/controller/http/middleware"
	"github.com/llravell/simple-cards/pkg/auth"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	router := chi.NewRouter()
	authMiddleware := middleware.NewAuthMiddleware(
		auth.NewJWTManager(testutils.JWTSecretKey),
		zerolog.Nop(),
	)

	router.Use(authMiddleware)
	router.Post("/", echoHandler(t))

	ts := httptest.NewServer(router)

	t.Run("Middleware return unauthorized status code if token does not exist", func(t *testing.T) {
		res, _ := testutils.SendTestRequest(t, ts, http.MethodPost, "/", http.NoBody, map[string]string{})
		defer res.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("Middleware call original handler if token exists", func(t *testing.T) {
		res, _ := testutils.SendTestRequest(
			t, ts, http.MethodPost, "/", http.NoBody, testutils.AuthHeaders(t),
		)
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
