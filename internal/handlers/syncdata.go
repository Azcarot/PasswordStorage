package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Azcarot/PasswordStorage/internal/cypher"
	"github.com/Azcarot/PasswordStorage/internal/storage"
)

// SyncBankData Ручка синхронизации банковских данных пользователя, сравнивает хэши
// данных с сервера и с клиента, при разнице в хэшах возвращает серверные данные
func SyncBankData(res http.ResponseWriter, req *http.Request) {
	var userData storage.UserData
	ctx := req.Context()
	dataLogin, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var syncHashes storage.SyncReq
	userData.Login = dataLogin
	data, err := io.ReadAll(req.Body)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &syncHashes)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	serverHash, err := storage.BCST.HashDatabaseData(ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if serverHash == syncHashes.BankCard {
		res.WriteHeader(http.StatusOK)
		return
	}
	allData, err := storage.BCST.GetAllRecords(ctx)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	decData, ok := allData.([]storage.BankCardData)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	for i, card := range decData {
		err = cypher.DeCypherBankData(ctx, &card)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		decData[i] = card
	}

	jsonData, err := json.Marshal(decData)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusAccepted)
	res.Header().Add("Content-Type", "application/json")
	res.Write(jsonData)
}

// SyncTextData Ручка синхронизации текстовых данных пользователя, сравнивает хэши
// данных с сервера и с клиента, при разнице в хэшах возвращает серверные данные
func SyncTextData(res http.ResponseWriter, req *http.Request) {
	var userData storage.UserData
	ctx := req.Context()
	dataLogin, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var syncHashes storage.SyncReq
	userData.Login = dataLogin
	data, err := io.ReadAll(req.Body)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &syncHashes)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	serverHash, err := storage.TST.HashDatabaseData(ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if serverHash == syncHashes.TextData {
		res.WriteHeader(http.StatusOK)
		return
	}
	allData, err := storage.TST.GetAllRecords(ctx)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	cyphData, ok := allData.([]storage.TextData)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	for i, text := range cyphData {
		err = cypher.DeCypherTextData(ctx, &text)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		cyphData[i] = text
	}

	jsonData, err := json.Marshal(cyphData)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusAccepted)
	res.Header().Add("Content-Type", "application/json")
	res.Write(jsonData)
}

// SyncLPWData Ручка синхронизации данных пользователя типа логин/пароль, сравнивает хэши
// данных с сервера и с клиента, при разнице в хэшах возвращает серверные данные
func SyncLPWData(res http.ResponseWriter, req *http.Request) {
	var userData storage.UserData
	ctx := req.Context()
	dataLogin, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var syncHashes storage.SyncReq
	userData.Login = dataLogin
	data, err := io.ReadAll(req.Body)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &syncHashes)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	serverHash, err := storage.LPST.HashDatabaseData(ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if serverHash == syncHashes.LoginPw {
		res.WriteHeader(http.StatusOK)
		return
	}
	allData, err := storage.LPST.GetAllRecords(ctx)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	cyphData, ok := allData.([]storage.LoginData)
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
	jsonData, err := json.Marshal(cyphData)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusAccepted)
	res.Header().Add("Content-Type", "application/json")
	res.Write(jsonData)
}

// SyncFileData Ручка синхронизации файловых данных пользователя, сравнивает хэши
// данных с сервера и с клиента, при разнице в хэшах возвращает серверные данные
func SyncFileData(res http.ResponseWriter, req *http.Request) {
	var userData storage.UserData
	ctx := req.Context()
	dataLogin, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var syncHashes storage.SyncReq
	userData.Login = dataLogin
	data, err := io.ReadAll(req.Body)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &syncHashes)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	serverHash, err := storage.FST.HashDatabaseData(ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if serverHash == syncHashes.FileData {
		res.WriteHeader(http.StatusOK)
		return
	}
	allData, err := storage.FST.GetAllRecords(ctx)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	files, ok := allData.([]storage.FileData)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	for i, file := range files {
		err := cypher.DeCypherFileData(ctx, &file)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		files[i] = file
	}

	jsonData, err := json.Marshal(files)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusAccepted)
	res.Header().Add("Content-Type", "application/json")
	res.Write(jsonData)
}
