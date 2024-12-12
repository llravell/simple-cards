package cards

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type createOrUpdateCardRequest struct {
	Term    string `json:"term"`
	Meaning string `json:"meaning"`
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
// @Summary      Get all module's cards
// @Tags         cards
// @Produce      json
// @Security     ApiKeyAuth
// @Param        module_uuid path string true "Module UUID"
// @Success      200  {array}  entity.Card
// @Failure      500
// @Router       /api/modules/{module_uuid}/cards/ [get]
func (routes *Routes) getCards(w http.ResponseWriter, r *http.Request) {}

// Swagger spec:
// @Summary      Add new card to module
// @Tags         cards
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request body createOrUpdateCardRequest true "Card params"
// @Success      201  {object}  entity.Card
// @Failure      500
// @Router       /api/modules/{module_uuid}/cards/ [post]
func (routes *Routes) addCard(w http.ResponseWriter, r *http.Request) {}

// Swagger spec:
// @Summary      Update card
// @Tags         cards
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        module_uuid path string true "Module UUID"
// @Param        card_uuid path string true "Card UUID"
// @Param        request body createOrUpdateCardRequest true "Card params"
// @Success      200  {object}  entity.Card
// @Failure      500
// @Router       /api/modules/{module_uuid}/cards/{card_uuid} [put]
func (routes *Routes) updateCard(w http.ResponseWriter, r *http.Request) {}

// Swagger spec:
// @Summary      Delete card
// @Tags         cards
// @Security     ApiKeyAuth
// @Param        module_uuid path string true "Module UUID"
// @Param        card_uuid path string true "Card UUID"
// @Success      202
// @Failure      500
// @Router       /api/modules/{module_uuid}/cards/{card_uuid} [delete]
func (routes *Routes) deleteCard(w http.ResponseWriter, r *http.Request) {}

func (routes *Routes) Apply(r chi.Router) {
	r.Route("/api/modules/{module_uuid}/cards", func(r chi.Router) {
		r.Get("/", routes.getCards)
		r.Post("/", routes.addCard)

		r.Route("/{card_uuid}", func(r chi.Router) {
			r.Put("/", routes.updateCard)
			r.Delete("/", routes.deleteCard)
		})
	})
}
