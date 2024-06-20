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

// AddNewText - ручка добавления новой текстовой записи пользователя
func AddNewText(res http.ResponseWriter, req *http.Request) {
	var userData storage.UserData
	ctx := req.Context()
	dataLogin, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var textData storage.TextData
	userData.Login = dataLogin
	data, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &textData)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	textData.User = userData.Login
	mut := sync.Mutex{}
	mut.Lock()
	defer mut.Unlock()
	textData.Date = time.Now().Format(time.RFC3339)
	err = cypher.CypherTextData(ctx, &textData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	err = storage.TST.AddData(textData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = storage.TST.CreateNewRecord(ctx)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	res.WriteHeader(http.StatusAccepted)
}

// GetText - ручка получения текстовой записи пользователя по id
func GetText(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	_, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var reqData storage.TextData

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

	err = storage.TST.AddData(reqData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	textData, err := storage.TST.GetRecord(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	cyphData, ok := textData.(storage.TextData)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = cypher.DeCypherTextData(ctx, &cyphData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(textData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(result)

}

// UpdateText - ручка для обновления тектовых данных пользователя по id
func UpdateText(res http.ResponseWriter, req *http.Request) {
	var userData storage.UserData
	ctx := req.Context()
	dataLogin, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var textData storage.TextData
	userData.Login = dataLogin
	data, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &textData)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	textData.User = userData.Login
	mut := sync.Mutex{}
	mut.Lock()
	defer mut.Unlock()
	textData.Date = time.Now().Format(time.RFC3339)
	err = storage.TST.AddData(textData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	old, err := storage.TST.GetRecord(ctx)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	oldData, ok := old.(storage.TextData)
	if ok {
		err = cypher.DeCypherTextData(ctx, &oldData)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if textData.Text == "" {
			textData.Text = oldData.Text
		}
		err = cypher.CypherTextData(ctx, &textData)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = storage.TST.AddData(textData)
		if err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}
	err = storage.TST.UpdateRecord(ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusAccepted)
}

// DeleteText - ручка для удаления текстовых данных пользователя по id
func DeleteText(res http.ResponseWriter, req *http.Request) {
	var userData storage.UserData
	ctx := req.Context()
	dataLogin, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var textData storage.TextData
	userData.Login = dataLogin
	data, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &textData)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	textData.User = userData.Login
	mut := sync.Mutex{}
	mut.Lock()
	defer mut.Unlock()
	textData.Date = time.Now().Format(time.RFC3339)
	err = storage.TST.AddData(textData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	err = storage.TST.DeleteRecord(ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

// GetAllTexts - ручка для получения всех текстовых записей пользователя
func GetAllTexts(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	_, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	textData, err := storage.TST.GetAllRecords(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	cyphData, ok := textData.([]storage.TextData)
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
	result, err := json.Marshal(cyphData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(result)

}
