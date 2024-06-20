package cypher

import (
	"context"

	"github.com/Azcarot/PasswordStorage/internal/storage"
)

// CypherBankData - шифрование секретом данных типа банковских карт
func CypherBankData(ctx context.Context, data *storage.BankCardData) error {
	var err error
	data.CardNumber, err = CypherData(ctx, data.CardNumber)
	if err != nil {
		return err
	}
	data.Cvc, err = CypherData(ctx, data.Cvc)
	if err != nil {
		return err
	}
	data.ExpDate, err = CypherData(ctx, data.ExpDate)
	if err != nil {
		return err
	}
	data.FullName, err = CypherData(ctx, data.FullName)
	if err != nil {
		return err
	}
	data.Comment, err = CypherData(ctx, data.Comment)
	if err != nil {
		return err
	}
	data.Str, err = CypherData(ctx, data.Str)
	if err != nil {
		return err
	}
	return err
}

// DeCypherBankData - дешифровка данных типа банковских карт
func DeCypherBankData(ctx context.Context, data *storage.BankCardData) error {
	var err error
	data.CardNumber, err = Dechypher(ctx, data.CardNumber)
	if err != nil {
		return err
	}
	data.Cvc, err = Dechypher(ctx, data.Cvc)
	if err != nil {
		return err
	}
	data.ExpDate, err = Dechypher(ctx, data.ExpDate)
	if err != nil {
		return err
	}
	data.FullName, err = Dechypher(ctx, data.FullName)
	if err != nil {
		return err
	}
	data.Comment, err = Dechypher(ctx, data.Comment)
	if err != nil {
		return err
	}
	return err
}
