package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/llravell/simple-cards/internal/controller/http/middleware"
	"github.com/llravell/simple-cards/internal/entity"
	"github.com/rs/zerolog"
)

const userTokenTTL = time.Hour * 3

type authUseCase interface {
	VerifyUser(ctx context.Context, login string, password string) (*entity.User, error)
	BuildUserToken(user *entity.User, ttl time.Duration) (string, error)
	RegisterUser(ctx context.Context, login string, password string) (*entity.User, error)
}

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Routes struct {
	authUC authUseCase
	log    zerolog.Logger
}

func NewRoutes(authUC authUseCase, log zerolog.Logger) *Routes {
	return &Routes{
		authUC: authUC,
		log:    log,
	}
}

func (routes *Routes) register(w http.ResponseWriter, r *http.Request) {
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

	token, err := routes.authUC.BuildUserToken(user, userTokenTTL)
	if err != nil {
		routes.log.Error().Err(err).Msg("build user token while registration failed")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     middleware.TokenCookieName,
		MaxAge:   int(userTokenTTL.Seconds()),
		HttpOnly: true,
		Value:    token,
	})
}

func (routes *Routes) login(w http.ResponseWriter, r *http.Request) {
	var requestData authRequest

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	user, err := routes.authUC.VerifyUser(r.Context(), requestData.Login, requestData.Password)
	if err != nil {
		routes.log.Error().Err(err).Msg("login user failed")
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	token, err := routes.authUC.BuildUserToken(user, userTokenTTL)
	if err != nil {
		routes.log.Error().Err(err).Msg("build user token while registration failed")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     middleware.TokenCookieName,
		MaxAge:   int(userTokenTTL.Seconds()),
		HttpOnly: true,
		Value:    token,
	})
}

func (routes *Routes) Apply(r chi.Router) {
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", routes.register)
		r.Post("/login", routes.login)
	})
}
