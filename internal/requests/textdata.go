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

// AddTextReq - запрос на создание новой текстовой записи на сервере
func AddTextReq(data storage.TextData) (bool, error) {
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx := context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	var cyphData storage.TextData
	var err error
	cyphData.Text, err = storage.CypherData(ctx, data.Text)

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
	regURL := "http://" + storage.ServURL + "/api/user/text/add"
	req, err := http.NewRequest("POST", regURL, bytes.NewBuffer(jsonData))
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

	if response.StatusCode != http.StatusAccepted && response.StatusCode != http.StatusUnauthorized && response.StatusCode != http.StatusUnprocessableEntity {
		return false, fmt.Errorf("unexpexteced reponse")
	}

	if response.StatusCode == http.StatusUnprocessableEntity {
		return false, nil
	}

	return true, nil
}

// DeleteTextReq - запрос на удаление записи на сервере
func DeleteTextReq(data storage.TextData) (bool, error) {
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx := context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	var cyphData storage.TextData
	var err error
	cyphData.Text, err = storage.CypherData(ctx, data.Text)

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
	regURL := "http://" + storage.ServURL + "/api/user/text/delete"
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

// UpdateTextReq - запрос на обновление текстовой записи на сервере
func UpdateTextReq(data storage.TextData) (bool, error) {
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx := context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	var cyphData storage.TextData
	var err error
	cyphData.Text, err = storage.CypherData(ctx, data.Text)

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
	regURL := "http://" + storage.ServURL + "/api/user/text/update"
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

// SyncTextReq - запрос на синхронизацию данных текстового типа, если хеши клиента и сервера
// не различались, данные не трогаем
func SyncTextReq() (bool, error) {
	var err error
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	storage.SyncClientHashes.TextData, err = storage.TLiteS.HashDatabaseData(ctx)
	if err != nil {
		return false, err
	}
	jsonData, err := json.Marshal(storage.SyncClientHashes)
	if err != nil {
		return false, err
	}
	regURL := "http://" + storage.ServURL + "/api/user/text/sync"
	req, err := http.NewRequest("GET", regURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", storage.AuthToken)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		return true, nil
	}

	if response.StatusCode != http.StatusAccepted && response.StatusCode != http.StatusUnauthorized {
		return false, fmt.Errorf("unexpexteced reponse")
	}

	if response.StatusCode == http.StatusAccepted {
		data, err := io.ReadAll(response.Body)
		if err != nil {
			return false, err
		}
		defer req.Body.Close()
		var respData []storage.TextData
		if err = json.Unmarshal(data, &respData); err != nil {
			return false, err
		}

		for _, txt := range respData {
			storage.TLiteS.AddData(txt)
			err := storage.TLiteS.CreateNewRecord(ctx)
			if err != nil {
				return false, err
			}
		}
		newData, err := storage.TLiteS.GetAllRecords(ctx)

		if err != nil {
			return false, err
		}
		var newTextData []storage.TextData
		for _, text := range newData.([]storage.TextResponse) {
			var data storage.TextData
			data.ID = text.ID
			newTextData = append(newTextData, data)
		}
		excessTexts := compareUnorderedTextSlices(respData, newTextData)
		for text := range excessTexts {
			storage.TLiteS.AddData(text)
			err = storage.TLiteS.DeleteRecord(ctx)
			if err != nil {
				return false, err
			}
		}
		return true, nil

	}
	return false, err
}

func textSliceToMap(slice []storage.TextData) (map[int]storage.TextData, map[int]int) {
	m := make(map[int]storage.TextData)
	c := make(map[int]int)
	for _, p := range slice {
		m[p.ID] = p
		c[p.ID]++
	}
	return m, c
}

// Получаем слайс структур банковских карт, которые есть только на клиенте
func compareUnorderedTextSlices(s, c []storage.TextData) []storage.TextData {
	if len(s) == len(c) {
		return nil
	}
	var exids []storage.TextData
	_, mapSIDs := textSliceToMap(s)
	mapClient, mapCIDS := textSliceToMap(c)

	for k, v := range mapCIDS {
		if mapSIDs[k] != v {
			exids = append(exids, mapClient[k])
		}
	}

	return exids
}
