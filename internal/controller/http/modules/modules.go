package modules

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/llravell/simple-cards/internal/controller/http/middleware"
	"github.com/llravell/simple-cards/internal/entity"
	"github.com/rs/zerolog"
)

type ModulesUseCase interface {
	GetAllModules(ctx context.Context, userUUID string) ([]*entity.Module, error)
	GetModuleWithCards(ctx context.Context, userUUID string, moduleUUID string) (*entity.ModuleWithCards, error)
	CreateNewModule(ctx context.Context, userUUID string, moduleName string) (*entity.Module, error)
	UpdateModule(ctx context.Context, userUUID string, moduleUUID string, moduleName string) (*entity.Module, error)
	DeleteModule(ctx context.Context, userUUID string, moduleUUID string) error
}

type createOrUpdateModuleRequest struct {
	Name string `json:"name"`
}

type Routes struct {
	log       zerolog.Logger
	modulesUC ModulesUseCase
}

func NewRoutes(modulesUC ModulesUseCase, log zerolog.Logger) *Routes {
	return &Routes{
		log:       log,
		modulesUC: modulesUC,
	}
}

// Swagger spec:
// @Summary      Get all user's modules
// @Security     UsersAuth
// @Tags         modules
// @Produce      json
// @Success      200  {array}  entity.Module
// @Failure      500
// @Router       /api/modules/ [get]
func (routes *Routes) getAllModules(w http.ResponseWriter, r *http.Request) {
	userUUID := middleware.GetUserUUIDFromRequest(r)

	modules, err := routes.modulesUC.GetAllModules(r.Context(), userUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = json.NewEncoder(w).Encode(modules)
	if err != nil {
		routes.log.Err(err).Msg("response write has been failed")
	}
}

// Swagger spec:
// @Summary      Create new module
// @Security     UsersAuth
// @Tags         modules
// @Accept       json
// @Produce      json
// @Param        request body createOrUpdateModuleRequest true "Module params"
// @Success      201  {object}  entity.Module
// @Failure      400
// @Failure      500
// @Router       /api/modules/ [post]
func (routes *Routes) createModule(w http.ResponseWriter, r *http.Request) {
	var moduleData createOrUpdateModuleRequest

	if err := json.NewDecoder(r.Body).Decode(&moduleData); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	module, err := routes.modulesUC.CreateNewModule(
		r.Context(),
		middleware.GetUserUUIDFromRequest(r),
		moduleData.Name,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(module)
	if err != nil {
		routes.log.Err(err).Msg("response write has been failed")
	}
}

// Swagger spec:
// @Summary      Update module
// @Security     UsersAuth
// @Tags         modules
// @Accept       json
// @Produce      json
// @Param        module_uuid path string true "Module UUID"
// @Param        request body createOrUpdateModuleRequest true "Module params"
// @Success      200  {object}  entity.Module
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /api/modules/{module_uuid}/ [put]
func (routes *Routes) updateModule(w http.ResponseWriter, r *http.Request) {
	moduleUUID := r.PathValue("module_uuid")
	if moduleUUID == "" {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var moduleData createOrUpdateModuleRequest

	if err := json.NewDecoder(r.Body).Decode(&moduleData); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	module, err := routes.modulesUC.UpdateModule(
		r.Context(),
		middleware.GetUserUUIDFromRequest(r),
		moduleUUID,
		moduleData.Name,
	)
	if err != nil {
		var notFoundErr *entity.ModuleNotFoundError

		if errors.As(err, &notFoundErr) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		routes.log.Error().Err(err).Msg("module creating failed")

		return
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(module)
	if err != nil {
		routes.log.Err(err).Msg("response write has been failed")
	}
}

// Swagger spec:
// @Summary      Delete module
// @Security     UsersAuth
// @Tags         modules
// @Param        module_uuid path string true "Module UUID"
// @Success      202
// @Success      400
// @Failure      500
// @Router       /api/modules/{module_uuid}/ [delete]
func (routes *Routes) deleteModule(w http.ResponseWriter, r *http.Request) {
	moduleUUID := r.PathValue("module_uuid")
	if moduleUUID == "" {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err := routes.modulesUC.DeleteModule(
		r.Context(),
		middleware.GetUserUUIDFromRequest(r),
		moduleUUID,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		routes.log.Error().Err(err).Msg("module deleting failed")

		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// Swagger spec:
// @Summary      Get module with cards
// @Security     UsersAuth
// @Tags         modules
// @Produce      json
// @Param        module_uuid path string true "Module UUID"
// @Success      200  {object}  entity.ModuleWithCards
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /api/modules/{module_uuid}/ [get]
func (routes *Routes) getModuleWithCards(w http.ResponseWriter, r *http.Request) {
	moduleUUID := r.PathValue("module_uuid")
	if moduleUUID == "" {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	moduleWithCards, err := routes.modulesUC.GetModuleWithCards(
		r.Context(),
		middleware.GetUserUUIDFromRequest(r),
		moduleUUID,
	)
	if err != nil {
		var notFoundErr *entity.ModuleNotFoundError

		if errors.As(err, &notFoundErr) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		routes.log.Error().Err(err).Msg("module searching failed")

		return
	}

	err = json.NewEncoder(w).Encode(moduleWithCards)
	if err != nil {
		routes.log.Err(err).Msg("response write has been failed")
	}
}

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
