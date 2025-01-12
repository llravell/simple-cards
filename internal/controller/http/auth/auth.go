package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	httpCommon "github.com/llravell/simple-cards/internal/controller/http"
	"github.com/llravell/simple-cards/internal/entity"
	"github.com/llravell/simple-cards/internal/entity/dto"
	"github.com/rs/zerolog"
)

const userTokenTTL = time.Hour * 3

type Routes struct {
	authUC    httpCommon.AuthUseCase
	log       zerolog.Logger
	validator *validator.Validate
}

func NewRoutes(authUC httpCommon.AuthUseCase, log zerolog.Logger) *Routes {
	return &Routes{
		authUC:    authUC,
		log:       log,
		validator: validator.New(),
	}
}

func (routes *Routes) parseAuthRequest(r *http.Request) (*dto.AuthRequest, error) {
	var requestData dto.AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		return nil, err
	}

	requestData.Login = strings.TrimSpace(requestData.Login)
	requestData.Password = strings.TrimSpace(requestData.Password)

	if err := routes.validator.Struct(requestData); err != nil {
		return nil, err
	}

	return &requestData, nil
}

// Swagger spec:
// @Summary      Register new user
// @Tags         auth
// @Accept       json
// @Param        request body dto.AuthRequest true "User creds"
// @Success      200  {object}  dto.AuthResponse
// @Failure      400  "invalid data"
// @Failure      409  "user with same login already exists"
// @Failure      500  "token building error"
// @Router       /api/user/register [post]
func (routes *Routes) register(w http.ResponseWriter, r *http.Request) {
	requestData, err := routes.parseAuthRequest(r)
	if err != nil {
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

	err = json.NewEncoder(w).Encode(dto.AuthResponse{Token: token})
	if err != nil {
		routes.log.Err(err).Msg("response write has been failed")
	}
}

// Swagger spec:
// @Summary      Verify user creds and login
// @Tags         auth
// @Accept       json
// @Param        request body dto.AuthRequest true "User creds"
// @Success      200  {object}  dto.AuthResponse
// @Failure      400  "invalid data"
// @Failure      401  "verification failed"
// @Failure      500  "token building error"
// @Router       /api/user/login [post]
func (routes *Routes) login(w http.ResponseWriter, r *http.Request) {
	requestData, err := routes.parseAuthRequest(r)
	if err != nil {
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

	err = json.NewEncoder(w).Encode(dto.AuthResponse{Token: token})
	if err != nil {
		routes.log.Err(err).Msg("response write has been failed")
	}
}

func (routes *Routes) Apply(r chi.Router) {
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", routes.register)
		r.Post("/login", routes.login)
	})
}
