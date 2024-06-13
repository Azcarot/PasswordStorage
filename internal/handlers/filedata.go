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

// AddNewFile - ручка для добавления новых данных типа "файл"
func AddNewFile(res http.ResponseWriter, req *http.Request) {
	var userData storage.UserData
	ctx := req.Context()
	dataLogin, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var fileData storage.FileData
	userData.Login = dataLogin
	data, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &fileData)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileData.User = userData.Login
	mut := sync.Mutex{}
	mut.Lock()
	defer mut.Unlock()
	fileData.Date = time.Now().Format(time.RFC3339)
	err = cypher.CypherFileData(ctx, &fileData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	err = storage.FST.AddData(fileData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = storage.FST.CreateNewRecord(ctx)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	res.WriteHeader(http.StatusAccepted)
}

// GetFile - ручка для получения фаловых данных по id
func GetFile(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	_, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var reqData storage.FileData

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
	err = storage.FST.AddData(reqData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	fileData, err := storage.FST.GetRecord(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	cyphFileData, ok := fileData.(storage.FileData)
	if !ok {

		res.WriteHeader(http.StatusInternalServerError)
		return

	}
	err = cypher.DeCypherFileData(ctx, &cyphFileData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(cyphFileData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(result)

}

// UpdateFile - ручка для обновления файловых данных по id
func UpdateFile(res http.ResponseWriter, req *http.Request) {
	var userData storage.UserData
	ctx := req.Context()
	dataLogin, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var fileData storage.FileData
	userData.Login = dataLogin
	data, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &fileData)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileData.User = userData.Login
	mut := sync.Mutex{}
	mut.Lock()
	defer mut.Unlock()
	fileData.Date = time.Now().Format(time.RFC3339)
	err = storage.FST.AddData(fileData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	old, err := storage.FST.GetRecord(ctx)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	oldData, ok := old.(storage.FileData)

	if !ok {

		res.WriteHeader(http.StatusInternalServerError)
		return

	}
	err = cypher.DeCypherFileData(ctx, &oldData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ok {
		if fileData.FileName == "" {
			fileData.FileName = oldData.FileName
		}
		if fileData.Data == "" {
			fileData.Data = oldData.Data
		}
		err := cypher.CypherFileData(ctx, &fileData)
		if err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		err = storage.FST.AddData(fileData)
		if err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}
	err = storage.FST.UpdateRecord(ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusAccepted)
}

// DeleteFile - ручека для удаления файла по id
func DeleteFile(res http.ResponseWriter, req *http.Request) {
	var userData storage.UserData
	ctx := req.Context()
	dataLogin, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var fileData storage.FileData
	userData.Login = dataLogin
	data, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &fileData)
	if err != nil {

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileData.User = userData.Login
	mut := sync.Mutex{}
	mut.Lock()
	defer mut.Unlock()
	fileData.Date = time.Now().Format(time.RFC3339)
	err = storage.FST.AddData(fileData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	err = storage.FST.DeleteRecord(ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

// SearchFile - ручка для поиска фйла по строке
func SearchFile(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	_, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var reqData storage.FileData

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
	err = storage.FST.AddData(reqData)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	fileData, err := storage.FST.SearchRecord(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	cyphData, ok := fileData.(storage.FileData)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = cypher.DeCypherFileData(ctx, &cyphData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
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

// GetAllFiles - ручка для получения всех сохраненных файлов
func GetAllFiles(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	_, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileData, err := storage.FST.GetAllRecords(ctx)

	if errors.Is(err, pgx.ErrNoRows) {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	files, ok := fileData.([]storage.FileData)
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

	result, err := json.Marshal(fileData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(result)

}
