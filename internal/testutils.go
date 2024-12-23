package testutils

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/llravell/simple-cards/pkg/auth"
	"github.com/stretchr/testify/require"
)

const (
	JWTSecretKey = "secret"
	UserUUID     = "test-uuid"
	UserTokenTTL = 24 * time.Hour
)

func BuildAuthHeader(t *testing.T) string {
	t.Helper()

	jwtToken, err := auth.NewJWTManager(JWTSecretKey).Issue(UserUUID, UserTokenTTL)
	require.NoError(t, err)

	return "Bearer " + jwtToken
}

func AuthHeaders(t *testing.T) map[string]string {
	t.Helper()

	headers := make(map[string]string, 1)
	headers["Authorization"] = BuildAuthHeader(t)

	return headers
}

func SendTestRequest(
	t *testing.T,
	ts *httptest.Server,
	method string,
	path string,
	body io.Reader,
	headers map[string]string,
) (*http.Response, []byte) {
	t.Helper()

	req, err := http.NewRequestWithContext(context.TODO(), method, ts.URL+path, body)
	require.NoError(t, err)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := ts.Client().Do(req)
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	return res, b
}

func ToJSON(t *testing.T, m any) string {
	t.Helper()

	data, err := json.Marshal(m)
	require.NoError(t, err)

	data = append(data, '\n')

	return string(data)
}
