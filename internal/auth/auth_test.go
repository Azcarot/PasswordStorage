package auth

import (
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestVerifyToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name  string
		args  args
		want  jwt.MapClaims
		want1 bool
	}{
		{name: "1", args: args{token: ""}, want: nil, want1: false},
		{name: "2", args: args{token: "123"}, want: jwt.MapClaims{"exp": 1.718889539e+09, "sub": "User"}, want1: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.args.token) != 0 {

				payload := jwt.MapClaims{
					"sub": "User",
					"exp": tt.want["exp"],
				}

				// Создаем новый JWT-токен и подписываем его по алгоритму HS256
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
				jwtSecretKey := []byte(SecretKey)
				authToken, err := token.SignedString(jwtSecretKey)
				assert.NoError(t, err)
				tt.args.token = authToken
			}
			got, got1 := VerifyToken(tt.args.token)
			if !reflect.DeepEqual(got, tt.want) {

				t.Errorf("VerifyToken() got = %v, want %v", got, tt.want)

			}
			if got1 != tt.want1 {
				t.Errorf("VerifyToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
