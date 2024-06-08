package storage

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Azcarot/PasswordStorage/internal/utils"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5"
)

// RegisterRequest - тип для запросов на регистрацию пользователя
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// CreateNewUser - создание нового пользователя на сервере
func (store SQLStore) CreateNewUser(ctx context.Context, data UserData) error {
	encodedPW := utils.ShaData(data.Password, SecretKey)
	for {
		select {
		case <-ctx.Done():
			return errTimeout
		default:
			mut.Lock()
			defer mut.Unlock()
			tx, err := store.DB.Begin(ctx)
			if err != nil {
				return err
			}

			_, err = store.DB.Exec(ctx, `INSERT into users (login, password, created) 
	values ($1, $2, $3);`,
				data.Login, encodedPW, data.Date)

			if err != nil {
				tx.Rollback(ctx)
				return err
			}
			err = tx.Commit(ctx)
			if err != nil {
				tx.Rollback(ctx)
				return err
			}
			return err
		}
	}

}

// CheckUserExists - проверка, зарегистрирован ли пользователь
func (store SQLStore) CheckUserExists(data UserData) (bool, error) {
	ctx := context.Background()
	var login string
	sqlQuery := fmt.Sprintf(`SELECT login FROM users WHERE login = '%s'`, data.Login)
	err := store.DB.QueryRow(ctx, sqlQuery).Scan(&login)

	if errors.Is(err, pgx.ErrNoRows) {

		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil

}

// CheckUserPassword - проверка пароля пользователя
func (store SQLStore) CheckUserPassword(ctx context.Context, data UserData) (bool, error) {
	encodedPw := utils.ShaData(data.Password, SecretKey)
	sqlQuery := fmt.Sprintf(`SELECT login, password FROM users WHERE login = '%s'`, data.Login)
	var login, pw string
	err := store.DB.QueryRow(ctx, sqlQuery).Scan(&login, &pw)
	if err != nil {
		return false, err
	}

	if encodedPw != pw {
		return false, nil
	}
	return true, nil
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
