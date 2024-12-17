package cards_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	testutils "github.com/llravell/simple-cards/internal"
	"github.com/llravell/simple-cards/internal/controller/http/cards"
	"github.com/llravell/simple-cards/internal/entity"
	"github.com/llravell/simple-cards/internal/mocks"
	"github.com/llravell/simple-cards/internal/usecase"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type testCase struct {
	name         string
	method       string
	path         string
	mock         func()
	body         io.Reader
	expectedCode int
	expectedBody string
}

var testCard = entity.Card{
	UUID:       "card-uuid",
	Term:       "term",
	Meaning:    "meaning",
	ModuleUUID: "module-uuid",
}

func prepareTestServer(
	modulesRepo usecase.ModulesRepository,
	cardsRepo usecase.CardsRepository,
) *httptest.Server {
	modulesUseCase := usecase.NewModulesUseCase(modulesRepo, cardsRepo)
	cardsUseCase := usecase.NewCardsUseCase(cardsRepo)
	router := chi.NewRouter()
	routes := cards.NewRoutes(modulesUseCase, cardsUseCase, zerolog.Nop())

	routes.Apply(router)

	return httptest.NewServer(router)
}

//nolint:funlen
func TestGetCards(t *testing.T) {
	modulesRepo := mocks.NewMockModulesRepository(gomock.NewController(t))
	cardsRepo := mocks.NewMockCardsRepository(gomock.NewController(t))
	ts := prepareTestServer(modulesRepo, cardsRepo)

	defer ts.Close()

	testCases := []testCase{
		{
			name:   "module checking error",
			method: http.MethodGet,
			path:   "/api/modules/module-uuid/cards",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(false, errors.New("boom"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "module checking failed",
			method: http.MethodGet,
			path:   "/api/modules/module-uuid/cards",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(false, nil)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:   "cards fetching error",
			method: http.MethodGet,
			path:   "/api/modules/module-uuid/cards",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(true, nil)

				cardsRepo.EXPECT().
					GetModuleCards(gomock.Any(), "module-uuid").
					Return([]*entity.Card{}, errors.New("boom"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "cards returned successfully",
			method: http.MethodGet,
			path:   "/api/modules/module-uuid/cards",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(true, nil)

				cardsRepo.EXPECT().
					GetModuleCards(gomock.Any(), "module-uuid").
					Return([]*entity.Card{&testCard}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: testutils.ToJSON(t, []entity.Card{testCard}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			res, body := testutils.SendTestRequest(t, ts, tc.method, tc.path, tc.body, map[string]string{})
			defer res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode)

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(body))
			}
		})
	}
}

//nolint:funlen
func TestAddCard(t *testing.T) {
	modulesRepo := mocks.NewMockModulesRepository(gomock.NewController(t))
	cardsRepo := mocks.NewMockCardsRepository(gomock.NewController(t))
	ts := prepareTestServer(modulesRepo, cardsRepo)

	defer ts.Close()

	testCases := []testCase{
		{
			name:   "module checking error",
			method: http.MethodPost,
			path:   "/api/modules/module-uuid/cards",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(false, errors.New("boom"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "module checking failed",
			method: http.MethodPost,
			path:   "/api/modules/module-uuid/cards",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(false, nil)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:   "unexpected format",
			method: http.MethodPost,
			path:   "/api/modules/module-uuid/cards",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(true, nil)
			},
			body:         strings.NewReader("not json"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "send empty term",
			method: http.MethodPost,
			path:   "/api/modules/module-uuid/cards",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(true, nil)
			},
			body: strings.NewReader(testutils.ToJSON(t, map[string]string{
				"term":    "",
				"meaning": "meaning",
			})),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "send empty meaning",
			method: http.MethodPost,
			path:   "/api/modules/module-uuid/cards",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(true, nil)
			},
			body: strings.NewReader(testutils.ToJSON(t, map[string]string{
				"term":    "term",
				"meaning": "",
			})),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "card creating error",
			method: http.MethodPost,
			path:   "/api/modules/module-uuid/cards",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(true, nil)

				cardsRepo.EXPECT().
					CreateCard(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("boom"))
			},
			body: strings.NewReader(testutils.ToJSON(t, map[string]string{
				"term":    "term",
				"meaning": "meaning",
			})),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "card created successfully",
			method: http.MethodPost,
			path:   "/api/modules/module-uuid/cards",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(true, nil)

				cardsRepo.EXPECT().
					CreateCard(gomock.Any(), gomock.Any()).
					Return(&testCard, nil)
			},
			body: strings.NewReader(testutils.ToJSON(t, map[string]string{
				"term":    "term",
				"meaning": "meaning",
			})),
			expectedCode: http.StatusCreated,
			expectedBody: testutils.ToJSON(t, testCard),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			res, body := testutils.SendTestRequest(t, ts, tc.method, tc.path, tc.body, map[string]string{})
			defer res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode)

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(body))
			}
		})
	}
}

