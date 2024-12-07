package health

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type healthUseCase interface {
	PingContext(ctx context.Context) error
}

type HealthRoutes struct {
	healthUC healthUseCase
	log      zerolog.Logger
}

func NewHealthRoutes(healthUC healthUseCase, log zerolog.Logger) *HealthRoutes {
	return &HealthRoutes{
		healthUC: healthUC,
		log:      log,
	}
}

func (hr *HealthRoutes) ping(w http.ResponseWriter, r *http.Request) {
	err := hr.healthUC.PingContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (routes *HealthRoutes) Apply(r chi.Router) {
	r.Get("/ping", routes.ping)
}
