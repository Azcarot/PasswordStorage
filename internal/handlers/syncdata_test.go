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

func TestSyncBankData(t *testing.T) {
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
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
		{name: "NoErr", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
		{name: "GetErr", login: "login",
			resp: http.StatusInternalServerError,
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
				case "NoErr":
					mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1).Return("21234234", nil)
					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return([]storage.BankCardData{}, pgx.ErrNoRows)
				case "GetErr":
					mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1).Return("21234234", nil)
					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return([]storage.BankCardData{{}}, nil)
				case "Unprocess":
					mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1).Return("", fmt.Errorf("error"))
				default:

					mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1)
				}

				ctx = context.WithValue(req.Context(), storage.UserLoginCtxKey, tt.login)
			}

			req = req.WithContext(ctx)

			res := httptest.NewRecorder()

			SyncBankData(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}

func TestSyncFileData(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		want     bool
		wantErr  bool
		resp     int
		respData storage.FileData
	}{
		{name: "withErr",
			resp: http.StatusInternalServerError,
			want: false, wantErr: true},
		{name: "WithData", login: "login",
			resp: http.StatusOK,
			respData: storage.FileData{
				FileName: "11",
				Comment:  "Test Card",
			},
			want: true, wantErr: false},
		{name: "Unprocess", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
		{name: "NoErr", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
		{name: "GetErr", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mock_storage.NewMockPgxStorage(ctrl)
			storage.FST = mock

			body, err := json.Marshal(tt.respData)
			assert.NoError(t, err)
			cypher.CypherFileData(context.Background(), &tt.respData)

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/file/sync", bytes.NewBuffer(body))
			ctx := req.Context()
			if len(tt.login) != 0 {
				switch tt.name {
				case "NoErr":
					mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1).Return("21234234", nil)
					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return([]storage.FileData{}, pgx.ErrNoRows)
				case "GetErr":
					mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1).Return("21234234", nil)
					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return([]storage.FileData{{}}, nil)
				case "Unprocess":
					mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1).Return("", fmt.Errorf("error"))
				default:

					mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1)
				}

				ctx = context.WithValue(req.Context(), storage.UserLoginCtxKey, tt.login)
			}

			req = req.WithContext(ctx)

			res := httptest.NewRecorder()

			SyncFileData(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}

func TestSyncTextData(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		want     bool
		wantErr  bool
		resp     int
		respData storage.TextData
	}{
		{name: "withErr",
			resp: http.StatusInternalServerError,
			want: false, wantErr: true},
		{name: "WithData", login: "login",
			resp: http.StatusInternalServerError,
			respData: storage.TextData{
				Text:    "11",
				Comment: "Test Card",
			},
			want: true, wantErr: false},
		{name: "Unprocess", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
		{name: "NoErr", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
		{name: "GetErr", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mock_storage.NewMockPgxStorage(ctrl)
			storage.TST = mock

			body, err := json.Marshal(tt.respData)
			assert.NoError(t, err)
			cypher.CypherTextData(context.Background(), &tt.respData)

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/text/sync", bytes.NewBuffer(body))
			ctx := req.Context()
			if len(tt.login) != 0 {
				switch tt.name {
				case "NoErr":
					mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1).Return("21234234", nil)
					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return([]storage.TextData{}, pgx.ErrNoRows)
				case "GetErr":
					mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1).Return("21234234", nil)
					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return([]storage.TextData{{}}, nil)
				case "Unprocess":
					mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1).Return("", fmt.Errorf("error"))
				default:

					mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1)
					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return(gomock.Any(), nil)
				}

				ctx = context.WithValue(req.Context(), storage.UserLoginCtxKey, tt.login)
			}

			req = req.WithContext(ctx)

			res := httptest.NewRecorder()

			SyncTextData(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}

func TestSyncLPWData(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		want     bool
		wantErr  bool
		resp     int
		respData storage.LoginData
	}{
		{name: "withErr",
			resp: http.StatusInternalServerError,
			want: false, wantErr: true},
		{name: "WithData", login: "login",
			resp: http.StatusOK,
			respData: storage.LoginData{
				Login:   "11",
				Comment: "Test Card",
			},
			want: true, wantErr: false},
		{name: "Unprocess", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
		{name: "NoErr", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
		{name: "GetErr", login: "login",
			resp: http.StatusOK,
			want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mock_storage.NewMockPgxStorage(ctrl)
			storage.LPST = mock

			body, err := json.Marshal(tt.respData)
			assert.NoError(t, err)
			cypher.CypherLPWData(context.Background(), &tt.respData)

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/lpw/sync", bytes.NewBuffer(body))
			ctx := req.Context()
			if len(tt.login) != 0 {
				switch tt.name {
				case "NoErr":
					mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1).Return("21234234", nil)
					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return([]storage.LoginData{}, pgx.ErrNoRows)
				case "CyphErr":
					mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1).Return("21234234", nil)
					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return([]storage.LoginData{{}}, nil)
				case "Unprocess":
					mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1).Return("", fmt.Errorf("error"))
				default:

					mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1)
				}

				ctx = context.WithValue(req.Context(), storage.UserLoginCtxKey, tt.login)
			}

			req = req.WithContext(ctx)

			res := httptest.NewRecorder()

			SyncLPWData(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}
