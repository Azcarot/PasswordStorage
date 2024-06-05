package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/Azcarot/PasswordStorage/internal/storage"
)

func AddCardReq(data storage.BankCardData) (bool, error) {
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx := context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	var cyphData storage.BankCardData
	var err error
	cyphData.CardNumber, err = storage.CypherData(ctx, data.CardNumber)

	if err != nil {
		return false, err
	}
	cyphData.ExpDate, err = storage.CypherData(ctx, data.ExpDate)

	if err != nil {
		return false, err
	}
	cyphData.Cvc, err = storage.CypherData(ctx, data.Cvc)

	if err != nil {
		return false, err
	}
	cyphData.Comment, err = storage.CypherData(ctx, data.Comment)
	if err != nil {
		return false, err
	}
	cyphData.FullName, err = storage.CypherData(ctx, data.FullName)
	if err != nil {
		return false, err
	}

	jsonData, err := json.Marshal(cyphData)
	if err != nil {
		return false, err
	}
	regURL := "http://" + storage.ServURL + "/api/user/card/add"
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

	// Check the response status code
	if response.StatusCode != http.StatusAccepted && response.StatusCode != http.StatusUnauthorized && response.StatusCode != http.StatusUnprocessableEntity {
		return false, fmt.Errorf("unexpexteced reponse")
	}

	if response.StatusCode == http.StatusUnprocessableEntity {
		return false, nil
	}

	return true, nil
}

func UpdateCardReq(data storage.BankCardData) (bool, error) {
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx := context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	var cyphData storage.BankCardData
	var err error
	cyphData.CardNumber, err = storage.CypherData(ctx, data.CardNumber)

	if err != nil {
		return false, err
	}
	cyphData.ExpDate, err = storage.CypherData(ctx, data.ExpDate)

	if err != nil {
		return false, err
	}
	cyphData.Cvc, err = storage.CypherData(ctx, data.Cvc)

	if err != nil {
		return false, err
	}
	cyphData.Comment, err = storage.CypherData(ctx, data.Comment)
	if err != nil {
		return false, err
	}
	cyphData.FullName, err = storage.CypherData(ctx, data.FullName)
	if err != nil {
		return false, err
	}

	cyphData.ID = data.ID

	jsonData, err := json.Marshal(cyphData)
	if err != nil {
		return false, err
	}
	regURL := "http://" + storage.ServURL + "/api/user/card/update"
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

	// Check the response status code
	if response.StatusCode != http.StatusAccepted && response.StatusCode != http.StatusUnauthorized && response.StatusCode != http.StatusUnprocessableEntity {
		return false, fmt.Errorf("unexpexteced reponse")
	}

	if response.StatusCode == http.StatusUnprocessableEntity {
		return false, nil
	}

	return true, nil
}

func SyncCardReq() (bool, error) {
	var err error
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	storage.SyncClientHashes.BankCard, err = storage.BCLiteS.HashDatabaseData(ctx)
	if err != nil {
		return false, err
	}
	jsonData, err := json.Marshal(storage.SyncClientHashes)
	if err != nil {
		return false, err
	}
	regURL := "http://" + storage.ServURL + "/api/user/card/sync"
	req, err := http.NewRequest("GET", regURL, bytes.NewBuffer(jsonData))
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
		var respData []storage.BankCardData
		if err = json.Unmarshal(data, &respData); err != nil {
			return false, err
		}
		for _, card := range respData {
			mut := sync.Mutex{}
			mut.Lock()
			defer mut.Unlock()
			storage.BCLiteS.AddData(card)
			err := storage.BCLiteS.CreateNewRecord(ctx)
			if err != nil {
				return false, err
			}
		}
		return true, nil
	}
	return false, err
}
