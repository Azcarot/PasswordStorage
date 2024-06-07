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

type FileLiteStorage struct {
	Storage PgxStorage
	DB      *sql.DB
	Data    FileData
}

var FLiteS PgxStorage

func (store *FileLiteStorage) CreateNewRecord(ctx context.Context) error {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)
	if !ok {
		return ErrNoLogin
	}

	tx, err := store.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = store.DB.ExecContext(ctx, `INSERT INTO file_data 
	(id, file_name, file_path, data, comment, username, created) 
	values ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT(id) DO UPDATE SET
	id = excluded.id,
	file_name = excluded.file_name,
	file_path = excluded.file_path,
	data = excluded.data,
	comment = excluded.comment,
	username = excluded.username,
	created = excluded.created;`,
		store.Data.ID, store.Data.FileName, store.Data.Path, store.Data.Data, store.Data.Comment, dataLogin, store.Data.Date)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (store *FileLiteStorage) GetRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := FileResponse{}

	for {
		select {
		case <-ctx.Done():
			return result, errTimeout
		default:
			query := `SELECT file_name, file_path, data, comment
	FROM file_data
	WHERE username = $1 AND id = $2`

			rows := store.DB.QueryRowContext(ctx, query, dataLogin, store.Data.ID)

			if err := rows.Scan(&store.Data.FileName, &store.Data.Path, &store.Data.Data, &store.Data.Comment); err != nil {
				return result, err
			}
			err := store.DeCypherFileData(ctx)
			if err != nil {
				return result, err
			}
			result.FileName = store.Data.FileName
			result.Path = store.Data.Path
			result.Data = store.Data.Data
			result.Comment = store.Data.Comment
			return result, nil
		}

	}
}

func (store *FileLiteStorage) UpdateRecord(ctx context.Context) error {
	tx, err := store.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	err = store.CypherFileData(ctx)
	if err != nil {
		return err
	}
	query := `UPDATE file_data SET
	file_name = $1, file_path =$2, data = $3, comment = $4, created = $5
	WHERE id = $6`
	_, err = store.DB.ExecContext(ctx, query,
		store.Data.FileName, store.Data.Path, store.Data.Data, store.Data.Comment, store.Data.Date, store.Data.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (store *FileLiteStorage) DeleteRecord(ctx context.Context) error {
	tx, err := store.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `DELETE FROM file_data
	WHERE id = $1`
	_, err = store.DB.ExecContext(ctx, query, store.Data.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (store *FileLiteStorage) SearchRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []FileResponse{}

	for {
		select {
		case <-ctx.Done():
			return result, errTimeout
		default:
			query := `SELECT file_name, file_path, data, comment
	FROM file_data
	WHERE username = $1 
	ORDER BY id DESC`

			rows, err := store.DB.QueryContext(ctx, query, dataLogin)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			for rows.Next() {
				var resp FileResponse
				myMap := make(map[string]string)
				if err := rows.Scan(&store.Data.FileName, &store.Data.Path, &store.Data.Data, &store.Data.Comment); err != nil {
					return result, err
				}
				err := store.DeCypherFileData(ctx)
				if err != nil {
					return result, err
				}
				myMap["Filename"] = store.Data.FileName
				myMap["Data"] = store.Data.Data
				myMap["Comment"] = store.Data.Comment
				myMap["Path"] = store.Data.Path
				for _, value := range myMap {
					if strings.Contains(strings.ToLower(value), strings.ToLower(store.Data.Str)) {
						resp.FileName = myMap["Filename"]
						resp.Data = myMap["Data"]
						resp.Comment = myMap["Comment"]
						resp.Path = myMap["Path"]
						result = append(result, resp)
					}
				}

			}
			if err = rows.Err(); err != nil {
				return result, err
			}
			return result, nil
		}
	}

}

func (store *FileLiteStorage) GetAllRecords(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []FileResponse{}

	for {
		select {
		case <-ctx.Done():
			return result, errTimeout
		default:
			query := `SELECT id, file_name, file_path, data, comment
	FROM file_data 
	WHERE username = $1
	ORDER BY id DESC`

			rows, err := store.DB.QueryContext(ctx, query, dataLogin)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			for rows.Next() {
				var resp FileResponse
				if err := rows.Scan(&store.Data.ID, &store.Data.FileName, &store.Data.Path, &store.Data.Data, &store.Data.Comment); err != nil {
					return result, err
				}

				resp.ID = store.Data.ID
				resp.FileName = store.Data.FileName
				resp.Path = store.Data.Path
				resp.Data = store.Data.Data
				resp.Comment = store.Data.Comment
				result = append(result, resp)
			}
			if err = rows.Err(); err != nil {
				return result, err
			}
			return result, nil
		}
	}

}

func (store *FileLiteStorage) CypherFileData(ctx context.Context) error {
	var err error
	store.Data.FileName, err = CypherData(ctx, store.Data.FileName)
	if err != nil {
		return err
	}

	store.Data.Path, err = CypherData(ctx, store.Data.Path)
	if err != nil {
		return err
	}

	store.Data.Data, err = CypherData(ctx, store.Data.Data)
	if err != nil {
		return err
	}

	store.Data.Comment, err = CypherData(ctx, store.Data.Comment)
	if err != nil {
		return err
	}
	store.Data.Str, err = CypherData(ctx, store.Data.Str)
	if err != nil {
		return err
	}
	return err
}

func (store *FileLiteStorage) DeCypherFileData(ctx context.Context) error {
	var err error
	store.Data.FileName, err = Dechypher(ctx, store.Data.FileName)
	if err != nil {
		return err
	}

	store.Data.Path, err = Dechypher(ctx, store.Data.Path)
	if err != nil {
		return err
	}

	store.Data.Data, err = Dechypher(ctx, store.Data.Data)
	if err != nil {
		return err
	}
	store.Data.Comment, err = Dechypher(ctx, store.Data.Comment)
	if err != nil {
		return err
	}

	return err
}

func (store *FileLiteStorage) HashDatabaseData(ctx context.Context) (string, error) {
	fileData, err := store.GetAllRecords(ctx)
	if err != nil {
		return "", err
	}
	jsonData, err := json.Marshal(fileData.([]FileResponse))
	if err != nil {
		return "", fmt.Errorf("failed to marshal file data: %v", err)
	}

	hash := sha256.Sum256(jsonData)

	hashString := hex.EncodeToString(hash[:])

	return hashString, nil
}

func NewFLiteStorage(storage PgxStorage, db *sql.DB) *FileLiteStorage {
	return &FileLiteStorage{
		Storage: storage,
		DB:      db,
	}
}
