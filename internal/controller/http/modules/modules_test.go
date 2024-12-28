package modules_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	testutils "github.com/llravell/simple-cards/internal"
	"github.com/llravell/simple-cards/internal/controller/http/modules"
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

var testModule = entity.Module{
	UUID:     "some-uuid",
	Name:     "module for testing",
	UserUUID: "some-user-uuid",
}

var testCard = entity.Card{
	UUID:       "card-uuid",
	Term:       "term",
	Meaning:    "meaning",
	ModuleUUID: "module-uuid",
}

func prepareTestServer(
	t *testing.T,
	modulesRepo usecase.ModulesRepository,
	cardsRepo usecase.CardsRepository,
	quizletModuleParser usecase.QuizletModuleParser,
	quizletImportWP usecase.QuizletImportWorkerPool,
	csvImportWP usecase.CSVImportWorkerPool,
) *httptest.Server {
	t.Helper()

	log := zerolog.Nop()

	modulesUseCase := usecase.NewModulesUseCase(
		modulesRepo,
		cardsRepo,
		quizletModuleParser,
		quizletImportWP,
		csvImportWP,
		&log,
	)
	router := chi.NewRouter()
	routes := modules.NewRoutes(modulesUseCase, log)

	routes.Apply(router)

	return httptest.NewServer(router)
}

func TestGetAllModules(t *testing.T) {
	modulesRepo := mocks.NewMockModulesRepository(gomock.NewController(t))
	cardsRepo := mocks.NewMockCardsRepository(gomock.NewController(t))
	quizletModuleParser := mocks.NewMockQuizletModuleParser(gomock.NewController(t))
	quizletImportWP := mocks.NewMockQuizletImportWorkerPool(gomock.NewController(t))
	csvImportWP := mocks.NewMockCSVImportWorkerPool(gomock.NewController(t))

	ts := prepareTestServer(t, modulesRepo, cardsRepo, quizletModuleParser, quizletImportWP, csvImportWP)

	defer ts.Close()

	testCases := []testCase{
		{
			name: "repo error",
			mock: func() {
				modulesRepo.EXPECT().
					GetAllModules(gomock.Any(), gomock.Any()).
					Return([]*entity.Module{}, errors.New("boom"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "modules returned successfully",
			mock: func() {
				modulesRepo.EXPECT().
					GetAllModules(gomock.Any(), gomock.Any()).
					Return([]*entity.Module{&testModule}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: testutils.ToJSON(t, []entity.Module{testModule}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			res, body := testutils.SendTestRequest(t, ts, http.MethodGet, "/api/modules", tc.body, map[string]string{})
			defer res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode)

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(body))
			}
		})
	}
}

//nolint:funlen
func TestCreateModule(t *testing.T) {
	modulesRepo := mocks.NewMockModulesRepository(gomock.NewController(t))
	cardsRepo := mocks.NewMockCardsRepository(gomock.NewController(t))
	quizletModuleParser := mocks.NewMockQuizletModuleParser(gomock.NewController(t))
	quizletImportWP := mocks.NewMockQuizletImportWorkerPool(gomock.NewController(t))
	csvImportWP := mocks.NewMockCSVImportWorkerPool(gomock.NewController(t))

	ts := prepareTestServer(t, modulesRepo, cardsRepo, quizletModuleParser, quizletImportWP, csvImportWP)

	defer ts.Close()

	testCases := []testCase{
		{
			name:         "unexpected format",
			mock:         func() {},
			body:         strings.NewReader("not json"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "send empty module name",
			mock: func() {},
			body: strings.NewReader(testutils.ToJSON(t, map[string]string{
				"name": "",
			})),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "send module name longer than 100",
			mock: func() {},
			body: strings.NewReader(testutils.ToJSON(t, map[string]string{
				"name": strings.Repeat("a", 101),
			})),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "repo error",
			mock: func() {
				modulesRepo.EXPECT().
					CreateNewModule(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("boom"))
			},
			body: strings.NewReader(testutils.ToJSON(t, map[string]string{
				"name": "module name",
			})),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "created successfully",
			mock: func() {
				modulesRepo.EXPECT().
					CreateNewModule(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&testModule, nil)
			},
			body: strings.NewReader(testutils.ToJSON(t, map[string]string{
				"name": "module name",
			})),
			expectedCode: http.StatusCreated,
			expectedBody: testutils.ToJSON(t, testModule),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			res, body := testutils.SendTestRequest(t, ts, http.MethodPost, "/api/modules", tc.body, map[string]string{})
			defer res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode)

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(body))
			}
		})
	}
}

