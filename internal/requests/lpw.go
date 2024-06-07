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

func AddLPWReq(data storage.LoginData) (bool, error) {
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx := context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	var cyphData storage.LoginData
	var err error
	cyphData.Login, err = storage.CypherData(ctx, data.Login)

	if err != nil {
		return false, err
	}
	cyphData.Password, err = storage.CypherData(ctx, data.Password)

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
	regURL := "http://" + storage.ServURL + "/api/user/lpw/add"
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

func UpdateLPWReq(data storage.LoginData) (bool, error) {
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx := context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	var cyphData storage.LoginData
	var err error
	cyphData.Login, err = storage.CypherData(ctx, data.Login)

	if err != nil {
		return false, err
	}
	cyphData.Password, err = storage.CypherData(ctx, data.Password)

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
	regURL := "http://" + storage.ServURL + "/api/user/lpw/update"
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

func DeleteLPWReq(data storage.LoginData) (bool, error) {
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx := context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	var cyphData storage.LoginData
	var err error
	cyphData.Login, err = storage.CypherData(ctx, data.Login)

	if err != nil {
		return false, err
	}
	cyphData.Password, err = storage.CypherData(ctx, data.Password)

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

	regURL := "http://" + storage.ServURL + "/api/user/lpw/delete"
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

func SyncLPWReq() (bool, error) {
	var err error
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	storage.SyncClientHashes.LoginPw, err = storage.LPWLiteS.HashDatabaseData(ctx)
	if err != nil {
		return false, err
	}
	jsonData, err := json.Marshal(storage.SyncClientHashes)
	if err != nil {
		return false, err
	}
	regURL := "http://" + storage.ServURL + "/api/user/lpw/sync"
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
		var respData []storage.LoginData
		if err = json.Unmarshal(data, &respData); err != nil {
			return false, err
		}

		for _, lpw := range respData {

			storage.LPWLiteS.AddData(lpw)
			err := storage.LPWLiteS.CreateNewRecord(ctx)
			if err != nil {
				return false, err
			}
		}
		newData, err := storage.LPWLiteS.GetAllRecords(ctx)

		if err != nil {
			return false, err
		}
		var newLPWData []storage.LoginData
		for _, lpw := range newData.([]storage.LoginResponse) {
			var data storage.LoginData
			data.ID = lpw.ID
			newLPWData = append(newLPWData, data)
		}
		excessLPWs := compareUnorderedLPWSlices(respData, newLPWData)
		for lpw := range excessLPWs {
			storage.LPWLiteS.AddData(lpw)
			err = storage.LPWLiteS.DeleteRecord(ctx)
			if err != nil {
				return false, err
			}
		}
		return true, nil
	}
	return false, err
}

func lpwSliceToMap(slice []storage.LoginData) (map[int]storage.LoginData, map[int]int) {
	m := make(map[int]storage.LoginData)
	c := make(map[int]int)
	for _, p := range slice {
		m[p.ID] = p
		c[p.ID]++
	}
	return m, c
}

// Получаем слайс структур login/pw, которые есть только на клиенте
func compareUnorderedLPWSlices(s, c []storage.LoginData) []storage.LoginData {
	if len(s) == len(c) {
		return nil
	}
	var exids []storage.LoginData
	_, mapSIDs := lpwSliceToMap(s)
	mapClient, mapCIDS := lpwSliceToMap(c)

	for k, v := range mapCIDS {
		if mapSIDs[k] != v {
			exids = append(exids, mapClient[k])
		}
	}

	return exids
}
