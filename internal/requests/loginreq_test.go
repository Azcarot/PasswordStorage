package requests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azcarot/PasswordStorage/internal/handlers"
	mock_storage "github.com/Azcarot/PasswordStorage/internal/mock"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/golang/mock/gomock"
)

func TestLoginReq(t *testing.T) {
	type args struct {
		data handlers.LoginRequest
		resp int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "NoErr",
			args: args{data: handlers.LoginRequest{Login: "Login", Password: "Password"}, resp: http.StatusOK},
			want: true, wantErr: false},
		{name: "Ok",
			args: args{data: handlers.LoginRequest{Login: "Login", Password: "Password"}, resp: http.StatusAccepted},
			want: false, wantErr: true},
		{name: "Witherr2",
			args: args{data: handlers.LoginRequest{Login: "Login", Password: "Password"}, resp: http.StatusUnprocessableEntity},
			want: false, wantErr: true},
		{name: "Witherr3",
			args: args{data: handlers.LoginRequest{Login: "", Password: ""}, resp: http.StatusUnprocessableEntity},
			want: false, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.resp)
				w.Write([]byte(`{"message": "Hello, World!"}`))
			}))
			storage.ServURL = mockServer.URL
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mock_storage.NewMockLiteConn(ctrl)
			if !tt.wantErr {

				mock.EXPECT().GetSecretKey(gomock.Any()).Times(1)
			}
			storage.LiteST = mock
			got, err := LoginReq(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LoginReq() = %v, want %v", got, tt.want)
			}
		})
	}
}