//nolint:funlen
func TestUpdateModule(t *testing.T) {
	modulesRepo := mocks.NewMockModulesRepository(gomock.NewController(t))
	cardsRepo := mocks.NewMockCardsRepository(gomock.NewController(t))
	quizletModuleParser := mocks.NewMockQuizletModuleParser(gomock.NewController(t))
	quizletImportWP := mocks.NewMockQuizletImportWorkerPool(gomock.NewController(t))
	csvImportWP := mocks.NewMockCSVImportWorkerPool(gomock.NewController(t))

	ts := prepareTestServer(t, modulesRepo, cardsRepo, quizletModuleParser, quizletImportWP, csvImportWP)

	defer ts.Close()

	testCases := []testCase{
		{
			name:         "unexpected format",
			mock:         func() {},
			body:         strings.NewReader("not json"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "send empty module name",
			mock: func() {},
			body: strings.NewReader(testutils.ToJSON(t, map[string]string{
				"name": "",
			})),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "send module name longer than 100",
			mock: func() {},
			body: strings.NewReader(testutils.ToJSON(t, map[string]string{
				"name": strings.Repeat("a", 101),
			})),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "repo error",
			mock: func() {
				modulesRepo.EXPECT().
					UpdateModule(gomock.Any(), gomock.Any(), "module-uuid", gomock.Any()).
					Return(nil, errors.New("boom"))
			},
			body: strings.NewReader(testutils.ToJSON(t, map[string]string{
				"name": "module name",
			})),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "created successfully",
			mock: func() {
				modulesRepo.EXPECT().
					UpdateModule(gomock.Any(), gomock.Any(), "module-uuid", gomock.Any()).
					Return(&testModule, nil)
			},
			body: strings.NewReader(testutils.ToJSON(t, map[string]string{
				"name": "module name",
			})),
			expectedCode: http.StatusOK,
			expectedBody: testutils.ToJSON(t, testModule),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			res, body := testutils.SendTestRequest(
				t, ts, http.MethodPut,
				"/api/modules/module-uuid", tc.body, map[string]string{},
			)
			defer res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode)

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(body))
			}
		})
	}
}

func TestDeleteModule(t *testing.T) {
	modulesRepo := mocks.NewMockModulesRepository(gomock.NewController(t))
	cardsRepo := mocks.NewMockCardsRepository(gomock.NewController(t))
	quizletModuleParser := mocks.NewMockQuizletModuleParser(gomock.NewController(t))
	quizletImportWP := mocks.NewMockQuizletImportWorkerPool(gomock.NewController(t))
	csvImportWP := mocks.NewMockCSVImportWorkerPool(gomock.NewController(t))

	ts := prepareTestServer(t, modulesRepo, cardsRepo, quizletModuleParser, quizletImportWP, csvImportWP)

	defer ts.Close()

	testCases := []testCase{
		{
			name: "repo error",
			mock: func() {
				modulesRepo.EXPECT().
					DeleteModule(gomock.Any(), gomock.Any(), "module-uuid").
					Return(errors.New("boom"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "deleted successfully",
			mock: func() {
				modulesRepo.EXPECT().
					DeleteModule(gomock.Any(), gomock.Any(), "module-uuid").
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
				"/api/modules/module-uuid", tc.body, map[string]string{},
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
func TestGetModuleWithCards(t *testing.T) {
	modulesRepo := mocks.NewMockModulesRepository(gomock.NewController(t))
	cardsRepo := mocks.NewMockCardsRepository(gomock.NewController(t))
	quizletModuleParser := mocks.NewMockQuizletModuleParser(gomock.NewController(t))
	quizletImportWP := mocks.NewMockQuizletImportWorkerPool(gomock.NewController(t))
	csvImportWP := mocks.NewMockCSVImportWorkerPool(gomock.NewController(t))

	ts := prepareTestServer(t, modulesRepo, cardsRepo, quizletModuleParser, quizletImportWP, csvImportWP)

	defer ts.Close()

	moduleWithCards := entity.ModuleWithCards{
		Module: testModule,
		Cards:  []*entity.Card{&testCard},
	}

	testCases := []testCase{
		{
			name: "repo error",
			mock: func() {
				modulesRepo.EXPECT().
					GetModule(gomock.Any(), gomock.Any(), "module-uuid").
					Return(nil, errors.New("boom"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "repo not found error",
			mock: func() {
				modulesRepo.EXPECT().
					GetModule(gomock.Any(), gomock.Any(), "module-uuid").
					Return(nil, &entity.ModuleNotFoundError{})
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "cards repo error",
			mock: func() {
				modulesRepo.EXPECT().
					GetModule(gomock.Any(), gomock.Any(), "module-uuid").
					Return(&testModule, nil)

				cardsRepo.EXPECT().
					GetModuleCards(gomock.Any(), "module-uuid").
					Return([]*entity.Card{}, errors.New("boom"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "module with cards returned successfully",
			mock: func() {
				modulesRepo.EXPECT().
					GetModule(gomock.Any(), gomock.Any(), "module-uuid").
					Return(&testModule, nil)

				cardsRepo.EXPECT().
					GetModuleCards(gomock.Any(), "module-uuid").
					Return([]*entity.Card{&testCard}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: testutils.ToJSON(t, moduleWithCards),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			res, body := testutils.SendTestRequest(
				t, ts, http.MethodGet,
				"/api/modules/module-uuid", tc.body, map[string]string{},
			)
			defer res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode)

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(body))
			}
		})
	}
}
