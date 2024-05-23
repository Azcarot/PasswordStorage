package storage

import (
	"context"
	"fmt"
	"time"
)

func (store Storage) CreateNewRecord(ctx context.Context) error {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)
	if !ok {
		return ErrNoLogin
	}
	data, ok := store.Data.(BankCardData)
	if !ok {

		return fmt.Errorf("error while asserting data to bank card type")
	}
	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = store.DB.Exec(ctx, `INSERT INTO bank_card 
	(card_number, cvc, exp_date, full_name, comment, user, created) 
	values ($1, $2, $3, $4, $5, $6, $7);`,
		data.CardNumber, data.Cvc, data.ExpDate, data.FullName, data.Comment, dataLogin, time.Now())
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

func (store Storage) GetRecord(ctx context.Context, id int) (any, error) {
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
			query := fmt.Sprintf(`SELECT card_number, cvc, exp_date, full_name, comment
	FROM bank_card
	WHERE customer = '%s' AND id = '%d'
	ORDER BY id DESC`, dataLogin, id)

			rows, err := store.DB.Query(ctx, query)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			if err := rows.Scan(&result.CardNumber, &result.Cvc, &result.ExpDate, &result.ExpDate, &result.Comment); err != nil {
				return result, err
			}

		}
		return result, nil
	}
}

func (store Storage) UpdateRecord(ctx context.Context, id int) error {
	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return err
	}
	data, ok := store.Data.(BankCardData)
	if !ok {

		return fmt.Errorf("error while asserting data to bank card type")
	}
	query := `UPDATE bank_card SET
	card_number = $1, cvc = $2 , exp_date = $3, full_name = $4, comment = $5, created = $6
	WHERE id = $7`
	_, err = store.DB.Exec(ctx, query,
		data.CardNumber, data.Cvc, data.ExpDate, data.FullName, data.Comment, time.Now(), id)
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

func (store Storage) DeleteRecord(ctx context.Context, id int) error {
	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return err
	}
	query := `DELETE FROM bank_card 
	WHERE id = $1`
	_, err = store.DB.Exec(ctx, query, id)
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

func (store Storage) SearchRecord(ctx context.Context, str string) (any, error) {
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
	FROM orders 
	WHERE customer = '$1' and 
	order_number LIKE '$2' OR
    cvc LIKE '$2' OR
	exp_date LIKE '$2' OR
	full_name LIKE '$2' OR
	commet LIKE '$2'
	ORDER BY id DESC`

			rows, err := store.DB.Query(ctx, query, dataLogin, str)
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

func (store Storage) GetAllRecords(ctx context.Context) (any, error) {
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
	FROM orders 
	WHERE customer = '$1' 
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
