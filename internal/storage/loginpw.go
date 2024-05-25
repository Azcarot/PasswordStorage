package storage

import (
	"context"
	"fmt"
)

func (store *LoginPwStorage) CreateNewRecord(ctx context.Context) error {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)
	if !ok {
		return ErrNoLogin
	}
	fmt.Println(store.Data)

	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = store.DB.Exec(ctx, `INSERT INTO login_pw 
	(login, pw, comment, username, created) 
	values ($1, $2, $3, $4, $5);`,
		store.Data.Login, store.Data.Password, store.Data.Comment, dataLogin, store.Data.Date)
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

func (store *LoginPwStorage) GetRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := LoginResponse{}

	for {
		select {
		case <-ctx.Done():
			return result, errTimeout
		default:
			query := `SELECT login, pw, comment
	FROM login_pw
	WHERE username = $1 AND id = $2`

			rows := store.DB.QueryRow(ctx, query, dataLogin, store.Data.ID)

			if err := rows.Scan(&result.Login, &result.Password, &result.Comment); err != nil {
				return result, err
			}
			return result, nil
		}

	}
}

func (store *LoginPwStorage) UpdateRecord(ctx context.Context) error {
	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return err
	}

	query := `UPDATE login_pw SET
	login = $1, pw = $2 , comment = $3, created = $4
	WHERE id = $5`
	_, err = store.DB.Exec(ctx, query,
		store.Data.Login, store.Data.Password, store.Data.Comment, store.Data.Date, store.Data.ID)
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

func (store *LoginPwStorage) DeleteRecord(ctx context.Context) error {
	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return err
	}
	query := `DELETE FROM login_pw 
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

func (store *LoginPwStorage) SearchRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []LoginResponse{}

	for {
		select {
		case <-ctx.Done():
			return result, errTimeout
		default:
			query := `SELECT login, pw, comment
	FROM login_pw 
	WHERE username = $1 and 
	login LIKE $2 OR
    pw LIKE $2 OR
	comment LIKE $2
	ORDER BY id DESC`

			rows, err := store.DB.Query(ctx, query, dataLogin, "%"+store.Data.Str+"%")
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			for rows.Next() {
				var resp LoginResponse
				if err := rows.Scan(&resp.Login, &resp.Password, &resp.Comment); err != nil {
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

func (store *LoginPwStorage) GetAllRecords(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []LoginResponse{}

	for {
		select {
		case <-ctx.Done():
			return result, errTimeout
		default:
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
				var resp LoginResponse
				if err := rows.Scan(&resp.Login, &resp.Password, &resp.Comment); err != nil {
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
