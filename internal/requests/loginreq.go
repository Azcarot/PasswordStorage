package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azcarot/PasswordStorage/internal/handlers"
	"github.com/Azcarot/PasswordStorage/internal/storage"
)

// LoginReq - запрос на авторизацию
func LoginReq(data handlers.LoginRequest) (bool, error) {
	if len(data.Login) == 0 || len(data.Password) == 0 {
		return false, fmt.Errorf("wrong login/password data")
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, err
	}
	regURL := storage.ServURL + "/api/user/login"
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
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusUnauthorized {
		return false, fmt.Errorf("unexpexteced reponse")
	}

	if response.StatusCode == http.StatusUnauthorized {
		return false, nil
	}

	storage.AuthToken = response.Header.Get("Authorization")
	err = storage.LiteST.GetSecretKey(data.Login)
	if err != nil {
		return false, fmt.Errorf("sqlite user fail")
	}
	return true, nil
}
