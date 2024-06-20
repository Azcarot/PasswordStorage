package cypher

import (
	"context"

	"github.com/Azcarot/PasswordStorage/internal/storage"
)

// CypherTextData - шифрование текстовых данных пользователя на сервере
func CypherTextData(ctx context.Context, data *storage.TextData) error {
	var err error
	data.Text, err = CypherData(ctx, data.Text)
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

// DeCypherTextData - дешифровка текстовых данных на сервере
func DeCypherTextData(ctx context.Context, data *storage.TextData) error {
	var err error
	data.Text, err = Dechypher(ctx, data.Text)
	if err != nil {
		return err
	}
	data.Comment, err = Dechypher(ctx, data.Comment)
	if err != nil {
		return err
	}

	return err
}
