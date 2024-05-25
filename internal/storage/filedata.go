package storage

import (
	"context"
	"fmt"
)

func (store *FileStorage) CreateNewRecord(ctx context.Context) error {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)
	if !ok {
		return ErrNoLogin
	}
	fmt.Println(store.Data)

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

			if err := rows.Scan(&result.FileName, &result.Data, &result.Comment); err != nil {
				return result, err
			}
			return result, nil
		}

	}
}

func (store *FileStorage) UpdateRecord(ctx context.Context) error {
	tx, err := store.DB.Begin(ctx)
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
	WHERE username = $1 and 
	file_name LIKE $2 OR
	comment LIKE $2
	ORDER BY id DESC`

			rows, err := store.DB.Query(ctx, query, dataLogin, "%"+store.Data.Str+"%")
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			for rows.Next() {
				var resp FileResponse
				if err := rows.Scan(&resp.FileName, &resp.Data, &resp.Comment); err != nil {
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
				if err := rows.Scan(&resp.FileName, &resp.Data, &resp.Comment); err != nil {
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
