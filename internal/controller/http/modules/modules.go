package modules

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type createOrUpdateModuleRequest struct {
	Name string `json:"name"`
}

type Routes struct {
	log zerolog.Logger
}

func NewRoutes(log zerolog.Logger) *Routes {
	return &Routes{
		log: log,
	}
}

// Swagger spec:
// @Summary      Get all user's modules
// @Tags         modules
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {array}  entity.Module
// @Failure      500
// @Router       /api/modules/ [get]
func (routes *Routes) getAllModules(w http.ResponseWriter, r *http.Request) {}

// Swagger spec:
// @Summary      Create new module
// @Tags         modules
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request body createOrUpdateModuleRequest true "Module params"
// @Success      201  {object}  entity.Module
// @Failure      500
// @Router       /api/modules/ [post]
func (routes *Routes) createModule(w http.ResponseWriter, r *http.Request) {}

// Swagger spec:
// @Summary      Update module
// @Tags         modules
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        module_uuid path string true "Module UUID"
// @Param        request body createOrUpdateModuleRequest true "Module params"
// @Success      200  {object}  entity.Module
// @Failure      500
// @Router       /api/modules/{module_uuid}/ [put]
func (routes *Routes) updateModule(w http.ResponseWriter, r *http.Request) {}

// Swagger spec:
// @Summary      Delete module
// @Tags         modules
// @Security     ApiKeyAuth
// @Param        module_uuid path string true "Module UUID"
// @Success      202
// @Failure      500
// @Router       /api/modules/{module_uuid}/ [delete]
func (routes *Routes) deleteModule(w http.ResponseWriter, r *http.Request) {}

// Swagger spec:
// @Summary      Get module with cards
// @Tags         modules
// @Produce      json
// @Security     ApiKeyAuth
// @Param        module_uuid path string true "Module UUID"
// @Success      200  {object}  entity.ModuleWithCards
// @Failure      500
// @Router       /api/modules/{module_uuid}/ [get]
func (routes *Routes) getModuleWithCards(w http.ResponseWriter, r *http.Request) {}

func (routes *Routes) Apply(r chi.Router) {
	r.Route("/api/modules", func(r chi.Router) {
		r.Get("/", routes.getAllModules)
		r.Post("/", routes.createModule)

		r.Route("/{module_uuid}", func(r chi.Router) {
			r.Get("/", routes.getModuleWithCards)
			r.Put("/", routes.updateModule)
			r.Delete("/", routes.deleteModule)
		})
	})
}
