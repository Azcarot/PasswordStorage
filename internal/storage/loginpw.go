package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// CreateNewRecord - создание новой записи с данными типа логин/пароль на сервере
func (store *LoginPwStorage) CreateNewRecord(ctx context.Context) error {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)
	if !ok {
		return ErrNoLogin
	}

	_, err := store.DB.Exec(ctx, `INSERT INTO login_pw 
	(login, pw, comment, username, created) 
	values ($1, $2, $3, $4, $5);`,
		store.Data.Login, store.Data.Password, store.Data.Comment, dataLogin, store.Data.Date)

	if err != nil {

		return err
	}
	return nil
}

// GetRecord - получение файловых данных типа логин/пароль на сервере по id
func (store *LoginPwStorage) GetRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := LoginData{}

	query := `SELECT login, pw, comment
	FROM login_pw
	WHERE username = $1 AND id = $2`

	rows := store.DB.QueryRow(ctx, query, dataLogin, store.Data.ID)

	if err := rows.Scan(&store.Data.Login, &store.Data.Password, &store.Data.Comment); err != nil {
		return result, err
	}

	result.Login = store.Data.Login
	result.Password = store.Data.Password
	result.Comment = store.Data.Comment
	return result, nil
}

// UpdateRecord - обновление данных типа логин/пароль на сервере по id
func (store *LoginPwStorage) UpdateRecord(ctx context.Context) error {

	query := `UPDATE login_pw SET
	login = $1, pw = $2 , comment = $3, created = $4
	WHERE id = $5`
	_, err := store.DB.Exec(ctx, query,
		store.Data.Login, store.Data.Password, store.Data.Comment, store.Data.Date, store.Data.ID)

	if err != nil {
		return err
	}
	return nil
}

// DeleteRecord - удаление данных типа логин/пароль с сервера по id
func (store *LoginPwStorage) DeleteRecord(ctx context.Context) error {

	query := `DELETE FROM login_pw 
	WHERE id = $1`
	_, err := store.DB.Exec(ctx, query, store.Data.ID)

	if err != nil {

		return err
	}
	return nil
}

// SearchRecord - поиск данных типа логин/пароль на сервере по строке
func (store *LoginPwStorage) SearchRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}

	result := []LoginData{}

	query := `SELECT login, pw, comment
	FROM login_pw 
	WHERE username = $1 
	ORDER BY id DESC`

	rows, err := store.DB.Query(ctx, query, dataLogin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var resp LoginData
		myMap := make(map[string]string)
		if err := rows.Scan(&store.Data.Login, &store.Data.Password, &store.Data.Comment); err != nil {
			return result, err
		}

		myMap["Login"] = store.Data.Login
		myMap["Password"] = store.Data.Password
		myMap["Comment"] = store.Data.Comment
		for _, value := range myMap {
			if strings.Contains(strings.ToLower(value), strings.ToLower(store.Data.Str)) {
				resp.Login = myMap["Login"]
				resp.Password = myMap["Password"]
				resp.Comment = myMap["Comment"]
				result = append(result, resp)
			}
		}

	}
	if err = rows.Err(); err != nil {
		return result, err
	}
	return result, nil
}

// GetAllRecords - получение всех данных типа логин/пароль пользователя на сервере
func (store *LoginPwStorage) GetAllRecords(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []LoginData{}

	query := `SELECT id, login, pw, comment
	FROM login_pw
	WHERE username = $1
	ORDER BY id DESC`

	rows, err := store.DB.Query(ctx, query, dataLogin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var resp LoginData
		if err := rows.Scan(&store.Data.ID, &store.Data.Login, &store.Data.Password, &store.Data.Comment); err != nil {
			return result, err
		}

		resp.Login = store.Data.Login
		resp.ID = store.Data.ID
		resp.Password = store.Data.Password
		resp.Comment = store.Data.Comment
		result = append(result, resp)
	}
	if err = rows.Err(); err != nil {
		return result, err
	}
	return result, nil
}

// HashDatabaseData - получение хэша из всех данных типа логин/пароль пользователя на сервере
func (store LoginPwStorage) HashDatabaseData(ctx context.Context) (string, error) {
	lpwData, err := store.GetAllRecords(ctx)
	if err != nil {
		return "", err
	}
	jsonData, err := json.Marshal(lpwData.([]LoginData))
	if err != nil {
		return "", fmt.Errorf("failed to marshal login/pw data: %v", err)
	}

	hash := sha256.Sum256(jsonData)

	hashString := hex.EncodeToString(hash[:])

	return hashString, nil
}
