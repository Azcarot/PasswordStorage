package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// CreateNewRecord - создание новой записи с файловыми данными на сервере
func (store *FileStorage) CreateNewRecord(ctx context.Context) error {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)
	if !ok {
		return ErrNoLogin
	}

	_, err := store.DB.Exec(ctx, `INSERT INTO file_data 
	(file_name, file_path, data, comment, username, created) 
	values ($1, $2, $3, $4, $5, $6);`,
		store.Data.FileName, store.Data.Path, store.Data.Data, store.Data.Comment, dataLogin, store.Data.Date)

	if err != nil {

		return err
	}
	return nil
}

// GetRecord - получение файловых данных на сервере по id
func (store *FileStorage) GetRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := FileData{}

	query := `SELECT file_name, file_path, data, comment
	FROM file_data
	WHERE username = $1 AND id = $2`

	rows := store.DB.QueryRow(ctx, query, dataLogin, store.Data.ID)

	if err := rows.Scan(&store.Data.FileName, &store.Data.Path, &store.Data.Data, &store.Data.Comment); err != nil {
		return result, err
	}

	result.FileName = store.Data.FileName
	result.Path = store.Data.Path
	result.Data = store.Data.Data
	result.Comment = store.Data.Comment
	return result, nil
}

// UpdateRecord - обновление файловых данных на сервере по id
func (store *FileStorage) UpdateRecord(ctx context.Context) error {

	query := `UPDATE file_data SET
	file_name = $1, file_path = $2, data = $3 , comment = $4, created = $5
	WHERE id = $6`
	_, err := store.DB.Exec(ctx, query,
		store.Data.FileName, store.Data.Path, store.Data.Data, store.Data.Comment, store.Data.Date, store.Data.ID)

	if err != nil {

		return err
	}
	return nil
}

// DeleteRecord - удаление фаловых данных с сервера по id
func (store *FileStorage) DeleteRecord(ctx context.Context) error {

	query := `DELETE FROM file_data 
	WHERE id = $1`
	_, err := store.DB.Exec(ctx, query, store.Data.ID)

	if err != nil {

		return err
	}
	return nil
}

// GetAllRecords - получение всех файловых данных пользователя на сервере
func (store *FileStorage) GetAllRecords(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []FileData{}

	query := `SELECT id, file_name, file_path, data, comment
	FROM file_data
	WHERE username = $1
	ORDER BY id DESC`

	rows, err := store.DB.Query(ctx, query, dataLogin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var resp FileData
		if err := rows.Scan(&store.Data.ID, &store.Data.FileName, &store.Data.Path, &store.Data.Data, &store.Data.Comment); err != nil {
			return result, err
		}

		resp.FileName = store.Data.FileName
		resp.Path = store.Data.Path
		resp.ID = store.Data.ID
		resp.Data = store.Data.Data
		resp.Comment = store.Data.Comment
		result = append(result, resp)
	}
	if err = rows.Err(); err != nil {
		return result, err
	}
	return result, nil
}

// HashDatabaseData - получение хэша из всех файловых данных пользователя на сервере
func (store FileStorage) HashDatabaseData(ctx context.Context) (string, error) {
	fileData, err := store.GetAllRecords(ctx)
	if err != nil {
		return "", err
	}
	jsonData, err := json.Marshal(fileData.([]FileData))
	if err != nil {
		return "", fmt.Errorf("failed to marshal file data: %v", err)
	}

	hash := sha256.Sum256(jsonData)

	hashString := hex.EncodeToString(hash[:])

	return hashString, nil
}
