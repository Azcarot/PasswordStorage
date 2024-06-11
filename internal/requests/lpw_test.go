package requests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mock_storage "github.com/Azcarot/PasswordStorage/internal/mock"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/golang/mock/gomock"
)

func TestAddLPWReq(t *testing.T) {
	type args struct {
		data storage.LoginData
		resp int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "withErr",
			args: args{data: storage.LoginData{Login: "Login", Password: "PWD", User: "User"}, resp: http.StatusOK},
			want: false, wantErr: true},
		{name: "Ok",
			args: args{data: storage.LoginData{Login: "Login", Password: "PWD", User: "User"}, resp: http.StatusAccepted},
			want: true, wantErr: false},
		{name: "Witherr2",
			args: args{data: storage.LoginData{Login: "Login", Password: "PWD", User: "User"}, resp: http.StatusUnprocessableEntity},
			want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.resp)
				w.Write([]byte(`{"message": "Hello, World!"}`))
			}))
			storage.ServURL = mockServer.URL
			got, err := AddLPWReq(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddLPWReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddLPWReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateLPWReq(t *testing.T) {
	type args struct {
		data storage.LoginData
		resp int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "withErr",
			args: args{data: storage.LoginData{Login: "Login", Password: "PWD", User: "User"}, resp: http.StatusOK},
			want: false, wantErr: true},
		{name: "Ok",
			args: args{data: storage.LoginData{Login: "Login", Password: "PWD", User: "User"}, resp: http.StatusAccepted},
			want: true, wantErr: false},
		{name: "Witherr2",
			args: args{data: storage.LoginData{Login: "Login", Password: "PWD", User: "User"}, resp: http.StatusUnprocessableEntity},
			want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.resp)
				w.Write([]byte(`{"message": "Hello, World!"}`))
			}))
			storage.ServURL = mockServer.URL
			got, err := UpdateLPWReq(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateLPWReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateLPWReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteLPWReq(t *testing.T) {
	type args struct {
		data storage.LoginData
		resp int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "withErr",
			args: args{data: storage.LoginData{Login: "Login", Password: "PWD", User: "User"}, resp: http.StatusOK},
			want: true, wantErr: false},
		{name: "Ok",
			args: args{data: storage.LoginData{Login: "Login", Password: "PWD", User: "User"}, resp: http.StatusAccepted},
			want: false, wantErr: true},
		{name: "Witherr2",
			args: args{data: storage.LoginData{Login: "Login", Password: "PWD", User: "User"}, resp: http.StatusUnprocessableEntity},
			want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.resp)
				w.Write([]byte(`{"message": "Hello, World!"}`))
			}))
			storage.ServURL = mockServer.URL
			got, err := DeleteLPWReq(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteLPWReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteLPWReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncLPWReq(t *testing.T) {
	tests := []struct {
		name    string
		want    bool
		wantErr bool
		resp    int
	}{
		{name: "withErr",
			resp: http.StatusInternalServerError,
			want: false, wantErr: true},
		{name: "NoData1",
			resp: http.StatusAccepted,
			want: false, wantErr: true},
		{name: "NoData2",
			resp: http.StatusOK,
			want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.resp)
				w.Write([]byte(`{"message": "Hello, World!"}`))
			}))
			storage.ServURL = mockServer.URL
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mock_storage.NewMockPgxStorage(ctrl)
			mock.EXPECT().HashDatabaseData(gomock.Any()).Times(1)
			storage.LPWLiteS = mock
			got, err := SyncLPWReq()
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncLPWReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SyncLPWReq() = %v, want %v", got, tt.want)
			}
		})
	}
}
