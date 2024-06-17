// Package middleware - пакет со всеми миддлвейр (добавление секрета, проверка атвторизации, логгирование запросов)

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Handler to test the middleware
func testHandler(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(storage.EncryptionCtxKey).([16]byte)
	if !ok {
		http.Error(w, "data not found in context", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data[:])
}

func TestAddParamToContext(t *testing.T) {
	tests := []struct {
		name     string
		data     [16]byte
		expected [16]byte
	}{
		{
			name:     "Valid data",
			data:     [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			expected: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		},
		{
			name:     "All zeros",
			data:     [16]byte{},
			expected: [16]byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request and response recorder
			req, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			// Create the middleware with the test data
			middleware := AddParamToContext(tt.data)

			// Serve the request with the middleware and test handler
			handler := middleware(http.HandlerFunc(testHandler))
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, http.StatusOK, rr.Code)

			// Check the response body
			assert.Equal(t, tt.expected[:], rr.Body.Bytes())
		})
	}
}
