package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/llravell/simple-cards/internal/entity"
	"github.com/rs/zerolog"
)

type authUseCase interface {
	VerifyUser(ctx context.Context, login string, password string) (*entity.User, error)
	BuildUserToken(user *entity.User) (string, error)
	RegisterUser(ctx context.Context, login string, password string) (*entity.User, error)
}

type (
	authRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	authResponse struct {
		Token string `json:"token"`
	}
)

type AuthRoutes struct {
	authUC authUseCase
	log    zerolog.Logger
}

func NewAuthRoutes(authUC authUseCase, log zerolog.Logger) *AuthRoutes {
	return &AuthRoutes{
		authUC: authUC,
		log:    log,
	}
}

func (routes *AuthRoutes) register(w http.ResponseWriter, r *http.Request) {
	var requestData authRequest

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	user, err := routes.authUC.RegisterUser(r.Context(), requestData.Login, requestData.Password)
	if err != nil {
		if errors.Is(err, entity.ErrUserConflict) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

		routes.log.Error().Err(err).Msg("register user failed")

		return
	}

	token, err := routes.authUC.BuildUserToken(user)
	if err != nil {
		routes.log.Error().Err(err).Msg("build user token while registration failed")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = json.NewEncoder(w).Encode(authResponse{Token: token})
	if err != nil {
		routes.log.Err(err).Msg("response writing failed")
	}
}

func (routes *AuthRoutes) login(w http.ResponseWriter, r *http.Request) {
	var requestData authRequest

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	user, err := routes.authUC.VerifyUser(r.Context(), requestData.Login, requestData.Password)
	if err != nil {
		routes.log.Error().Err(err).Msg("register user failed")
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	token, err := routes.authUC.BuildUserToken(user)
	if err != nil {
		routes.log.Error().Err(err).Msg("build user token while registration failed")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = json.NewEncoder(w).Encode(authResponse{Token: token})
	if err != nil {
		routes.log.Err(err).Msg("response writing failed")
	}
}

func (routes *AuthRoutes) Apply(r chi.Router) {
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", routes.register)
		r.Post("/login", routes.login)
	})
}
