package health

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpCommon "github.com/llravell/simple-cards/internal/controller/http"
	"github.com/rs/zerolog"
)

type Routes struct {
	healthUC httpCommon.HealthUseCase
	log      zerolog.Logger
}

func NewRoutes(healthUC httpCommon.HealthUseCase, log zerolog.Logger) *Routes {
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
