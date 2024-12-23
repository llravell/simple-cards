package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	testutils "github.com/llravell/simple-cards/internal"
	"github.com/llravell/simple-cards/internal/controller/http/auth"
	"github.com/llravell/simple-cards/internal/entity"
	"github.com/llravell/simple-cards/internal/mocks"
	"github.com/llravell/simple-cards/internal/usecase"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserRepository = errors.New("somethig went wrong")
	ErrJWTIssuer      = errors.New("somethig went wrong")
)

func generateAuthBody(t *testing.T) io.Reader {
	t.Helper()

	data, err := json.Marshal(map[string]string{
		"login":    "login",
		"password": "password",
	})
	require.NoError(t, err)

	data = append(data, '\n')

	return bytes.NewReader(data)
}

func generatePasswordHash(t *testing.T) string {
	t.Helper()

	passwordBytes, err := bcrypt.GenerateFromPassword([]byte("password"), 14)
	require.NoError(t, err)

	return string(passwordBytes)
}

//nolint:funlen
func TestAuthRoutes(t *testing.T) {
	repo := mocks.NewMockUserRepository(gomock.NewController(t))
	jwtIssuer := mocks.NewMockJWTIssuer(gomock.NewController(t))
	authUseCase := usecase.NewAuthUseCase(repo, jwtIssuer)
	authRoutes := auth.NewRoutes(authUseCase, zerolog.Nop())
	router := chi.NewRouter()

	authRoutes.Apply(router)

	ts := httptest.NewServer(router)
	defer ts.Close()

	testCases := []struct {
		name         string
		method       string
		path         string
		prepareMocks func()
		body         io.Reader
		expectedCode int
	}{
		{
			name:         "register: send empty body",
			method:       http.MethodPost,
			path:         "/api/user/register",
			prepareMocks: func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "register: user store failed",
			method: http.MethodPost,
			path:   "/api/user/register",
			body:   generateAuthBody(t),
			prepareMocks: func() {
				repo.EXPECT().
					StoreUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, ErrUserRepository)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "register: user already exists",
			method: http.MethodPost,
			path:   "/api/user/register",
			body:   generateAuthBody(t),
			prepareMocks: func() {
				repo.EXPECT().
					StoreUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, entity.ErrUserConflict)
			},
			expectedCode: http.StatusConflict,
		},
		{
			name:   "register: token issuing failed",
			method: http.MethodPost,
			path:   "/api/user/register",
			body:   generateAuthBody(t),
			prepareMocks: func() {
				repo.EXPECT().
					StoreUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&entity.User{}, nil)

				jwtIssuer.EXPECT().
					Issue(gomock.Any(), gomock.Any()).
					Return("", ErrJWTIssuer)
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "register: user creating has been succeeded",
			method: http.MethodPost,
			path:   "/api/user/register",
			body:   generateAuthBody(t),
			prepareMocks: func() {
				repo.EXPECT().
					StoreUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&entity.User{}, nil)

				jwtIssuer.EXPECT().
					Issue(gomock.Any(), gomock.Any()).
					Return("abc", nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "login: send empty body",
			method:       http.MethodPost,
			path:         "/api/user/login",
			prepareMocks: func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "login: user verify failed",
			method: http.MethodPost,
			path:   "/api/user/login",
			body:   generateAuthBody(t),
			prepareMocks: func() {
				repo.EXPECT().
					FindUserByLogin(gomock.Any(), gomock.Any()).
					Return(nil, ErrUserRepository)
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:   "login: token issuing failed",
			method: http.MethodPost,
			path:   "/api/user/login",
			body:   generateAuthBody(t),
			prepareMocks: func() {
				user := &entity.User{
					UUID:     testutils.UserUUID,
					Password: generatePasswordHash(t),
				}

				repo.EXPECT().
					FindUserByLogin(gomock.Any(), gomock.Any()).
					Return(user, nil)

				jwtIssuer.EXPECT().
					Issue(gomock.Any(), gomock.Any()).
					Return("", ErrJWTIssuer)
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "login: success",
			method: http.MethodPost,
			path:   "/api/user/login",
			body:   generateAuthBody(t),
			prepareMocks: func() {
				user := &entity.User{
					UUID:     testutils.UserUUID,
					Password: generatePasswordHash(t),
				}

				repo.EXPECT().
					FindUserByLogin(gomock.Any(), gomock.Any()).
					Return(user, nil)

				jwtIssuer.EXPECT().
					Issue(gomock.Any(), gomock.Any()).
					Return("abc", nil)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepareMocks()

			res, _ := testutils.SendTestRequest(t, ts, tc.method, tc.path, tc.body, map[string]string{})
			defer res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode)
		})
	}
}
