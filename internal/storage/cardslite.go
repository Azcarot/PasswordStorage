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

type BankCardLiteStorage struct {
	Storage PgxStorage
	DB      *sql.DB
	Data    BankCardData
}

var BCLiteS PgxStorage

func (store *BankCardLiteStorage) CreateNewRecord(ctx context.Context) error {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)
	if !ok {
		return ErrNoLogin
	}

	tx, err := store.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = store.DB.ExecContext(ctx, `INSERT INTO bank_card 
	(id, card_number, cvc, exp_date, full_name, comment, username, created) 
	values ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT(id) DO UPDATE SET
	id = excluded.id,
	card_number = excluded.card_number,
	cvc = excluded.cvc,
	exp_date = excluded.exp_date,
	full_name = excluded.full_name,
	comment = excluded.comment,
	username = excluded.username,
	created = excluded.created;`,
		store.Data.ID, store.Data.CardNumber, store.Data.Cvc, store.Data.ExpDate, store.Data.FullName, store.Data.Comment, dataLogin, store.Data.Date)
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

func (store *BankCardLiteStorage) GetRecord(ctx context.Context) (any, error) {
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

			rows := store.DB.QueryRowContext(ctx, query, dataLogin, store.Data.ID)

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

func (store *BankCardLiteStorage) UpdateRecord(ctx context.Context) error {
	tx, err := store.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	err = store.CypherBankData(ctx)
	if err != nil {
		return err
	}

	query := `UPDATE bank_card SET
	card_number = $1, cvc = $2 , exp_date = $3, full_name = $4, comment = $5, created = $6
	WHERE id = $7`
	_, err = store.DB.ExecContext(ctx, query,
		store.Data.CardNumber, store.Data.Cvc, store.Data.ExpDate, store.Data.FullName, store.Data.Comment, store.Data.Date, store.Data.ID)
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

func (store *BankCardLiteStorage) DeleteRecord(ctx context.Context) error {
	tx, err := store.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `DELETE FROM bank_card 
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

func (store *BankCardLiteStorage) SearchRecord(ctx context.Context) (any, error) {
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

			rows, err := store.DB.QueryContext(ctx, query, dataLogin)
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

func (store *BankCardLiteStorage) GetAllRecords(ctx context.Context) (any, error) {
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
			query := `SELECT id, card_number, cvc, exp_date, full_name, comment
	FROM bank_card 
	WHERE username = $1
	ORDER BY id DESC`

			rows, err := store.DB.QueryContext(ctx, query, dataLogin)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			for rows.Next() {
				var resp BankCardResponse
				if err := rows.Scan(&store.Data.ID, &store.Data.CardNumber, &store.Data.Cvc, &store.Data.ExpDate, &store.Data.FullName, &store.Data.Comment); err != nil {
					return result, err
				}

				resp.ID = store.Data.ID
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

func (store *BankCardLiteStorage) CypherBankData(ctx context.Context) error {
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

func (store *BankCardLiteStorage) DeCypherBankData(ctx context.Context) error {
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

func (store BankCardLiteStorage) HashDatabaseData(ctx context.Context) (string, error) {
	bankData, err := store.GetAllRecords(ctx)
	if err != nil {
		return "", err
	}
	jsonData, err := json.Marshal(bankData.([]BankCardResponse))
	if err != nil {
		return "", fmt.Errorf("failed to marshal card data: %v", err)
	}

	hash := sha256.Sum256(jsonData)

	hashString := hex.EncodeToString(hash[:])

	return hashString, nil
}

func NewBCLiteStorage(storage PgxStorage, db *sql.DB) *BankCardLiteStorage {
	return &BankCardLiteStorage{
		Storage: storage,
		DB:      db,
	}
}
