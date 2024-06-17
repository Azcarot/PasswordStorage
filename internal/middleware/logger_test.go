package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
	fmt.Println(entries)
	require.Len(t, entries, 1)
	entry := entries[0]
	fmt.Println(entry)
	for key, value := range entry.ContextMap() {
		fmt.Println("key ", key, " value ", value)
	}

	assert.Equal(t, "uri  method GET body <nil> status 200 duration 0s size 13", entry.Message)
}
