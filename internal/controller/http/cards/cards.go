package cards

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
	ModuleExists(ctx context.Context, userUUID string, moduleUUID string) (bool, error)
}

type cardsUseCase interface {
	GetModuleCards(ctx context.Context, moduleUUID string) ([]*entity.Card, error)
	CreateCard(ctx context.Context, card *entity.Card) (*entity.Card, error)
	SaveCard(ctx context.Context, card *entity.Card) (*entity.Card, error)
	DeleteCard(ctx context.Context, moduleUUID string, cardUUID string) error
}

type createCardRequest struct {
	Term    string `json:"term"    validate:"required"`
	Meaning string `json:"meaning" validate:"required"`
}

type updateCardRequest struct {
	Term    string `json:"term"    validate:"required_without=Meaning"`
	Meaning string `json:"meaning" validate:"required_without=Term"`
}

type Routes struct {
	log       zerolog.Logger
	modulesUC modulesUseCase
	cardsUC   cardsUseCase
	validator *validator.Validate
}

func NewRoutes(
	modulesUC modulesUseCase,
	cardsUC cardsUseCase,
	log zerolog.Logger,
) *Routes {
	return &Routes{
		log:       log,
		modulesUC: modulesUC,
		cardsUC:   cardsUC,
		validator: validator.New(),
	}
}

func (routes *Routes) checkModuleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		moduleUUID := r.PathValue("module_uuid")
		if moduleUUID == "" {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		isModuleExists, err := routes.modulesUC.ModuleExists(
			r.Context(),
			middleware.GetUserUUIDFromRequest(r),
			moduleUUID,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			routes.log.Error().Err(err).Msg("module checking failed")

			return
		}

		if !isModuleExists {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (routes *Routes) jsonResponse(w http.ResponseWriter, v any) {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		routes.log.Err(err).Msg("response write has been failed")
	}
}

// Swagger spec:
// @Summary      Get all module's cards
// @Security     UsersAuth
// @Tags         cards
// @Produce      json
// @Param        module_uuid path string true "Module UUID"
// @Success      200  {array}  entity.Card
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /api/modules/{module_uuid}/cards/ [get]
func (routes *Routes) getCards(w http.ResponseWriter, r *http.Request) {
	cards, err := routes.cardsUC.GetModuleCards(r.Context(), r.PathValue("module_uuid"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		routes.log.Error().Err(err).Msg("cards fetching failed")

		return
	}

	routes.jsonResponse(w, cards)
}

// Swagger spec:
// @Summary      Add new card to module
// @Security     UsersAuth
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param        module_uuid path string true "Module UUID"
// @Param        request body createCardRequest true "Card params"
// @Success      201  {object}  entity.Card
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /api/modules/{module_uuid}/cards/ [post]
func (routes *Routes) addCard(w http.ResponseWriter, r *http.Request) {
	var req createCardRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	req.Term = strings.TrimSpace(req.Term)
	req.Meaning = strings.TrimSpace(req.Meaning)

	if err := routes.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	card, err := routes.cardsUC.CreateCard(r.Context(), &entity.Card{
		ModuleUUID: r.PathValue("module_uuid"),
		Term:       req.Term,
		Meaning:    req.Meaning,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		routes.log.Error().Err(err).Msg("card creating failed")

		return
	}

	w.WriteHeader(http.StatusCreated)
	routes.jsonResponse(w, card)
}

// Swagger spec:
// @Summary      Update card
// @Security     UsersAuth
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param        module_uuid path string true "Module UUID"
// @Param        card_uuid path string true "Card UUID"
// @Param        request body updateCardRequest true "Card update params"
// @Success      200  {object}  entity.Card
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /api/modules/{module_uuid}/cards/{card_uuid} [put]
func (routes *Routes) updateCard(w http.ResponseWriter, r *http.Request) {
	cardUUID := r.PathValue("card_uuid")
	if cardUUID == "" {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var req updateCardRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	req.Term = strings.TrimSpace(req.Term)
	req.Meaning = strings.TrimSpace(req.Meaning)

	if err := routes.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	card, err := routes.cardsUC.SaveCard(r.Context(), &entity.Card{
		UUID:       cardUUID,
		ModuleUUID: r.PathValue("module_uuid"),
		Term:       req.Term,
		Meaning:    req.Meaning,
	})
	if err != nil {
		var notFoundErr *entity.CardNotFoundError

		if errors.As(err, &notFoundErr) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		routes.log.Error().Err(err).Msg("card updating failed")

		return
	}

	routes.jsonResponse(w, card)
}

// Swagger spec:
// @Summary      Delete card
// @Security     UsersAuth
// @Tags         cards
// @Param        module_uuid path string true "Module UUID"
// @Param        card_uuid path string true "Card UUID"
// @Success      202
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /api/modules/{module_uuid}/cards/{card_uuid} [delete]
func (routes *Routes) deleteCard(w http.ResponseWriter, r *http.Request) {
	cardUUID := r.PathValue("card_uuid")
	if cardUUID == "" {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err := routes.cardsUC.DeleteCard(r.Context(), r.PathValue("module_uuid"), cardUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		routes.log.Error().Err(err).Msg("card deleting failed")

		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (routes *Routes) Apply(r chi.Router) {
	r.Route("/api/modules/{module_uuid}/cards", func(r chi.Router) {
		r.Use(routes.checkModuleMiddleware)

		r.Get("/", routes.getCards)
		r.Post("/", routes.addCard)

		r.Route("/{card_uuid}", func(r chi.Router) {
			r.Put("/", routes.updateCard)
			r.Delete("/", routes.deleteCard)
		})
	})
}
