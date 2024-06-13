// Package handlers - все обработчики запросов на сервере
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

// AddNewCard - ручка для добавления новых данных типа банковских карт
func AddNewCard(res http.ResponseWriter, req *http.Request) {
	var userData storage.UserData
	ctx := req.Context()
	dataLogin, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var bankData storage.BankCardData
	userData.Login = dataLogin
	data, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &bankData)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	bankData.User = userData.Login
	mut := sync.Mutex{}
	mut.Lock()
	defer mut.Unlock()
	bankData.Date = time.Now().Format(time.RFC3339)
	err = cypher.CypherBankData(ctx, &bankData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	err = storage.BCST.AddData(bankData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = storage.BCST.CreateNewRecord(ctx)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	res.WriteHeader(http.StatusAccepted)
}

// GetBankCard - ручка для получения конкретной банковской карты по id
func GetBankCard(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	_, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var reqData storage.BankCardData

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
	err = storage.BCST.AddData(reqData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	cardData, err := storage.BCST.GetRecord(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	bankData := cardData.(storage.BankCardData)
	err = cypher.DeCypherBankData(ctx, &bankData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(bankData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(result)

}

// UpdateCard - ручка для обновления записи банковской карты по id
func UpdateCard(res http.ResponseWriter, req *http.Request) {
	var userData storage.UserData
	ctx := req.Context()
	dataLogin, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var bankData storage.BankCardData
	userData.Login = dataLogin
	data, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &bankData)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	bankData.User = userData.Login
	mut := sync.Mutex{}
	mut.Lock()
	defer mut.Unlock()
	bankData.Date = time.Now().Format(time.RFC3339)
	err = storage.BCST.AddData(bankData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	old, err := storage.BCST.GetRecord(ctx)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	oldData, ok := old.(storage.BankCardData)
	if ok {
		err := cypher.DeCypherBankData(ctx, &oldData)
		if err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		if bankData.CardNumber == "" {
			bankData.CardNumber = oldData.CardNumber
		}
		if bankData.Cvc == "" {
			bankData.Cvc = oldData.Cvc
		}
		if bankData.ExpDate == "" {
			bankData.ExpDate = oldData.ExpDate
		}
		err = cypher.CypherBankData(ctx, &bankData)
		if err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		err = storage.BCST.AddData(bankData)
		if err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}

	err = storage.BCST.UpdateRecord(ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusAccepted)
}

// DeleteCard - ручка для удаления банковской карты по id
func DeleteCard(res http.ResponseWriter, req *http.Request) {
	var userData storage.UserData
	ctx := req.Context()
	dataLogin, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var bankData storage.BankCardData
	userData.Login = dataLogin
	data, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &bankData)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	bankData.User = userData.Login
	mut := sync.Mutex{}
	mut.Lock()
	defer mut.Unlock()
	bankData.Date = time.Now().Format(time.RFC3339)
	err = storage.BCST.AddData(bankData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	err = storage.BCST.DeleteRecord(ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

// SearchBankCard - ручка для поиска банковской карты по str
func SearchBankCard(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	_, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var reqData storage.BankCardData

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
	err = storage.BCST.AddData(reqData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	cardData, err := storage.BCST.SearchRecord(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	cypData := cardData.(storage.BankCardData)
	err = cypher.DeCypherBankData(ctx, &cypData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	result, err := json.Marshal(cypData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(result)

}

// GetAllBankCards - ручка для получения полного списка сохраненных карт
func GetAllBankCards(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	_, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	cardData, err := storage.BCST.GetAllRecords(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	decData, ok := cardData.([]storage.BankCardData)
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

	result, err := json.Marshal(decData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(result)

}
