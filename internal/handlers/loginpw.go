package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/Azcarot/PasswordStorage/internal/cypher"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/jackc/pgx/v5"
)

// AddNewLoginPw - ручка сохранения данных типа логин/пароль
func AddNewLoginPw(res http.ResponseWriter, req *http.Request) {
	var userData storage.UserData
	ctx := req.Context()
	dataLogin, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var loginPw storage.LoginData
	userData.Login = dataLogin
	data, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &loginPw)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	loginPw.User = userData.Login
	mut := sync.Mutex{}
	mut.Lock()
	defer mut.Unlock()
	loginPw.Date = time.Now().Format(time.RFC3339)

	err = cypher.CypherLPWData(ctx, &loginPw)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	err = storage.LPST.AddData(loginPw)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = storage.LPST.CreateNewRecord(ctx)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	res.WriteHeader(http.StatusAccepted)
}

// GetLoginPW - ручка получения записи логин/пароль по id
func GetLoginPW(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	_, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var reqData storage.LoginData

	data, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	if err = json.Unmarshal(data, &reqData); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = storage.LPST.AddData(reqData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	lpwData, err := storage.LPST.GetRecord(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	cyphData, ok := lpwData.(storage.LoginData)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = cypher.DeCypherLPWData(ctx, &cyphData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(lpwData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(result)

}

// UpdateLoginPW - ручка обновления записи типа логин/пароль по id
func UpdateLoginPW(res http.ResponseWriter, req *http.Request) {
	var userData storage.UserData
	ctx := req.Context()
	dataLogin, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var loginPWData storage.LoginData
	userData.Login = dataLogin
	data, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &loginPWData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	loginPWData.User = userData.Login

	mut := sync.Mutex{}
	mut.Lock()
	defer mut.Unlock()
	loginPWData.Date = time.Now().Format(time.RFC3339)
	err = storage.LPST.AddData(loginPWData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	old, err := storage.LPST.GetRecord(ctx)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	oldData, ok := old.(storage.LoginData)
	if ok {
		err = cypher.DeCypherLPWData(ctx, &oldData)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if loginPWData.Login == "" {
			loginPWData.Login = oldData.Login
		}
		if loginPWData.Password == "" {
			loginPWData.Password = oldData.Password
		}
		err = cypher.CypherLPWData(ctx, &loginPWData)
		if err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		err = storage.LPST.AddData(loginPWData)
		if err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}

	err = storage.LPST.UpdateRecord(ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusAccepted)
}

// DeleteLoginPW - ручка удаления записи типа логин/пароль по id
func DeleteLoginPW(res http.ResponseWriter, req *http.Request) {
	var userData storage.UserData
	ctx := req.Context()
	dataLogin, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var loginPWData storage.LoginData
	userData.Login = dataLogin
	data, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &loginPWData)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	loginPWData.User = userData.Login
	mut := sync.Mutex{}
	mut.Lock()
	defer mut.Unlock()
	loginPWData.Date = time.Now().Format(time.RFC3339)
	err = storage.LPST.AddData(loginPWData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	err = storage.LPST.DeleteRecord(ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

// Поиск записи типа логин/пароль по строке
func SearchLoginPW(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	_, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var reqData storage.LoginData

	data, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	if err = json.Unmarshal(data, &reqData); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = storage.LPST.AddData(reqData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	lpwData, err := storage.LPST.SearchRecord(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(lpwData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(result)

}

// GetAllLoginPWs - ручка получения всех данных типа логин/пароль для поьзователя
func GetAllLoginPWs(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	_, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	lpwData, err := storage.LPST.GetAllRecords(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	cyphData, ok := lpwData.([]storage.LoginData)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	for i, lpw := range cyphData {
		err = cypher.DeCypherLPWData(ctx, &lpw)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		cyphData[i] = lpw
	}

	result, err := json.Marshal(cyphData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(result)

}
