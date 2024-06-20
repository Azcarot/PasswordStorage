package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azcarot/PasswordStorage/internal/auth"
	"github.com/Azcarot/PasswordStorage/internal/cypher"
	"github.com/Azcarot/PasswordStorage/internal/storage"
)

// ReistrationReq - запрос на регистрацию пользователя
func RegistrationReq(data storage.RegisterRequest) (bool, error) {
	if len(data.Login) == 0 || len(data.Password) == 0 {
		return false, fmt.Errorf("wrong login/password data")
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	regURL := storage.ServURL + "/api/user/register"
	req, err := http.NewRequest("POST", regURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request using http.Client
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusConflict {
		return false, fmt.Errorf("unexpexteced reponse")
	}

	if response.StatusCode == http.StatusConflict {
		return false, nil
	}
	ctx := req.Context()
	storage.AuthToken = response.Header.Get("Authorization")
	data.Password = cypher.ShaData(data.Password, auth.SecretKey)
	err = storage.LiteST.CreateNewUser(ctx, data)
	if err != nil {
		return false, fmt.Errorf("sqlite user fail")
	}
	err = storage.LiteST.GetSecretKey(data.Login)
	if err != nil {

		return false, fmt.Errorf("sqlite user fail")
	}
	return true, nil
}
