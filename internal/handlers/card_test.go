// Package handlers - все обработчики запросов на сервере

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azcarot/PasswordStorage/internal/cypher"
	mock_storage "github.com/Azcarot/PasswordStorage/internal/mock"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestAddNewCard(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		want     bool
		wantErr  bool
		resp     int
		respData storage.BankCardData
	}{
		{name: "withErr",
			resp: http.StatusInternalServerError,
			want: false, wantErr: true},
		{name: "NoData1", login: "login",
			resp: http.StatusAccepted,
			want: false, wantErr: true},
		{name: "WithData", login: "login",
			resp: http.StatusAccepted,
			respData: storage.BankCardData{
				CardNumber: "1234567890123456",
				Cvc:        "123",
				ExpDate:    "12/24",
				FullName:   "John Doe",
				Comment:    "Test Card",
			},
			want: true, wantErr: false},
		{name: "Unprocess", login: "login",
			resp: http.StatusUnprocessableEntity,
			want: true, wantErr: false},
		{name: "CreateErr", login: "login",
			resp: http.StatusUnprocessableEntity,
			want: true, wantErr: false},
		{name: "UnmErr", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mock_storage.NewMockPgxStorage(ctrl)
			storage.BCST = mock

			var body []byte
			var err error
			if tt.name == "UnmErr" {
				body, err = json.Marshal("Something")
			} else {
				body, err = json.Marshal(tt.respData)
			}
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/card/add", bytes.NewBuffer(body))
			ctx := req.Context()
			if len(tt.login) != 0 {
				switch tt.name {
				case "UnmErr":
				case "CreateErr":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().CreateNewRecord(gomock.Any()).Times(1).Return(fmt.Errorf("error"))
				case "Unprocess":
					mock.EXPECT().AddData(gomock.Any()).Times(1).Return(fmt.Errorf("error"))
				default:
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().CreateNewRecord(gomock.Any()).Times(1)
				}
				ctx = context.WithValue(req.Context(), storage.UserLoginCtxKey, tt.login)
			}

			req = req.WithContext(ctx)

			res := httptest.NewRecorder()

			AddNewCard(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}

func TestGetBankCard(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		want     bool
		wantErr  bool
		resp     int
		respData storage.BankCardData
	}{
		{name: "withErr",
			resp: http.StatusInternalServerError,
			want: false, wantErr: true},
		{name: "WithData", login: "login",
			resp: http.StatusOK,
			respData: storage.BankCardData{
				CardNumber: "1234567890123456",
				Cvc:        "123",
				ExpDate:    "12/24",
				FullName:   "John Doe",
				Comment:    "Test Card",
			},
			want: true, wantErr: false},
		{name: "Unprocess", login: "login",
			resp: http.StatusUnprocessableEntity,
			want: true, wantErr: false},
		{name: "GetErr", login: "login",
			resp: http.StatusNoContent,
			want: true, wantErr: false},
		{name: "PGerr", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
		{name: "UnmErr", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mock_storage.NewMockPgxStorage(ctrl)
			storage.BCST = mock
			var body []byte
			var err error
			if tt.name == "UnmErr" {
				body, err = json.Marshal("Something")
			} else {
				body, err = json.Marshal(tt.respData)
			}
			assert.NoError(t, err)
			cypher.CypherBankData(context.Background(), &tt.respData)

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/card/get", bytes.NewBuffer(body))
			ctx := req.Context()
			if len(tt.login) != 0 {
				switch tt.name {
				case "UnmErr":

				case "GetErr":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return(storage.BankCardData{}, pgx.ErrNoRows)
				case "PGerr":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return(storage.BankCardData{}, fmt.Errorf("err"))

				case "Unprocess":
					mock.EXPECT().AddData(gomock.Any()).Times(1).Return(fmt.Errorf("error"))
				default:

					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return(tt.respData, nil)
				}

				ctx = context.WithValue(req.Context(), storage.UserLoginCtxKey, tt.login)
			}

			req = req.WithContext(ctx)

			res := httptest.NewRecorder()

			GetBankCard(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}

func TestUpdateBankCard(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		want     bool
		wantErr  bool
		resp     int
		respData storage.BankCardData
	}{
		{name: "withErr",
			resp: http.StatusInternalServerError,
			want: false, wantErr: true},
		{name: "WithData", login: "login",
			resp: http.StatusAccepted,
			respData: storage.BankCardData{
				CardNumber: "1234567890123456",
				Cvc:        "123",
				ExpDate:    "12/24",
				FullName:   "John Doe",
				Comment:    "Test Card",
			},
			want: true, wantErr: false},
		{name: "WithData2", login: "login",
			resp: http.StatusUnprocessableEntity,
			respData: storage.BankCardData{
				CardNumber: "",
				Cvc:        "",
				ExpDate:    "",
				FullName:   "",
				Comment:    "",
			},
			want: true, wantErr: false},
		{name: "WithData3", login: "login",
			resp: http.StatusInternalServerError,
			respData: storage.BankCardData{
				CardNumber: "",
				Cvc:        "",
				ExpDate:    "",
				FullName:   "",
				Comment:    "",
			},
			want: true, wantErr: false},
		{name: "Unprocess", login: "login",
			resp: http.StatusUnprocessableEntity,
			want: true, wantErr: false},
		{name: "GetErr", login: "login",
			resp: http.StatusUnprocessableEntity,
			want: true, wantErr: false},
		{name: "UnmErr", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mock_storage.NewMockPgxStorage(ctrl)
			storage.BCST = mock
			var body []byte
			var err error
			if tt.name == "UnmErr" {
				body, err = json.Marshal("Something")
			} else {
				body, err = json.Marshal(tt.respData)
			}
			assert.NoError(t, err)
			cypher.CypherBankData(context.Background(), &tt.respData)

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/card/get", bytes.NewBuffer(body))
			ctx := req.Context()
			if len(tt.login) != 0 {
				switch tt.name {
				case "UnmErr":

				case "GetErr":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return(storage.BankCardData{}, pgx.ErrNoRows)
				case "Unprocess":
					mock.EXPECT().AddData(gomock.Any()).Times(1).Return(fmt.Errorf("error"))
				case "WithData2":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return(tt.respData, nil)
					mock.EXPECT().AddData(gomock.Any()).Times(1).Return(fmt.Errorf("error"))
				case "WithData3":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return(tt.respData, nil)
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().UpdateRecord(gomock.Any()).Times(1).Return(fmt.Errorf("error"))
				default:

					mock.EXPECT().AddData(gomock.Any()).Times(2)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return(tt.respData, nil)
					mock.EXPECT().UpdateRecord(gomock.Any()).Times(1).Return(nil)
				}

				ctx = context.WithValue(req.Context(), storage.UserLoginCtxKey, tt.login)
			}

			req = req.WithContext(ctx)

			res := httptest.NewRecorder()

			UpdateCard(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}

func TestDeleteBankCard(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		want     bool
		wantErr  bool
		resp     int
		respData storage.BankCardData
	}{
		{name: "withErr",
			resp: http.StatusInternalServerError,
			want: false, wantErr: true},
		{name: "WithData", login: "login",
			resp: http.StatusOK,
			respData: storage.BankCardData{
				CardNumber: "1234567890123456",
				Cvc:        "123",
				ExpDate:    "12/24",
				FullName:   "John Doe",
				Comment:    "Test Card",
			},
			want: true, wantErr: false},
		{name: "Unprocess", login: "login",
			resp: http.StatusUnprocessableEntity,
			want: true, wantErr: false},
		{name: "GetErr", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
		{name: "UnmErr", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mock_storage.NewMockPgxStorage(ctrl)
			storage.BCST = mock

			var body []byte
			var err error
			if tt.name == "UnmErr" {
				body, err = json.Marshal("Something")
			} else {
				body, err = json.Marshal(tt.respData)
			}
			assert.NoError(t, err)
			cypher.CypherBankData(context.Background(), &tt.respData)

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/card/get", bytes.NewBuffer(body))
			ctx := req.Context()
			if len(tt.login) != 0 {
				switch tt.name {
				case "UnmErr":

				case "GetErr":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().DeleteRecord(gomock.Any()).Times(1).Return(fmt.Errorf("error"))
				case "Unprocess":
					mock.EXPECT().AddData(gomock.Any()).Times(1).Return(fmt.Errorf("error"))
				default:

					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().DeleteRecord(gomock.Any()).Times(1).Return(nil)
				}

				ctx = context.WithValue(req.Context(), storage.UserLoginCtxKey, tt.login)
			}

			req = req.WithContext(ctx)

			res := httptest.NewRecorder()

			DeleteCard(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}

func TestGetAllBankCards(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		want     bool
		wantErr  bool
		resp     int
		respData storage.BankCardData
	}{
		{name: "withErr",
			resp: http.StatusInternalServerError,
			want: false, wantErr: true},
		{name: "WithData", login: "login",
			resp: http.StatusOK,
			respData: storage.BankCardData{
				CardNumber: "1234567890123456",
				Cvc:        "123",
				ExpDate:    "12/24",
				FullName:   "John Doe",
				Comment:    "Test Card",
			},
			want: true, wantErr: false},

		{name: "GetErr", login: "login",
			resp: http.StatusNoContent,
			want: true, wantErr: false},
		{name: "GetErr2", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},

		{name: "WithDataErr", login: "login",
			resp: http.StatusInternalServerError,
			respData: storage.BankCardData{
				CardNumber: "1234567890123456",
				Cvc:        "123",
				ExpDate:    "12/24",
				FullName:   "John Doe",
				Comment:    "Test Card",
			},
			want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mock_storage.NewMockPgxStorage(ctrl)
			storage.BCST = mock

			body, err := json.Marshal(tt.respData)
			assert.NoError(t, err)
			cypher.CypherBankData(context.Background(), &tt.respData)

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/card/get", bytes.NewBuffer(body))
			ctx := req.Context()
			if len(tt.login) != 0 {
				switch tt.name {
				case "GetErr":

					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return(storage.BankCardData{}, pgx.ErrNoRows)
				case "GetErr2":

					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return(storage.BankCardData{}, fmt.Errorf("error"))
				case "WithDataErr":

					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return(storage.BankCardData{}, nil)
				default:
					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return([]storage.BankCardData{}, nil)

				}

				ctx = context.WithValue(req.Context(), storage.UserLoginCtxKey, tt.login)
			}

			req = req.WithContext(ctx)

			res := httptest.NewRecorder()

			GetAllBankCards(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}
