package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestWithLogging(t *testing.T) {
	// Set up a zap test logger
	core, logs := observer.New(zapcore.InfoLevel)
	logger := zap.New(core).Sugar()
	Sugar = *logger

	// Define a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	// Wrap the test handler with the WithLogging middleware
	handler := WithLogging(testHandler)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Hello, World!", rr.Body.String())

	// Check the logs
	entries := logs.All()
	require.Len(t, entries, 1)
	entry := entries[0]
	result := truncateDuration(entry.Message)
	assert.Equal(t, "uri  method GET body <nil> status 200 size 13", result)
}

func truncateDuration(input string) string {

	parts := strings.Split(input, " duration ")
	if len(parts) != 2 {

		return input
	}

	durationParts := strings.SplitN(parts[1], " ", 2)

	if len(durationParts) < 2 {
		return parts[0]
	}

	return parts[0] + " " + durationParts[1]
}
