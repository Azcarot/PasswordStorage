package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

func (store *TextStorage) CreateNewRecord(ctx context.Context) error {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)
	if !ok {
		return ErrNoLogin
	}
	err := store.CypherTextData(ctx)
	if err != nil {
		return err
	}

	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = store.DB.Exec(ctx, `INSERT INTO text_data 
	(text, comment, username, created) 
	values ($1, $2, $3, $4);`,
		store.Data.Text, store.Data.Comment, dataLogin, store.Data.Date)
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

func (store *TextStorage) GetRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := TextResponse{}

	for {
		select {
		case <-ctx.Done():
			return result, errTimeout
		default:
			query := `SELECT text, comment
	FROM text_data
	WHERE username = $1 AND id = $2`

			rows := store.DB.QueryRow(ctx, query, dataLogin, store.Data.ID)

			if err := rows.Scan(&store.Data.Text, &store.Data.Comment); err != nil {
				return result, err
			}
			err := store.DeCypherTextData(ctx)
			if err != nil {
				return result, err
			}
			result.Text = store.Data.Text
			result.Comment = store.Data.Comment
			return result, nil
		}

	}
}

func (store *TextStorage) UpdateRecord(ctx context.Context) error {
	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return err
	}
	err = store.CypherTextData(ctx)
	if err != nil {
		return err
	}
	query := `UPDATE text_data SET
	text = $1, comment = $2, created = $3
	WHERE id = $4`
	_, err = store.DB.Exec(ctx, query,
		store.Data.Text, store.Data.Comment, store.Data.Date, store.Data.ID)
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

func (store *TextStorage) DeleteRecord(ctx context.Context) error {
	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return err
	}
	query := `DELETE FROM text_data
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

func (store *TextStorage) SearchRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []TextResponse{}

	for {
		select {
		case <-ctx.Done():
			return result, errTimeout
		default:
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
				var resp TextResponse
				myMap := make(map[string]string)
				if err := rows.Scan(&store.Data.Text, &store.Data.Comment); err != nil {
					return result, err
				}
				err := store.DeCypherTextData(ctx)
				if err != nil {
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
	}

}

func (store *TextStorage) GetAllRecords(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []TextResponse{}

	for {
		select {
		case <-ctx.Done():
			return result, errTimeout
		default:
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
				var resp TextResponse
				if err := rows.Scan(&store.Data.Text, &store.Data.Comment); err != nil {
					return result, err
				}
				err := store.DeCypherTextData(ctx)
				if err != nil {
					return result, err
				}
				resp.Text = store.Data.Text
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

func (store *TextStorage) CypherTextData(ctx context.Context) error {
	var err error
	store.Data.Text, err = CypherData(ctx, store.Data.Text)
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

func (store *TextStorage) DeCypherTextData(ctx context.Context) error {
	var err error
	store.Data.Text, err = Dechypher(ctx, store.Data.Text)
	if err != nil {
		return err
	}
	store.Data.Comment, err = Dechypher(ctx, store.Data.Comment)
	if err != nil {
		return err
	}

	return err
}

func (store *TextStorage) HashDatabaseData(ctx context.Context) (string, error) {
	textData, err := store.GetAllRecords(ctx)
	if err != nil {
		return "", err
	}
	jsonData, err := json.Marshal(textData.([]TextResponse))
	if err != nil {
		return "", fmt.Errorf("failed to marshal text data: %v", err)
	}

	hash := sha256.Sum256(jsonData)

	hashString := hex.EncodeToString(hash[:])

	return hashString, nil
}
