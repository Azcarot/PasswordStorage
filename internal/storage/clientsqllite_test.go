package storage

import (
	"context"
	"testing"

	_ "modernc.org/sqlite"
)

func TestGenerateSecretKey(t *testing.T) {
	type args struct {
		data RegisterRequest
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "name", args: args{RegisterRequest{Login: "Login", Password: "Password"}}, want: "947a01e4ba47c8c4a84facb1a94273ae80cf636782deec8eb1251c91b7b87d8c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateSecretKey(tt.args.data); got != tt.want {
				t.Errorf("GenerateSecretKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLLiteStore_GetSecretKey(t *testing.T) {
	type args struct {
		login string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "1", args: args{login: "1111"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LiteST.GetSecretKey(tt.args.login); (err != nil) != tt.wantErr {
				t.Errorf("SQLLiteStore.GetSecretKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLLiteStore_CreateNewUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		data RegisterRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "1", args: args{ctx: context.Background(), data: RegisterRequest{Login: "111", Password: "2222"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LiteST.CreateNewUser(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("SQLLiteStore.CreateNewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
