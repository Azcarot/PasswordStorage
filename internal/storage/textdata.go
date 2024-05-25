package storage

import (
	"context"
	"fmt"
)

func (store *TextStorage) CreateNewRecord(ctx context.Context) error {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)
	if !ok {
		return ErrNoLogin
	}
	fmt.Println(store.Data)

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

			if err := rows.Scan(&result.Text, &result.Comment); err != nil {
				return result, err
			}
			return result, nil
		}

	}
}

func (store *TextStorage) UpdateRecord(ctx context.Context) error {
	tx, err := store.DB.Begin(ctx)
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
	WHERE username = $1 and 
	text LIKE $2 OR
    comment LIKE $2
	ORDER BY id DESC`

			rows, err := store.DB.Query(ctx, query, dataLogin, "%"+store.Data.Str+"%")
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			for rows.Next() {
				var resp TextResponse
				if err := rows.Scan(&resp.Text, &resp.Comment); err != nil {
					return result, err
				}
				result = append(result, resp)
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
				if err := rows.Scan(&resp.Text, &resp.Comment); err != nil {
					return result, err
				}
				result = append(result, resp)
			}
			if err = rows.Err(); err != nil {
				return result, err
			}
			return result, nil
		}
	}

}
