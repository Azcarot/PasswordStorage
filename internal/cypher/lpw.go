package cypher

import (
	"context"

	"github.com/Azcarot/PasswordStorage/internal/storage"
)

// CypherFileData - шифрование данных типа логин/пароль пользователя на сервере
func CypherLPWData(ctx context.Context, data *storage.LoginData) error {
	var err error
	data.Login, err = CypherData(ctx, data.Login)
	if err != nil {
		return err
	}
	data.Password, err = CypherData(ctx, data.Password)
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

// DeCypherFileData - дешифровка данных типа логин/пароль на сервере
func DeCypherLPWData(ctx context.Context, data *storage.LoginData) error {
	var err error
	data.Login, err = Dechypher(ctx, data.Login)
	if err != nil {
		return err
	}
	data.Password, err = Dechypher(ctx, data.Password)
	if err != nil {
		return err
	}

	data.Comment, err = Dechypher(ctx, data.Comment)
	if err != nil {
		return err
	}
	return err
}
