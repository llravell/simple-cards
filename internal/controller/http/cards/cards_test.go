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
			name: "module checking error",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(false, errors.New("boom"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "module checking failed",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(false, nil)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "cards fetching error",
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
			name: "cards returned successfully",
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

			res, body := testutils.SendTestRequest(
				t, ts, http.MethodGet,
				"/api/modules/module-uuid/cards", tc.body, map[string]string{},
			)
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
			name: "module checking error",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(false, errors.New("boom"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "module checking failed",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(false, nil)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "unexpected format",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(true, nil)
			},
			body:         strings.NewReader("not json"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "send empty term",
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
			name: "send empty meaning",
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
			name: "card creating error",
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
			name: "card created successfully",
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

			res, body := testutils.SendTestRequest(
				t, ts, http.MethodPost,
				"/api/modules/module-uuid/cards", tc.body, map[string]string{},
			)
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
			name: "module checking error",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(false, errors.New("boom"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "module checking failed",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(false, nil)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "unexpected format",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(true, nil)
			},
			body:         strings.NewReader("not json"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "send empty term and meaning",
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
			name: "card updating error",
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
			name: "card not found error",
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
			name: "card updated successfully",
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

			res, body := testutils.SendTestRequest(
				t, ts, http.MethodPut,
				"/api/modules/module-uuid/cards/card-uuid", tc.body, map[string]string{},
			)
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
			name: "module checking error",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(false, errors.New("boom"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "module checking failed",
			mock: func() {
				modulesRepo.EXPECT().
					ModuleExists(gomock.Any(), gomock.Any(), "module-uuid").
					Return(false, nil)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "card deleting error",
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
			name: "cards deleted successfully",
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

			res, body := testutils.SendTestRequest(
				t, ts, http.MethodDelete,
				"/api/modules/module-uuid/cards/card-uuid", tc.body, map[string]string{},
			)
			defer res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode)

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(body))
			}
		})
	}
}
