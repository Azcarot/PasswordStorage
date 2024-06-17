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

func TestAddNewFile(t *testing.T) {
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
		{name: "NoData1", login: "login",
			resp: http.StatusAccepted,
			want: false, wantErr: true},
		{name: "WithData", login: "login",
			resp: http.StatusAccepted,
			respData: storage.FileData{
				FileName: "Name",
				Path:     "Path",
				Data:     "Data",
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
			storage.FST = mock

			body, err := json.Marshal(tt.respData)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/file/add", bytes.NewBuffer(body))
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

			AddNewFile(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}

func TestGetFile(t *testing.T) {
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
				FileName: "Name",
				Path:     "Path",
				Data:     "Data",
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
		{name: "UnmErr", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mock_storage.NewMockPgxStorage(ctrl)
			storage.FST = mock

			var body []byte
			var err error
			if tt.name == "UnmErr" {
				body, err = json.Marshal("Something")
			} else {
				body, err = json.Marshal(tt.respData)
			}
			assert.NoError(t, err)
			cypher.CypherFileData(context.Background(), &tt.respData)

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/file/get", bytes.NewBuffer(body))
			ctx := req.Context()
			if len(tt.login) != 0 {
				switch tt.name {
				case "UnmErr":
				case "GetErr":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return(storage.FileData{}, pgx.ErrNoRows)
				case "PGerr":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return(storage.FileData{}, fmt.Errorf("err"))

				case "Format":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return([]storage.FileData{}, nil)
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

			GetFile(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}

func TestUpdateFile(t *testing.T) {
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
			resp: http.StatusAccepted,
			respData: storage.FileData{
				FileName: "Name",
				Path:     "Path",
				Data:     "Data",
			},
			want: true, wantErr: false},
		{name: "WithData2", login: "login",
			resp: http.StatusUnprocessableEntity,
			respData: storage.FileData{
				FileName: "",
				Path:     "",
				Data:     "",
				Comment:  "",
			},
			want: true, wantErr: false},
		{name: "WithData3", login: "login",
			resp: http.StatusInternalServerError,
			respData: storage.FileData{
				FileName: "",
				Path:     "",
				Data:     "",
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

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/file/get", bytes.NewBuffer(body))
			ctx := req.Context()
			if len(tt.login) != 0 {
				switch tt.name {
				case "GetErr":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return(storage.FileData{}, pgx.ErrNoRows)
				case "Unprocess":
					mock.EXPECT().AddData(gomock.Any()).Times(1).Return(fmt.Errorf("error"))
				case "WithData2":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return(tt.respData, nil)
					mock.EXPECT().AddData(gomock.Any()).Times(1).Return(fmt.Errorf("error"))
				case "Format":
					mock.EXPECT().AddData(gomock.Any()).Times(1)
					mock.EXPECT().GetRecord(gomock.Any()).Times(1).Return([]storage.FileData{}, nil)
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

			UpdateFile(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}

func TestDeleteFile(t *testing.T) {
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
				FileName: "Name",
				Path:     "Path",
				Data:     "Data",
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
			storage.FST = mock

			body, err := json.Marshal(tt.respData)
			assert.NoError(t, err)
			cypher.CypherFileData(context.Background(), &tt.respData)

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/file/get", bytes.NewBuffer(body))
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

			DeleteFile(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}

func TestGetAllFiles(t *testing.T) {
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
				FileName: "Name",
				Path:     "Path",
				Data:     "Data",
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
			respData: storage.FileData{
				FileName: "Name",
				Path:     "Path",
				Data:     "Data",
			},
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

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/file/get", bytes.NewBuffer(body))
			ctx := req.Context()
			if len(tt.login) != 0 {
				switch tt.name {
				case "GetErr":

					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return(storage.FileData{}, pgx.ErrNoRows)
				case "GetErr2":

					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return(storage.FileData{}, fmt.Errorf("error"))
				case "WithDataErr":

					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return(storage.FileData{}, nil)
				default:
					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return([]storage.FileData{}, nil)

				}

				ctx = context.WithValue(req.Context(), storage.UserLoginCtxKey, tt.login)
			}

			req = req.WithContext(ctx)

			res := httptest.NewRecorder()

			GetAllFiles(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}
