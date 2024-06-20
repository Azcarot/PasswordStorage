package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_storage "github.com/Azcarot/PasswordStorage/internal/mock"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegistration(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		want     bool
		wantErr  bool
		resp     int
		respData storage.RegisterRequest
	}{
		{name: "WithErr",
			resp: http.StatusInternalServerError,
			want: false, wantErr: true},
		{name: "NoData1", login: "login",
			resp: http.StatusOK,
			want: false, wantErr: true},
		{name: "WithData", login: "login",
			resp: http.StatusOK,
			respData: storage.RegisterRequest{
				Login:    "1234567890123456",
				Password: "122131234",
			},
			want: true, wantErr: false},

		{name: "CreateErr", login: "login",
			resp: http.StatusConflict,
			want: true, wantErr: false},

		{name: "CreateErr2", login: "login",
			resp: http.StatusInternalServerError,
			want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mock_storage.NewMockPgxConn(ctrl)
			storage.ST = mock

			body, err := json.Marshal(tt.respData)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/card/add", bytes.NewBuffer(body))

			switch tt.name {
			case "CreateErr":
				mock.EXPECT().CheckUserExists(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
			case "WithErr":
				mock.EXPECT().CheckUserExists(gomock.Any(), gomock.Any()).Times(1).Return(false, fmt.Errorf("error"))
			case "CreateErr2":
				mock.EXPECT().CheckUserExists(gomock.Any(), gomock.Any()).Times(1)
				mock.EXPECT().CreateNewUser(gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("error"))

			default:
				mock.EXPECT().CheckUserExists(gomock.Any(), gomock.Any()).Times(1)
				mock.EXPECT().CreateNewUser(gomock.Any(), gomock.Any()).Times(1)
			}
			ctx := context.WithValue(req.Context(), storage.UserLoginCtxKey, tt.login)

			req = req.WithContext(ctx)

			res := httptest.NewRecorder()

			Registration(res, req)

			assert.Equal(t, tt.resp, res.Code)

		})
	}
}
