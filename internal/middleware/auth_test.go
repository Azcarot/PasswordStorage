package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azcarot/PasswordStorage/internal/auth"
	mock_storage "github.com/Azcarot/PasswordStorage/internal/mock"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	validToken   = "valid-token"
	invalidToken = "invalid-token"
)

func testAuthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func TestCheckAuthorization(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		user         string
		exists       bool
		expectedCode int
	}{
		{
			name:         "Valid token and user exists",
			token:        validToken,
			user:         "User",
			exists:       true,
			expectedCode: http.StatusOK,
		},
		{
			name:  "Invalid token",
			token: invalidToken,

			expectedCode: http.StatusUnauthorized,
		},
		{
			name:   "Valid token but user does not exist",
			token:  validToken,
			exists: false,

			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mock_storage.NewMockPgxConn(ctrl)
			storage.ST = mock

			if tt.token == validToken {
				if tt.exists {
					mock.EXPECT().CheckUserExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				} else {
					mock.EXPECT().CheckUserExists(gomock.Any(), gomock.Any()).Times(1)
				}
				payload := jwt.MapClaims{
					"sub": "User",
					"exp": 1.718889539e+09,
				}

				// Создаем новый JWT-токен и подписываем его по алгоритму HS256

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
				jwtSecretKey := []byte(auth.SecretKey)
				authToken, err := token.SignedString(jwtSecretKey)
				assert.NoError(t, err)
				tt.token = authToken

			}
			// Mock the CheckUserExists method

			// Create a new request and response recorder
			req, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", tt.token)

			rr := httptest.NewRecorder()

			// Create the middleware and test handler
			middleware := CheckAuthorization(http.HandlerFunc(testAuthHandler))

			// Serve the request with the middleware and test handler
			middleware.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}
