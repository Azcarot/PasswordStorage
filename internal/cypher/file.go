package cypher

import (
	"context"

	"github.com/Azcarot/PasswordStorage/internal/storage"
)

// CypherFileData - шифрование файловых данных пользователя на сервере
func CypherFileData(ctx context.Context, data *storage.FileData) error {
	var err error
	data.FileName, err = CypherData(ctx, data.FileName)
	if err != nil {
		return err
	}
	data.Path, err = CypherData(ctx, data.Path)
	if err != nil {
		return err
	}

	data.Data, err = CypherData(ctx, data.Data)
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

// DeCypherFileData - дешифровка файловых данных на сервере
func DeCypherFileData(ctx context.Context, data *storage.FileData) error {
	var err error

	data.FileName, err = Dechypher(ctx, data.FileName)
	if err != nil {
		return err
	}

	data.Path, err = Dechypher(ctx, data.Path)
	if err != nil {
		return err
	}

	data.Data, err = Dechypher(ctx, data.Data)
	if err != nil {
		return err
	}
	data.Comment, err = Dechypher(ctx, data.Comment)
	if err != nil {
		return err
	}

	return err
}
