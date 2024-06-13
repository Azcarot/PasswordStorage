package auth

import (
	"log"

	"github.com/golang-jwt/jwt"
)

// SecretKey - ключ для шифрования пользователя
const SecretKey string = "super-secret"

// MyCustomClaims - тип для хранения токена авторизации
type MyCustomClaims struct {
	jwt.MapClaims
}

// VerifyToken - проверка токена авторизации
func VerifyToken(token string) (jwt.MapClaims, bool) {
	hmacSecretString := SecretKey
	hmacSecret := []byte(hmacSecretString)
	gettoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := gettoken.Claims.(jwt.MapClaims); ok && gettoken.Valid {
		return claims, true

	}
	log.Printf("Invalid JWT Token")
	return nil, false

}
