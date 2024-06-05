package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

func (store *FileStorage) CreateNewRecord(ctx context.Context) error {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)
	if !ok {
		return ErrNoLogin
	}

	err := store.CypherFileData(ctx)
	if err != nil {
		return err
	}

	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = store.DB.Exec(ctx, `INSERT INTO file_data 
	(file_name, data, comment, username, created) 
	values ($1, $2, $3, $4, $5);`,
		store.Data.FileName, store.Data.Data, store.Data.Comment, dataLogin, store.Data.Date)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	return nil
}

func (store *FileStorage) GetRecord(ctx context.Context) (any, error) {
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
			query := `SELECT file_name, data, comment
	FROM file_data
	WHERE username = $1 AND id = $2`

			rows := store.DB.QueryRow(ctx, query, dataLogin, store.Data.ID)

			if err := rows.Scan(&store.Data.FileName, &store.Data.Data, &store.Data.Comment); err != nil {
				return result, err
			}
			err := store.DeCypherFileData(ctx)
			if err != nil {
				return result, err
			}
			result.FileName = store.Data.FileName
			result.Data = store.Data.Data
			result.Comment = store.Data.Comment
			return result, nil
		}

	}
}

func (store *FileStorage) UpdateRecord(ctx context.Context) error {
	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return err
	}
	err = store.CypherFileData(ctx)
	if err != nil {
		return err
	}

	query := `UPDATE file_data SET
	file_name = $1, data = $2 , comment = $3, created = $4
	WHERE id = $5`
	_, err = store.DB.Exec(ctx, query,
		store.Data.FileName, store.Data.Data, store.Data.Comment, store.Data.Date, store.Data.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	return nil
}

func (store *FileStorage) DeleteRecord(ctx context.Context) error {
	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return err
	}

	query := `DELETE FROM file_data 
	WHERE id = $1`
	_, err = store.DB.Exec(ctx, query, store.Data.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	return nil
}

func (store *FileStorage) SearchRecord(ctx context.Context) (any, error) {
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
			query := `SELECT file_name, data, comment
	FROM file_data
	WHERE username = $1 
	ORDER BY id DESC`

			rows, err := store.DB.Query(ctx, query, dataLogin)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			for rows.Next() {
				var resp FileResponse
				myMap := make(map[string]string)
				if err := rows.Scan(&store.Data.FileName, &store.Data.Data, &store.Data.Comment); err != nil {
					return result, err
				}
				err := store.DeCypherFileData(ctx)
				if err != nil {
					return result, err
				}
				myMap["FileName"] = store.Data.FileName
				myMap["Data"] = store.Data.Data
				myMap["Comment"] = store.Data.Comment
				for _, value := range myMap {
					if strings.Contains(strings.ToLower(value), strings.ToLower(store.Data.Str)) {
						resp.FileName = myMap["FileName"]
						resp.Data = myMap["Data"]
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
	}

}

func (store *FileStorage) GetAllRecords(ctx context.Context) (any, error) {
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
			query := `SELECT id, file_name, data, comment
	FROM file_data
	WHERE username = $1
	ORDER BY id DESC`

			rows, err := store.DB.Query(ctx, query, dataLogin)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			for rows.Next() {
				var resp FileResponse
				if err := rows.Scan(&store.Data.ID, &store.Data.FileName, &store.Data.Data, &store.Data.Comment); err != nil {
					return result, err
				}
				err := store.DeCypherFileData(ctx)
				if err != nil {
					return result, err
				}
				resp.FileName = store.Data.FileName
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
	}

}

func (store *FileStorage) CypherFileData(ctx context.Context) error {
	var err error
	store.Data.FileName, err = CypherData(ctx, store.Data.FileName)
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

func (store *FileStorage) DeCypherFileData(ctx context.Context) error {
	var err error
	store.Data.FileName, err = Dechypher(ctx, store.Data.FileName)
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

func (store FileStorage) HashDatabaseData(ctx context.Context) (string, error) {
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
