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

type Routes struct {
	healthUC healthUseCase
	log      zerolog.Logger
}

func NewRoutes(healthUC healthUseCase, log zerolog.Logger) *Routes {
	return &Routes{
		healthUC: healthUC,
		log:      log,
	}
}

// Swagger spec:
// @Summary      Check database connection
// @Tags         health
// @Success      200
// @Failure      500
// @Router       /ping [get]
func (routes *Routes) ping(w http.ResponseWriter, r *http.Request) {
	err := routes.healthUC.PingContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (routes *Routes) Apply(r chi.Router) {
	r.Get("/ping", routes.ping)
}
