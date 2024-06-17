package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// BankCardLiteStorage - тип хранилиае банковских кард на клиенте
type BankCardLiteStorage struct {
	Storage PgxStorage
	DB      *sql.DB
	Data    BankCardData
}

// BCLiteS - реализация хранилища банковских карт на клиенте
var BCLiteS PgxStorage

// CreateNewRecord - создание новой записи для банковских карт в бд клиента
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

// GetRecord - получение банковской карты по id из бд клиента
func (store *BankCardLiteStorage) GetRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}

	result := BankCardData{}

	query := `SELECT card_number, cvc, exp_date, full_name, comment
	FROM bank_card
	WHERE username = $1 AND id = $2`

	rows := store.DB.QueryRowContext(ctx, query, dataLogin, store.Data.ID)

	if err := rows.Scan(&store.Data.CardNumber, &store.Data.Cvc, &store.Data.ExpDate, &store.Data.FullName, &store.Data.Comment); err != nil {
		return result, err
	}

	result.CardNumber = store.Data.CardNumber
	result.Cvc = store.Data.Cvc
	result.ExpDate = store.Data.ExpDate
	result.FullName = store.Data.FullName
	result.Comment = store.Data.Comment
	return result, nil
}

// UpdateRecord - обновление банковской карты в бд клиента по id
func (store *BankCardLiteStorage) UpdateRecord(ctx context.Context) error {
	tx, err := store.DB.BeginTx(ctx, nil)
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

// DeleteRecord - удаление банковской карты из бд клиента по id
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

// GetAllRecords - получение всех данных банковских карт пользователя на клиенте
func (store *BankCardLiteStorage) GetAllRecords(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []BankCardData{}

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
		var resp BankCardData
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

// HashDatabaseData получение хеша из банковских карт пользователя на клиенте
func (store BankCardLiteStorage) HashDatabaseData(ctx context.Context) (string, error) {
	bankData, err := store.GetAllRecords(ctx)
	if err != nil {
		return "", err
	}
	jsonData, err := json.Marshal(bankData.([]BankCardData))
	if err != nil {
		return "", fmt.Errorf("failed to marshal card data: %v", err)
	}

	hash := sha256.Sum256(jsonData)

	hashString := hex.EncodeToString(hash[:])

	return hashString, nil
}

// NewBCLiteStorage - создание хранилища банковских карт на клиенте
func NewBCLiteStorage(storage PgxStorage, db *sql.DB) *BankCardLiteStorage {
	return &BankCardLiteStorage{
		Storage: storage,
		DB:      db,
	}
}
