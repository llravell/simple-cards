package middleware_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func echoHandler(t *testing.T) http.HandlerFunc {
	t.Helper()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "text/plain")

		_, err = w.Write(body)
		require.NoError(t, err)
	})
}

func decompress(t *testing.T, data []byte) string {
	t.Helper()

	buf := bytes.NewBuffer(data)
	gr, err := gzip.NewReader(buf)
	require.NoError(t, err)

	res, err := io.ReadAll(gr)
	require.NoError(t, err)

	err = gr.Close()
	require.NoError(t, err)

	return string(res)
}
