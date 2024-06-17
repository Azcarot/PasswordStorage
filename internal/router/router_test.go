package router

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azcarot/PasswordStorage/internal/cfg"
	"github.com/Azcarot/PasswordStorage/internal/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMakeRouter(t *testing.T) {
	// Setup
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	middleware.Sugar = *logger.Sugar()
	defer logger.Sync()

	// Define a mock flag
	flag := cfg.Flags{
		SecretKey: [16]byte{'s', 'e', 'c', 'r', 'e', 't', 'k', 'e', 'y', '1', '2', '3', '4', '5', '6', '7'},
	}

	// Create the router
	router := MakeRouter(flag)

	// Define test cases
	testCases := []struct {
		method       string
		url          string
		expectedCode int
	}{
		{"POST", "/api/user/register", http.StatusBadRequest},
		{"POST", "/api/user/login", http.StatusBadRequest},
		// Add more test cases as needed
	}

	// Execute test cases
	for _, tc := range testCases {

		req, err := http.NewRequest(tc.method, tc.url, bytes.NewBuffer(([]byte{213})))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Check the status code
		assert.Equal(t, tc.expectedCode, rec.Code)
	}
}
