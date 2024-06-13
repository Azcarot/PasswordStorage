package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// TextLiteStorage - хранилище текстовых данных на клиенте
type TextLiteStorage struct {
	Storage PgxStorage
	DB      *sql.DB
	Data    TextData
}

// TLiteS - реализация хранилища текстовых данных на клиенте
var TLiteS PgxStorage

// CreateNewRecord - создание новой записи текстовых данных на клиенте
func (store *TextLiteStorage) CreateNewRecord(ctx context.Context) error {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)
	if !ok {
		return ErrNoLogin
	}

	_, err := store.DB.ExecContext(ctx, `INSERT INTO text_data 
	(id, text, comment, username, created) 
	values ($1, $2, $3, $4, $5)
	ON CONFLICT(id) DO UPDATE SET
	id = excluded.id,
	text = excluded.text,
	comment = excluded.comment,
	username = excluded.username,
	created = excluded.created;`,
		store.Data.ID, store.Data.Text, store.Data.Comment, dataLogin, store.Data.Date)

	if err != nil {

		return err
	}
	return nil
}

// GetRecord - получение текстовых данных на клиенте по id
func (store *TextLiteStorage) GetRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := TextData{}

	query := `SELECT text, comment
	FROM text_data
	WHERE username = $1 AND id = $2`

	rows := store.DB.QueryRowContext(ctx, query, dataLogin, store.Data.ID)

	if err := rows.Scan(&store.Data.Text, &store.Data.Comment); err != nil {
		return result, err
	}

	result.Text = store.Data.Text
	result.Comment = store.Data.Comment
	return result, nil
}

// UpdateRecord - обновление текстовых данных на клиенте по id
func (store *TextLiteStorage) UpdateRecord(ctx context.Context) error {

	query := `UPDATE text_data SET
	text = $1, comment = $2, created = $3
	WHERE id = $4`
	_, err := store.DB.ExecContext(ctx, query,
		store.Data.Text, store.Data.Comment, store.Data.Date, store.Data.ID)

	if err != nil {

		return err
	}
	return nil
}

// DeleteRecord - удаление текстовых данных с клиента по id
func (store *TextLiteStorage) DeleteRecord(ctx context.Context) error {

	query := `DELETE FROM text_data
	WHERE id = $1`
	_, err := store.DB.ExecContext(ctx, query, store.Data.ID)

	if err != nil {

		return err
	}
	return nil
}

// SearchRecord - поиск текстовых данных на клиенте по строке
func (store *TextLiteStorage) SearchRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []TextData{}

	query := `SELECT text, comment
	FROM text_data
	WHERE username = $1 
	ORDER BY id DESC`

	rows, err := store.DB.QueryContext(ctx, query, dataLogin)
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

// GetAllRecords - получение всех текстовых данных пользователя на клиенте
func (store *TextLiteStorage) GetAllRecords(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []TextData{}

	query := `SELECT id, text, comment
	FROM text_data 
	WHERE username = $1
	ORDER BY id DESC`

	rows, err := store.DB.QueryContext(ctx, query, dataLogin)
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

// HashDatabaseData - получение хэша из всех текстовых данных пользователя на клиенте
func (store *TextLiteStorage) HashDatabaseData(ctx context.Context) (string, error) {
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

// NewTLiteStorage - реализация нового хранилища текстовых данных на клиенте
func NewTLiteStorage(storage PgxStorage, db *sql.DB) *TextLiteStorage {
	return &TextLiteStorage{
		Storage: storage,
		DB:      db,
	}
}
