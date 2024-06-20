package requests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mock_storage "github.com/Azcarot/PasswordStorage/internal/mock"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/golang/mock/gomock"
)

func TestRegistrationReq(t *testing.T) {
	type args struct {
		data storage.RegisterRequest
		resp int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "NoErr",
			args: args{data: storage.RegisterRequest{Login: "Login", Password: "Password"}, resp: http.StatusOK},
			want: true, wantErr: false},
		{name: "Ok",
			args: args{data: storage.RegisterRequest{Login: "Login", Password: "Password"}, resp: http.StatusAccepted},
			want: false, wantErr: true},
		{name: "Witherr2",
			args: args{data: storage.RegisterRequest{Login: "Login", Password: "Password"}, resp: http.StatusUnprocessableEntity},
			want: false, wantErr: true},
		{name: "Witherr3",
			args: args{data: storage.RegisterRequest{Login: "", Password: ""}, resp: http.StatusUnprocessableEntity},
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
				mock.EXPECT().CreateNewUser(gomock.Any(), gomock.Any()).Times(1)
				mock.EXPECT().GetSecretKey(gomock.Any()).Times(1)
			}
			storage.LiteST = mock
			got, err := RegistrationReq(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegistrationReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RegistrationReq() = %v, want %v", got, tt.want)
			}
		})
	}
}
