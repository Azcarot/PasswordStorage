package cypher

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/Azcarot/PasswordStorage/internal/storage"
)

func TestCypherLPWData(t *testing.T) {
	type args struct {
		ctx     context.Context
		data    *storage.LoginData
		randErr bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{ctx: context.Background(), randErr: false, data: &storage.LoginData{ID: 1, Login: "", Password: "22222"}}, wantErr: false},
		{name: "test1", args: args{ctx: context.Background(), randErr: true, data: &storage.LoginData{ID: 1, Login: "", Password: "22222"}}, wantErr: true},
	}
	for _, tt := range tests {
		if tt.args.randErr == true {
			oldRandReader := rand.Reader
			rand.Reader = &errorReader{}
			defer func() { rand.Reader = oldRandReader }()
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := CypherLPWData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("CypherLPWData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeCypherLPWData(t *testing.T) {
	type args struct {
		ctx    context.Context
		data   *storage.LoginData
		secret bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{ctx: context.Background(), data: &storage.LoginData{ID: 1, Login: "1111", Password: "22222"}}, wantErr: true},
		{name: "test1", args: args{ctx: context.Background(), data: &storage.LoginData{ID: 1, Password: "22222"}}, wantErr: true},
		{name: "test1", args: args{ctx: context.Background(), secret: true, data: &storage.LoginData{ID: 1, Login: "1111 11111 1111 1111", Password: "22222"}}, wantErr: false},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if tt.args.secret {
				CypherLPWData(tt.args.ctx, tt.args.data)
			}
			if err := DeCypherLPWData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("DeCypherLPWData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
