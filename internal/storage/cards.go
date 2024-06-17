// Package storage - описание и реализация всех методов взаимодействия с хранилищами
// на сервере и клиенте
package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// CreateNewRecord - создание новой записи для банковских карт в бд сервера
func (store *BankCardStorage) CreateNewRecord(ctx context.Context) error {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)
	if !ok {
		return ErrNoLogin
	}

	_, err := store.DB.Exec(ctx, `INSERT INTO bank_card 
	(card_number, cvc, exp_date, full_name, comment, username, created) 
	values ($1, $2, $3, $4, $5, $6, $7);`,
		store.Data.CardNumber, store.Data.Cvc, store.Data.ExpDate, store.Data.FullName, store.Data.Comment, dataLogin, store.Data.Date)

	if err != nil {

		return err
	}
	return nil
}

// GetRecord - получение записи для банковских карт в бд сервера по id
func (store *BankCardStorage) GetRecord(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := BankCardData{}

	query := `SELECT card_number, cvc, exp_date, full_name, comment
	FROM bank_card
	WHERE username = $1 AND id = $2`

	rows := store.DB.QueryRow(ctx, query, dataLogin, store.Data.ID)

	if err := rows.Scan(&store.Data.CardNumber, &store.Data.Cvc, &store.Data.ExpDate, &store.Data.FullName, &store.Data.Comment); err != nil {
		return result, err
	}

	result.CardNumber = store.Data.CardNumber
	result.Cvc = store.Data.Cvc
	result.ExpDate = store.Data.ExpDate
	result.FullName = store.Data.FullName
	result.Comment = store.Data.Comment
	result.ID = store.Data.ID
	return result, nil
}

// UpdateRecord - обновление записи для банковских карт в бд сервера  по ее id
func (store *BankCardStorage) UpdateRecord(ctx context.Context) error {

	query := `UPDATE bank_card SET
	card_number = $1, cvc = $2 , exp_date = $3, full_name = $4, comment = $5, created = $6
	WHERE id = $7`
	_, err := store.DB.Exec(ctx, query,
		store.Data.CardNumber, store.Data.Cvc, store.Data.ExpDate, store.Data.FullName, store.Data.Comment, store.Data.Date, store.Data.ID)

	if err != nil {

		return err
	}
	return nil
}

// DeleteRecord - удаление записи для банковских карт в бд сервера  по id
func (store *BankCardStorage) DeleteRecord(ctx context.Context) error {

	query := `DELETE FROM bank_card 
	WHERE id = $1`
	_, err := store.DB.Exec(ctx, query, store.Data.ID)

	if err != nil {

		return err
	}
	return nil
}

// GetAllRecords - получение всех запрошенных данных пользователя для банковских карт из бд сервера
func (store *BankCardStorage) GetAllRecords(ctx context.Context) (any, error) {
	dataLogin, ok := ctx.Value(UserLoginCtxKey).(string)

	if !ok {
		return nil, ErrNoLogin
	}
	result := []BankCardData{}

	query := `SELECT id, card_number, cvc, exp_date, full_name, comment
	FROM bank_card 
	WHERE username = $1
	ORDER BY id DESC`

	rows, err := store.DB.Query(ctx, query, dataLogin)
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

// HashDataBaseData - получения хэша из данных пользователя для типа банковских карт на сервере
func (store *BankCardStorage) HashDatabaseData(ctx context.Context) (string, error) {
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
