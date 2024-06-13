package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// RegisterRequest - тип для запросов на регистрацию пользователя
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// CreateNewUser - создание нового пользователя на сервере
func (store SQLStore) CreateNewUser(ctx context.Context, data UserData) error {

	mut.Lock()
	defer mut.Unlock()

	_, err := store.DB.Exec(ctx, `INSERT into users (login, password, created) 
	values ($1, $2, $3);`,
		data.Login, data.Password, data.Date)

	if err != nil {
		return err
	}

	return nil
}

// CheckUserExists - проверка, зарегистрирован ли пользователь
func (store SQLStore) CheckUserExists(ctx context.Context, data UserData) (bool, error) {

	var login string
	sqlQuery := `SELECT login FROM users WHERE login = $1`
	err := store.DB.QueryRow(ctx, sqlQuery, data.Login).Scan(&login)

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

	sqlQuery := `SELECT login, password FROM users WHERE login = $1`
	var login, pw string
	err := store.DB.QueryRow(ctx, sqlQuery, data.Login).Scan(&login, &pw)
	if err != nil {
		return false, err
	}

	if data.Password != pw {
		return false, nil
	}
	return true, nil
}
