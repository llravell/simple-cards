package modules

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/llravell/simple-cards/internal/controller/http/middleware"
	"github.com/llravell/simple-cards/internal/entity"
	"github.com/rs/zerolog"
)

type modulesUseCase interface {
	GetAllModules(ctx context.Context, userUUID string) ([]*entity.Module, error)
	GetModuleWithCards(ctx context.Context, userUUID string, moduleUUID string) (*entity.ModuleWithCards, error)
	CreateNewModule(ctx context.Context, userUUID string, moduleName string) (*entity.Module, error)
	UpdateModule(ctx context.Context, userUUID string, moduleUUID string, moduleName string) (*entity.Module, error)
	DeleteModule(ctx context.Context, userUUID string, moduleUUID string) error
	ImportModuleFromQuizlet(ctx context.Context, userUUID string, moduleName string, quizletModuleID string)
}

type createOrUpdateModuleRequest struct {
	Name string `json:"name" validate:"required,max=100"`
}

type quizletImportRequest struct {
	ModuleName      string `json:"module_name"       validate:"required,max=100"`
	QuizletModuleID string `json:"quizlet_module_id" validate:"required"`
}

type Routes struct {
	log       zerolog.Logger
	modulesUC modulesUseCase
	validator *validator.Validate
}

func NewRoutes(modulesUC modulesUseCase, log zerolog.Logger) *Routes {
	return &Routes{
		log:       log,
		modulesUC: modulesUC,
		validator: validator.New(),
	}
}

func (routes *Routes) jsonResponse(w http.ResponseWriter, v any) {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		routes.log.Err(err).Msg("response write has been failed")
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
		routes.log.Error().Err(err).Msg("modules fetching failed")

		return
	}

	routes.jsonResponse(w, modules)
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
	var req createOrUpdateModuleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	req.Name = strings.TrimSpace(req.Name)

	if err := routes.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	module, err := routes.modulesUC.CreateNewModule(
		r.Context(),
		middleware.GetUserUUIDFromRequest(r),
		req.Name,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)

	routes.jsonResponse(w, module)
}

// Swagger spec:
// @Summary      Get module with cards
// @Security     UsersAuth
// @Tags         modules
// @Accept       json
// @Param        request body quizletImportRequest true "Import module params"
// @Success      200
// @Failure      400
// @Router       /api/modules/import/quizlet [post]
func (routes *Routes) importModuleFromQuizlet(w http.ResponseWriter, r *http.Request) {
	var req quizletImportRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	req.ModuleName = strings.TrimSpace(req.ModuleName)
	req.QuizletModuleID = strings.TrimSpace(req.QuizletModuleID)

	if err := routes.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	routes.modulesUC.ImportModuleFromQuizlet(
		r.Context(),
		middleware.GetUserUUIDFromRequest(r),
		req.ModuleName,
		req.QuizletModuleID,
	)
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
	var req createOrUpdateModuleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	req.Name = strings.TrimSpace(req.Name)

	if err := routes.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	module, err := routes.modulesUC.UpdateModule(
		r.Context(),
		middleware.GetUserUUIDFromRequest(r),
		r.PathValue("module_uuid"),
		req.Name,
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

	routes.jsonResponse(w, module)
}

// Swagger spec:
// @Summary      Delete module
// @Security     UsersAuth
// @Tags         modules
// @Param        module_uuid path string true "Module UUID"
// @Success      202
// @Failure      500
// @Router       /api/modules/{module_uuid}/ [delete]
func (routes *Routes) deleteModule(w http.ResponseWriter, r *http.Request) {
	err := routes.modulesUC.DeleteModule(
		r.Context(),
		middleware.GetUserUUIDFromRequest(r),
		r.PathValue("module_uuid"),
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
// @Failure      404
// @Failure      500
// @Router       /api/modules/{module_uuid}/ [get]
func (routes *Routes) getModuleWithCards(w http.ResponseWriter, r *http.Request) {
	moduleWithCards, err := routes.modulesUC.GetModuleWithCards(
		r.Context(),
		middleware.GetUserUUIDFromRequest(r),
		r.PathValue("module_uuid"),
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

	routes.jsonResponse(w, moduleWithCards)
}

func (routes *Routes) Apply(r chi.Router) {
	r.Route("/api/modules", func(r chi.Router) {
		r.Get("/", routes.getAllModules)
		r.Post("/", routes.createModule)

		r.Route("/import", func(r chi.Router) {
			r.Post("/quizlet", routes.importModuleFromQuizlet)
		})

		r.Route("/{module_uuid}", func(r chi.Router) {
			r.Get("/", routes.getModuleWithCards)
			r.Put("/", routes.updateModule)
			r.Delete("/", routes.deleteModule)
		})
	})
}
