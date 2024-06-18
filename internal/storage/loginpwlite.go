package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// LPWLiteStorage - хранилище данных типа логин/пароль на клиенте
type LPWLiteStorage struct {
	Storage PgxStorage
	DB      *sql.DB
	Data    LoginData
}

// LPWLiteS - реализация хранилища данных типа логин/пароль на клиенте
var LPWLiteS PgxStorage

// CreateNewRecord - создание новой записи с файловыми данными на клиенте
func (store *LPWLiteStorage) CreateNewRecord(ctx context.Context) error {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)
	if !ok {
		return ErrNoLogin
	}

	_, err := store.DB.ExecContext(ctx, `INSERT INTO login_pw 
	(id, login, pw, comment, username, created) 
	values ($1, $2, $3, $4, $5, $6)
	ON CONFLICT(id) DO UPDATE SET
	id = excluded.id,
	login = excluded.login,
	pw = excluded.pw,
	comment = excluded.comment,
	username = excluded.username,
	created = excluded.created;`,
		store.Data.ID, store.Data.Login, store.Data.Password, store.Data.Comment, dataLogin, store.Data.Date)

	if err != nil {

		return err
	}
	return nil
}

// GetRecord - получение данных типа логин/пароль на клиенте по id
func (store *LPWLiteStorage) GetRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := LoginData{}

	query := `SELECT login, pw, comment
	FROM login_pw
	WHERE username = $1 AND id = $2`

	rows := store.DB.QueryRowContext(ctx, query, dataLogin, store.Data.ID)

	if err := rows.Scan(&store.Data.Login, &store.Data.Password, &store.Data.Comment); err != nil {
		return result, err
	}

	result.Login = store.Data.Login
	result.Password = store.Data.Password
	result.Comment = store.Data.Comment
	return result, nil
}

// UpdateRecord - обновление данных типа логин/пароль на клиенте по id
func (store *LPWLiteStorage) UpdateRecord(ctx context.Context) error {

	query := `UPDATE login_pw SET
	login = $1, pw = $2, comment = $3, created = $4
	WHERE id = $5`
	_, err := store.DB.ExecContext(ctx, query,
		store.Data.Login, store.Data.Password, store.Data.Comment, store.Data.Date, store.Data.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRecord - удаление данных типа логин/пароль с клиента по id
func (store *LPWLiteStorage) DeleteRecord(ctx context.Context) error {

	query := `DELETE FROM login_pw
	WHERE id = $1`
	_, err := store.DB.ExecContext(ctx, query, store.Data.ID)

	if err != nil {

		return err
	}
	return nil
}

// GetAllRecords - получение всех данных типа логин/пароль пользователя на клиенте
func (store *LPWLiteStorage) GetAllRecords(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []LoginData{}

	query := `SELECT id, login, pw, comment
	FROM login_pw 
	WHERE username = $1
	ORDER BY id DESC`

	rows, err := store.DB.QueryContext(ctx, query, dataLogin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var resp LoginData
		if err := rows.Scan(&store.Data.ID, &store.Data.Login, &store.Data.Password, &store.Data.Comment); err != nil {
			return result, err
		}
		resp.ID = store.Data.ID
		resp.Login = store.Data.Login
		resp.Password = store.Data.Password
		resp.Comment = store.Data.Comment
		result = append(result, resp)
	}
	if err = rows.Err(); err != nil {
		return result, err
	}
	return result, nil
}

// HashDatabaseData - получение хэша из всех данных типа логин/пароль пользователя на клиенте
func (store *LPWLiteStorage) HashDatabaseData(ctx context.Context) (string, error) {
	textData, err := store.GetAllRecords(ctx)
	if err != nil {
		return "", err
	}
	jsonData, err := json.Marshal(textData.([]LoginData))
	if err != nil {
		return "", fmt.Errorf("failed to marshal text data: %v", err)
	}

	hash := sha256.Sum256(jsonData)

	hashString := hex.EncodeToString(hash[:])

	return hashString, nil
}

// NewLPLiteStorage - реализация нового хранилища данных типа логин/пароль на клиенте
func NewLPLiteStorage(storage PgxStorage, db *sql.DB) *LPWLiteStorage {
	return &LPWLiteStorage{
		Storage: storage,
		DB:      db,
	}
}
