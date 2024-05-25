package storage

import (
	"context"
	"fmt"
)

func (store *BankCardStorage) CreateNewRecord(ctx context.Context) error {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)
	if !ok {
		return ErrNoLogin
	}
	fmt.Println(store.Data)

	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = store.DB.Exec(ctx, `INSERT INTO bank_card 
	(card_number, cvc, exp_date, full_name, comment, username, created) 
	values ($1, $2, $3, $4, $5, $6, $7);`,
		store.Data.CardNumber, store.Data.Cvc, store.Data.ExpDate, store.Data.FullName, store.Data.Comment, dataLogin, store.Data.Date)
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

func (store *BankCardStorage) GetRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := BankCardResponse{}

	for {
		select {
		case <-ctx.Done():
			return result, errTimeout
		default:
			query := `SELECT card_number, cvc, exp_date, full_name, comment
	FROM bank_card
	WHERE username = $1 AND id = $2`

			rows := store.DB.QueryRow(ctx, query, dataLogin, store.Data.ID)

			if err := rows.Scan(&result.CardNumber, &result.Cvc, &result.ExpDate, &result.FullName, &result.Comment); err != nil {
				return result, err
			}
			return result, nil
		}

	}
}

func (store *BankCardStorage) UpdateRecord(ctx context.Context) error {
	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return err
	}

	query := `UPDATE bank_card SET
	card_number = $1, cvc = $2 , exp_date = $3, full_name = $4, comment = $5, created = $6
	WHERE id = $7`
	_, err = store.DB.Exec(ctx, query,
		store.Data.CardNumber, store.Data.Cvc, store.Data.ExpDate, store.Data.FullName, store.Data.Comment, store.Data.Date, store.Data.ID)
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

func (store *BankCardStorage) DeleteRecord(ctx context.Context) error {
	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return err
	}
	query := `DELETE FROM bank_card 
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

func (store *BankCardStorage) SearchRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []BankCardResponse{}

	for {
		select {
		case <-ctx.Done():
			return result, errTimeout
		default:
			query := `SELECT card_number, cvc, exp_date, full_name, comment
	FROM bank_card 
	WHERE username = $1 and 
	card_number LIKE $2 OR
    cvc LIKE $2 OR
	exp_date LIKE $2 OR
	full_name LIKE $2 OR
	comment LIKE $2
	ORDER BY id DESC`

			rows, err := store.DB.Query(ctx, query, dataLogin, "%"+store.Data.Str+"%")
			fmt.Println(err)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			for rows.Next() {
				var resp BankCardResponse
				if err := rows.Scan(&resp.CardNumber, &resp.Cvc, &resp.ExpDate, &resp.FullName, &resp.Comment); err != nil {
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

func (store *BankCardStorage) GetAllRecords(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []BankCardResponse{}

	for {
		select {
		case <-ctx.Done():
			return result, errTimeout
		default:
			query := `SELECT card_number, cvc, exp_date, full_name, comment
	FROM bank_card 
	WHERE username = $1
	ORDER BY id DESC`

			rows, err := store.DB.Query(ctx, query, dataLogin)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			for rows.Next() {
				var resp BankCardResponse
				if err := rows.Scan(&resp.CardNumber, &resp.Cvc, &resp.ExpDate, &resp.FullName, &resp.Comment); err != nil {
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
