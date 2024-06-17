package cypher

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/Azcarot/PasswordStorage/internal/storage"
)

type errorReader struct{}

// Read method for errorReader always returns an error
func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("random source error")
}
func TestCypherBankData(t *testing.T) {
	type args struct {
		ctx     context.Context
		data    *storage.BankCardData
		randErr bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{ctx: context.Background(), randErr: false, data: &storage.BankCardData{ID: 1, CardNumber: "", FullName: "22222"}}, wantErr: false},
		{name: "test1", args: args{ctx: context.Background(), randErr: true, data: &storage.BankCardData{ID: 1, CardNumber: "", FullName: "22222"}}, wantErr: true},
	}
	for _, tt := range tests {
		if tt.args.randErr == true {
			oldRandReader := rand.Reader
			rand.Reader = &errorReader{}
			defer func() { rand.Reader = oldRandReader }()
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := CypherBankData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("CypherBankData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeCypherBankData(t *testing.T) {
	type args struct {
		ctx    context.Context
		data   *storage.BankCardData
		secret bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{ctx: context.Background(), data: &storage.BankCardData{ID: 1, CardNumber: "1111", FullName: "22222"}}, wantErr: true},
		{name: "test1", args: args{ctx: context.Background(), secret: true, data: &storage.BankCardData{ID: 1, CardNumber: "1111 11111 1111 1111", FullName: "22222", Cvc: "1111", ExpDate: "1111"}}, wantErr: false},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if tt.args.secret {
				CypherBankData(tt.args.ctx, tt.args.data)
			}
			if err := DeCypherBankData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("DeCypherBankData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
