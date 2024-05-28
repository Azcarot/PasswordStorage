package storage

import (
	"context"
	"fmt"
	"strings"
)

func (store *BankCardStorage) CreateNewRecord(ctx context.Context) error {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)
	if !ok {
		return ErrNoLogin
	}
	err := store.CypherBankData(ctx)
	if err != nil {
		return err
	}
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

			if err := rows.Scan(&store.Data.CardNumber, &store.Data.Cvc, &store.Data.ExpDate, &store.Data.FullName, &store.Data.Comment); err != nil {
				return result, err
			}
			err := store.DeCypherBankData(ctx)
			if err != nil {
				return result, err
			}
			result.CardNumber = store.Data.CardNumber
			result.Cvc = store.Data.Cvc
			result.ExpDate = store.Data.ExpDate
			result.FullName = store.Data.FullName
			result.Comment = store.Data.Comment
			return result, nil
		}

	}
}

func (store *BankCardStorage) UpdateRecord(ctx context.Context) error {
	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return err
	}
	err = store.CypherBankData(ctx)
	fmt.Println(err)
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
	WHERE username = $1
	ORDER BY id DESC`

			rows, err := store.DB.Query(ctx, query, dataLogin)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			for rows.Next() {
				var resp BankCardResponse
				myMap := make(map[string]string)
				if err := rows.Scan(&store.Data.CardNumber, &store.Data.Cvc, &store.Data.ExpDate, &store.Data.FullName, &store.Data.Comment); err != nil {
					return result, err
				}
				err := store.DeCypherBankData(ctx)
				if err != nil {
					return result, err
				}
				myMap["CardNumber"] = store.Data.CardNumber
				myMap["Cvc"] = store.Data.Cvc
				myMap["ExpDate"] = store.Data.ExpDate
				myMap["FullName"] = store.Data.FullName
				myMap["Comment"] = store.Data.Comment
				for _, value := range myMap {
					if strings.Contains(strings.ToLower(value), strings.ToLower(store.Data.Str)) {
						resp.CardNumber = myMap["CardNumber"]
						resp.Cvc = myMap["Cvc"]
						resp.Comment = myMap["Comment"]
						resp.ExpDate = myMap["ExpDate"]
						resp.FullName = myMap["FullName"]
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
				if err := rows.Scan(&store.Data.CardNumber, &store.Data.Cvc, &store.Data.ExpDate, &store.Data.FullName, &store.Data.Comment); err != nil {
					return result, err
				}
				err := store.DeCypherBankData(ctx)
				if err != nil {
					return result, err
				}
				resp.CardNumber = store.Data.CardNumber
				resp.Cvc = store.Data.Cvc
				resp.ExpDate = store.Data.ExpDate
				resp.FullName = store.Data.FullName
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

func (store *BankCardStorage) CypherBankData(ctx context.Context) error {
	var err error
	store.Data.CardNumber, err = CypherData(ctx, store.Data.CardNumber)
	if err != nil {
		return err
	}
	store.Data.Cvc, err = CypherData(ctx, store.Data.Cvc)
	if err != nil {
		return err
	}
	store.Data.ExpDate, err = CypherData(ctx, store.Data.ExpDate)
	if err != nil {
		return err
	}
	store.Data.FullName, err = CypherData(ctx, store.Data.FullName)
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

func (store *BankCardStorage) DeCypherBankData(ctx context.Context) error {
	var err error
	store.Data.CardNumber, err = Dechypher(ctx, store.Data.CardNumber)
	if err != nil {
		return err
	}
	store.Data.Cvc, err = Dechypher(ctx, store.Data.Cvc)
	if err != nil {
		return err
	}
	store.Data.ExpDate, err = Dechypher(ctx, store.Data.ExpDate)
	if err != nil {
		return err
	}
	store.Data.FullName, err = Dechypher(ctx, store.Data.FullName)
	if err != nil {
		return err
	}
	store.Data.Comment, err = Dechypher(ctx, store.Data.Comment)
	if err != nil {
		return err
	}
	return err
}
