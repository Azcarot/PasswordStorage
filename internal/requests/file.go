package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Azcarot/PasswordStorage/internal/storage"
)

// AddFileReq - запрос на добавление файла на сервер
func AddFileReq(data storage.FileData) (bool, error) {
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx := context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	var cyphData storage.FileData
	var err error
	cyphData.FileName, err = storage.CypherData(ctx, data.FileName)

	if err != nil {
		return false, err
	}

	cyphData.Path, err = storage.CypherData(ctx, data.Path)

	if err != nil {
		return false, err
	}

	cyphData.Data, err = storage.CypherData(ctx, data.Data)

	if err != nil {
		return false, err
	}

	cyphData.Comment, err = storage.CypherData(ctx, data.Comment)
	if err != nil {
		return false, err
	}

	jsonData, err := json.Marshal(cyphData)
	if err != nil {
		return false, err
	}
	regURL := "http://" + storage.ServURL + "/api/user/file/add"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, regURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", storage.AuthToken)
	// Send the request using http.Client
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusAccepted && response.StatusCode != http.StatusUnauthorized && response.StatusCode != http.StatusUnprocessableEntity {
		return false, fmt.Errorf("unexpexteced reponse")
	}

	if response.StatusCode == http.StatusUnprocessableEntity {
		return false, nil
	}

	return true, nil
}

// UpdateFileReq - запрос на обновление файловых данных на сервере
func UpdateFileReq(data storage.FileData) (bool, error) {
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx := context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	var cyphData storage.FileData
	var err error
	cyphData.FileName, err = storage.CypherData(ctx, data.FileName)

	if err != nil {
		return false, err
	}

	cyphData.Path, err = storage.CypherData(ctx, data.Path)

	if err != nil {
		return false, err
	}

	cyphData.Data, err = storage.CypherData(ctx, data.Data)

	if err != nil {
		return false, err
	}

	cyphData.Comment, err = storage.CypherData(ctx, data.Comment)
	if err != nil {
		return false, err
	}

	cyphData.ID = data.ID

	jsonData, err := json.Marshal(cyphData)
	if err != nil {
		return false, err
	}
	regURL := "http://" + storage.ServURL + "/api/user/file/update"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, regURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", storage.AuthToken)
	// Send the request using http.Client
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusAccepted && response.StatusCode != http.StatusUnauthorized && response.StatusCode != http.StatusUnprocessableEntity {
		return false, fmt.Errorf("unexpexteced reponse")
	}

	if response.StatusCode == http.StatusUnprocessableEntity {
		return false, nil
	}

	return true, nil
}

// DeleteFileReq - запрос на удаление файловых данных на сервере
func DeleteFileReq(data storage.FileData) (bool, error) {
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx := context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	var cyphData storage.FileData
	var err error
	cyphData.FileName, err = storage.CypherData(ctx, data.FileName)

	if err != nil {
		return false, err
	}
	cyphData.Path, err = storage.CypherData(ctx, data.Path)

	if err != nil {
		return false, err
	}
	cyphData.Data, err = storage.CypherData(ctx, data.Data)

	if err != nil {
		return false, err
	}

	cyphData.Comment, err = storage.CypherData(ctx, data.Comment)
	if err != nil {
		return false, err
	}

	cyphData.ID = data.ID

	jsonData, err := json.Marshal(cyphData)
	if err != nil {
		return false, err
	}

	regURL := "http://" + storage.ServURL + "/api/user/file/delete"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, regURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", storage.AuthToken)
	// Send the request using http.Client
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusAccepted && response.StatusCode != http.StatusUnauthorized && response.StatusCode != http.StatusUnprocessableEntity {
		return false, fmt.Errorf("unexpexteced reponse")
	}

	if response.StatusCode == http.StatusUnprocessableEntity {
		return false, nil
	}

	return true, nil
}

// SyncFileReq - запрос на синхронизацию файловых данных, если хеши клиента и сервера
// не различались, данные не трогаем
func SyncFileReq() (bool, error) {
	var err error
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	storage.SyncClientHashes.FileData, err = storage.FLiteS.HashDatabaseData(ctx)
	if err != nil {
		return false, err
	}
	jsonData, err := json.Marshal(storage.SyncClientHashes)
	if err != nil {
		return false, err
	}
	regURL := "http://" + storage.ServURL + "/api/user/file/sync"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, regURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", storage.AuthToken)
	// Send the request using http.Client

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		return true, nil
	}
	// Check the response status code
	if response.StatusCode != http.StatusAccepted && response.StatusCode != http.StatusUnauthorized {
		return false, fmt.Errorf("unexpexteced reponse")
	}

	if response.StatusCode == http.StatusAccepted {
		data, err := io.ReadAll(response.Body)
		if err != nil {
			return false, err
		}
		defer req.Body.Close()
		var respData []storage.FileData
		if err = json.Unmarshal(data, &respData); err != nil {
			return false, err
		}

		for _, file := range respData {

			storage.FLiteS.AddData(file)
			err := storage.FLiteS.CreateNewRecord(ctx)
			if err != nil {
				return false, err
			}
		}
		newData, err := storage.FLiteS.GetAllRecords(ctx)

		if err != nil {
			return false, err
		}
		var newFileData []storage.FileData
		for _, file := range newData.([]storage.FileResponse) {
			var data storage.FileData
			data.ID = file.ID
			newFileData = append(newFileData, data)
		}
		excessFiles := compareUnorderedFileSlices(respData, newFileData)
		for file := range excessFiles {
			storage.FLiteS.AddData(file)
			err = storage.FLiteS.DeleteRecord(ctx)
			if err != nil {
				return false, err
			}
		}
		return true, nil
	}
	return false, err
}

func fileSliceToMap(slice []storage.FileData) (map[int]storage.FileData, map[int]int) {
	m := make(map[int]storage.FileData)
	c := make(map[int]int)
	for _, p := range slice {
		m[p.ID] = p
		c[p.ID]++
	}
	return m, c
}

// Получаем слайс структур file, которые есть только на клиенте
func compareUnorderedFileSlices(s, c []storage.FileData) []storage.FileData {
	if len(s) == len(c) {
		return nil
	}
	var exids []storage.FileData
	_, mapSIDs := fileSliceToMap(s)
	mapClient, mapCIDS := fileSliceToMap(c)

	for k, v := range mapCIDS {
		if mapSIDs[k] != v {
			exids = append(exids, mapClient[k])
		}
	}

	return exids
}
