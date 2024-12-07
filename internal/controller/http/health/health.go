package health

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type HealthUseCase interface {
	PingContext(ctx context.Context) error
}

type healthRoutes struct {
	healthUC HealthUseCase
	log      zerolog.Logger
}

func NewHealthRoutes(r chi.Router, healthUC HealthUseCase, log zerolog.Logger) {
	routes := &healthRoutes{
		healthUC: healthUC,
		log:      log,
	}

	r.Get("/ping", routes.ping)
}

func (hr *healthRoutes) ping(w http.ResponseWriter, r *http.Request) {
	err := hr.healthUC.PingContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