//nolint:funlen
func TestUpdateCard(t *testing.T) {
	modulesRepo := mocks.NewMockModulesRepository(gomock.NewController(t))
	cardsRepo := mocks.NewMockCardsRepository(gomock.NewController(t))
	ts := prepareTestServer(modulesRepo, cardsRepo)

	defer ts.Close()

	testCases := []testCase{
		{
			name:   "module checking error",
			method: http.MethodPut,
			path:   "/api/modules/module-uuid/cards/card-uuid",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(false, errors.New("boom"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "module checking failed",
			method: http.MethodPut,
			path:   "/api/modules/module-uuid/cards/card-uuid",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(false, nil)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:   "unexpected format",
			method: http.MethodPut,
			path:   "/api/modules/module-uuid/cards/card-uuid",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(true, nil)
			},
			body:         strings.NewReader("not json"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "send empty term and meaning",
			method: http.MethodPut,
			path:   "/api/modules/module-uuid/cards/card-uuid",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(true, nil)
			},
			body: strings.NewReader(testutils.ToJSON(t, map[string]string{
				"term":    "",
				"meaning": "",
			})),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "card updating error",
			method: http.MethodPut,
			path:   "/api/modules/module-uuid/cards/card-uuid",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(true, nil)

				cardsRepo.EXPECT().
					SaveCard(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("boom"))
			},
			body: strings.NewReader(testutils.ToJSON(t, map[string]string{
				"term":    "term",
				"meaning": "meaning",
			})),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "card not found error",
			method: http.MethodPut,
			path:   "/api/modules/module-uuid/cards/card-uuid",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(true, nil)

				cardsRepo.EXPECT().
					SaveCard(gomock.Any(), gomock.Any()).
					Return(nil, &entity.CardNotFoundError{})
			},
			body: strings.NewReader(testutils.ToJSON(t, map[string]string{
				"term":    "term",
				"meaning": "meaning",
			})),
			expectedCode: http.StatusNotFound,
		},
		{
			name:   "card updated successfully",
			method: http.MethodPut,
			path:   "/api/modules/module-uuid/cards/card-uuid",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(true, nil)

				cardsRepo.EXPECT().
					SaveCard(gomock.Any(), gomock.Any()).
					Return(&testCard, nil)
			},
			body: strings.NewReader(testutils.ToJSON(t, map[string]string{
				"term":    "term",
				"meaning": "meaning",
			})),
			expectedCode: http.StatusOK,
			expectedBody: testutils.ToJSON(t, testCard),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			res, body := testutils.SendTestRequest(t, ts, tc.method, tc.path, tc.body, map[string]string{})
			defer res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode)

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(body))
			}
		})
	}
}

//nolint:funlen
func TestDeleteCard(t *testing.T) {
	modulesRepo := mocks.NewMockModulesRepository(gomock.NewController(t))
	cardsRepo := mocks.NewMockCardsRepository(gomock.NewController(t))
	ts := prepareTestServer(modulesRepo, cardsRepo)

	defer ts.Close()

	testCases := []testCase{
		{
			name:   "module checking error",
			method: http.MethodDelete,
			path:   "/api/modules/module-uuid/cards/card-uuid",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(false, errors.New("boom"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "module checking failed",
			method: http.MethodDelete,
			path:   "/api/modules/module-uuid/cards/card-uuid",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(false, nil)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:   "card deleting error",
			method: http.MethodDelete,
			path:   "/api/modules/module-uuid/cards/card-uuid",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(true, nil)

				cardsRepo.EXPECT().
					DeleteCard(gomock.Any(), "module-uuid", "card-uuid").
					Return(errors.New("boom"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "cards deleted successfully",
			method: http.MethodDelete,
			path:   "/api/modules/module-uuid/cards/card-uuid",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(true, nil)

				cardsRepo.EXPECT().
					DeleteCard(gomock.Any(), "module-uuid", "card-uuid").
					Return(nil)
			},
			expectedCode: http.StatusAccepted,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			res, body := testutils.SendTestRequest(t, ts, tc.method, tc.path, tc.body, map[string]string{})
			defer res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode)

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(body))
			}
		})
	}
}
