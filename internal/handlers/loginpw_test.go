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

func TestAddNewLPW(t *testing.T) {
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
		{name: "NoData1", login: "login",
			resp: http.StatusAccepted,
			want: false, wantErr: true},
		{name: "WithData", login: "login",
			resp: http.StatusAccepted,
			respData: storage.LoginData{
				Login:    "Name",
				Password: "Password",
				Comment:  "Data",
			},
			want: true, wantErr: false},
		{name: "Unprocess", login: "login",
			resp: http.StatusUnprocessableEntity,
			want: true, wantErr: false},
		{name: "CreateErr", login: "login",
			resp: http.StatusUnprocessableEntity,
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

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/lpw/add", bytes.NewBuffer(body))
			ctx := req.Context()
			if len(tt.login) != 0 {
				switch tt.name {
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

			AddNewLoginPw(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}

func TestGetLPW(t *testing.T) {
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
				Login:    "Name",
				Password: "Password",
				Comment:  "Data",
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
		{name: "Format", login: "login",
			resp: http.StatusInternalServerError,
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

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/lpw/get", bytes.NewBuffer(body))
			ctx := req.Context()
			if len(tt.login) != 0 {
				switch tt.name {
				case "GetErr":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return(storage.LoginData{}, pgx.ErrNoRows)
				case "PGerr":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return(storage.LoginData{}, fmt.Errorf("err"))

				case "Format":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return([]storage.LoginData{}, nil)
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

			GetLoginPW(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}

func TestUpdateLPW(t *testing.T) {
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
			resp: http.StatusAccepted,
			respData: storage.LoginData{
				Login:    "Name",
				Password: "Password",
				Comment:  "Data",
			},
			want: true, wantErr: false},
		{name: "WithData2", login: "login",
			resp: http.StatusUnprocessableEntity,
			respData: storage.LoginData{
				Login:    "",
				Password: "",
				Comment:  "",
			},
			want: true, wantErr: false},
		{name: "WithData3", login: "login",
			resp: http.StatusInternalServerError,
			respData: storage.LoginData{
				Login:    "",
				Password: "",
				Comment:  "",
			},
			want: true, wantErr: false},
		{name: "Unprocess", login: "login",
			resp: http.StatusUnprocessableEntity,
			want: true, wantErr: false},
		{name: "GetErr", login: "login",
			resp: http.StatusUnprocessableEntity,
			want: true, wantErr: false},
		{name: "Format", login: "login",
			resp: http.StatusAccepted,
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

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/lpw/get", bytes.NewBuffer(body))
			ctx := req.Context()
			if len(tt.login) != 0 {
				switch tt.name {
				case "GetErr":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return(storage.LoginData{}, pgx.ErrNoRows)
				case "Unprocess":
					mock.EXPECT().AddData(gomock.Any()).Times(1).Return(fmt.Errorf("error"))
				case "WithData2":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return(tt.respData, nil)
					mock.EXPECT().AddData(gomock.Any()).Times(1).Return(fmt.Errorf("error"))
				case "Format":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return([]storage.LoginData{}, nil)
					mock.EXPECT().UpdateRecord(gomock.Any()).Times(1).Return(nil)
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

			UpdateLoginPW(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}

func TestDeleteLPW(t *testing.T) {
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
				Login:    "Name",
				Password: "Password",
				Comment:  "Data",
			},
			want: true, wantErr: false},
		{name: "Unprocess", login: "login",
			resp: http.StatusUnprocessableEntity,
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
			storage.LPST = mock

			body, err := json.Marshal(tt.respData)
			assert.NoError(t, err)
			cypher.CypherLPWData(context.Background(), &tt.respData)

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/lpw/get", bytes.NewBuffer(body))
			ctx := req.Context()
			if len(tt.login) != 0 {
				switch tt.name {
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

			DeleteLoginPW(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}

func TestGetAllLPW(t *testing.T) {
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
				Login:    "Name",
				Password: "Password",
				Comment:  "Data",
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
			respData: storage.LoginData{
				Login:    "Name",
				Password: "Password",
				Comment:  "Data",
			},
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

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/lpw/get", bytes.NewBuffer(body))
			ctx := req.Context()
			if len(tt.login) != 0 {
				switch tt.name {
				case "GetErr":

					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return(storage.LoginData{}, pgx.ErrNoRows)
				case "GetErr2":

					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return(storage.LoginData{}, fmt.Errorf("error"))
				case "WithDataErr":

					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return(storage.LoginData{}, nil)
				default:
					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return([]storage.LoginData{}, nil)

				}

				ctx = context.WithValue(req.Context(), storage.UserLoginCtxKey, tt.login)
			}

			req = req.WithContext(ctx)

			res := httptest.NewRecorder()

			GetAllLoginPWs(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}
