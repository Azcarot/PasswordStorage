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
func (store *TextStorage) CreateNewRecord(ctx context.Context) error {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)
	if !ok {
		return ErrNoLogin
	}

	_, err := store.DB.Exec(ctx, `INSERT INTO text_data 
	(text, comment, username, created) 
	values ($1, $2, $3, $4);`,
		store.Data.Text, store.Data.Comment, dataLogin, store.Data.Date)
	if err != nil {

		return err
	}

	return nil
}

// GetRecord - получение файловых текстовых данных на сервере по id
func (store *TextStorage) GetRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := TextData{}

	query := `SELECT text, comment
	FROM text_data
	WHERE username = $1 AND id = $2`

	rows := store.DB.QueryRow(ctx, query, dataLogin, store.Data.ID)

	if err := rows.Scan(&store.Data.Text, &store.Data.Comment); err != nil {
		return result, err
	}

	result.Text = store.Data.Text
	result.Comment = store.Data.Comment
	return result, nil
}

// UpdateRecord - обновление текстовых данных на сервере по id
func (store *TextStorage) UpdateRecord(ctx context.Context) error {

	query := `UPDATE text_data SET
	text = $1, comment = $2, created = $3
	WHERE id = $4`
	_, err := store.DB.Exec(ctx, query,
		store.Data.Text, store.Data.Comment, store.Data.Date, store.Data.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRecord - удаление текстовых данных с сервера по id
func (store *TextStorage) DeleteRecord(ctx context.Context) error {

	query := `DELETE FROM text_data
	WHERE id = $1`
	_, err := store.DB.Exec(ctx, query, store.Data.ID)
	if err != nil {

		return err
	}

	return nil
}

// SearchRecord - поиск текстовых данных на сервере по строке
func (store *TextStorage) SearchRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []TextData{}

	query := `SELECT text, comment
	FROM text_data
	WHERE username = $1 
	ORDER BY id DESC`

	rows, err := store.DB.Query(ctx, query, dataLogin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var resp TextData
		myMap := make(map[string]string)
		if err := rows.Scan(&store.Data.Text, &store.Data.Comment); err != nil {
			return result, err
		}

		myMap["Text"] = store.Data.Text
		myMap["Comment"] = store.Data.Comment
		for _, value := range myMap {
			if strings.Contains(strings.ToLower(value), strings.ToLower(store.Data.Str)) {
				resp.Text = myMap["Text"]
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

// GetAllRecords - получение всех текстовых данных пользователя на сервере
func (store *TextStorage) GetAllRecords(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []TextData{}

	query := `SELECT id, text, comment
	FROM text_data
	WHERE username = $1
	ORDER BY id DESC`

	rows, err := store.DB.Query(ctx, query, dataLogin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var resp TextData
		if err := rows.Scan(&store.Data.ID, &store.Data.Text, &store.Data.Comment); err != nil {
			return result, err
		}

		resp.ID = store.Data.ID
		resp.Text = store.Data.Text
		resp.Comment = store.Data.Comment
		result = append(result, resp)
	}
	if err = rows.Err(); err != nil {
		return result, err
	}
	return result, nil
}

// HashDatabaseData - получение хэша из всех текстовых данных пользователя на сервере
func (store *TextStorage) HashDatabaseData(ctx context.Context) (string, error) {
	textData, err := store.GetAllRecords(ctx)
	if err != nil {
		return "", err
	}
	jsonData, err := json.Marshal(textData.([]TextData))
	if err != nil {
		return "", fmt.Errorf("failed to marshal text data: %v", err)
	}

	hash := sha256.Sum256(jsonData)

	hashString := hex.EncodeToString(hash[:])

	return hashString, nil
}
