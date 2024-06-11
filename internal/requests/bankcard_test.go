// Package requests - модуль с запросами к серверу

// Включает запросы на регистрацию/авторизацию пользователя

// Создание/обновление/Удаление/Синхронизация всех данных пользователя

package requests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_storage "github.com/Azcarot/PasswordStorage/internal/mock"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/golang/mock/gomock"
)

func TestAddCardReq(t *testing.T) {
	type args struct {
		data storage.BankCardData
		resp int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "withErr",
			args: args{data: storage.BankCardData{CardNumber: "111", Cvc: "222", ExpDate: "22/23", User: "User"}, resp: http.StatusOK},
			want: false, wantErr: true},
		{name: "Ok",
			args: args{data: storage.BankCardData{CardNumber: "111", Cvc: "222", ExpDate: "22/23", User: "User"}, resp: http.StatusAccepted},
			want: true, wantErr: false},
		{name: "Witherr2",
			args: args{data: storage.BankCardData{CardNumber: "111", Cvc: "222", ExpDate: "22/23", User: "User"}, resp: http.StatusUnprocessableEntity},
			want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.resp)

				w.Write([]byte(`{"message": "Hello, World!"}`))

			}))
			storage.ServURL = mockServer.URL
			got, err := AddCardReq(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCardReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddCardReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateCardReq(t *testing.T) {
	type args struct {
		data storage.BankCardData
		resp int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "withErr",
			args: args{data: storage.BankCardData{CardNumber: "111", Cvc: "222", ExpDate: "22/23", User: "User"}, resp: http.StatusOK},
			want: false, wantErr: true},
		{name: "Ok",
			args: args{data: storage.BankCardData{CardNumber: "111", Cvc: "222", ExpDate: "22/23", User: "User"}, resp: http.StatusAccepted},
			want: true, wantErr: false},
		{name: "Witherr2",
			args: args{data: storage.BankCardData{CardNumber: "111", Cvc: "222", ExpDate: "22/23", User: "User"}, resp: http.StatusUnprocessableEntity},
			want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.resp)
				w.Write([]byte(`{"message": "Hello, World!"}`))
			}))
			storage.ServURL = mockServer.URL
			got, err := UpdateCardReq(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCardReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateCardReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteCardReq(t *testing.T) {
	type args struct {
		data storage.BankCardData
		resp int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "withErr",
			args: args{data: storage.BankCardData{CardNumber: "111", Cvc: "222", ExpDate: "22/23", User: "User"}, resp: http.StatusOK},
			want: true, wantErr: false},
		{name: "Ok",
			args: args{data: storage.BankCardData{CardNumber: "111", Cvc: "222", ExpDate: "22/23", User: "User"}, resp: http.StatusAccepted},
			want: false, wantErr: true},
		{name: "Witherr2",
			args: args{data: storage.BankCardData{CardNumber: "111", Cvc: "222", ExpDate: "22/23", User: "User"}, resp: http.StatusUnprocessableEntity},
			want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.resp)
				w.Write([]byte(`{"message": "Hello, World!"}`))
			}))
			storage.ServURL = mockServer.URL
			got, err := DeleteCardReq(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteCardReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteCardReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncCardReq(t *testing.T) {
	tests := []struct {
		name     string
		want     bool
		wantErr  bool
		resp     int
		respData []storage.BankCardData
	}{
		{name: "withErr",
			resp: http.StatusInternalServerError,
			want: false, wantErr: true},
		{name: "NoData1",
			resp: http.StatusAccepted,
			want: false, wantErr: true},
		{name: "WithData",
			resp:     http.StatusAccepted,
			respData: []storage.BankCardData{{ID: 1, CardNumber: "111", Cvc: "222", ExpDate: "333"}},
			want:     false, wantErr: false},
		{name: "NoData2",
			resp: http.StatusOK,
			want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mock_storage.NewMockPgxStorage(ctrl)
			mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1)
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				w.WriteHeader(tt.resp)
				if len(tt.respData) != 0 {
					mock.EXPECT().AddData(tt.respData[0]).Times(1)
					mock.EXPECT().CreateNewRecord(gomock.Any()).Times(1)
					mock.EXPECT().GetAllRecords(gomock.Any()).Times(1)
					data, _ := json.Marshal(tt.respData)
					w.Write(data)
				} else {
					w.Write([]byte(`{"message": "Hello, World!"}`))
				}
			}))
			storage.ServURL = mockServer.URL

			storage.BCLiteS = mock
			got, err := SyncCardReq()
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncCardReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SyncCardReq() = %v, want %v", got, tt.want)
			}
		})
	}
}
