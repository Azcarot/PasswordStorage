package cypher

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/Azcarot/PasswordStorage/internal/storage"
)

func TestCypherFileData(t *testing.T) {
	type args struct {
		ctx     context.Context
		data    *storage.FileData
		randErr bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{ctx: context.Background(), randErr: false, data: &storage.FileData{ID: 1, FileName: "", Path: "22222"}}, wantErr: false},
		{name: "test1", args: args{ctx: context.Background(), randErr: true, data: &storage.FileData{ID: 1, FileName: "", Path: "22222"}}, wantErr: true},
	}
	for _, tt := range tests {
		if tt.args.randErr == true {
			oldRandReader := rand.Reader
			rand.Reader = &errorReader{}
			defer func() { rand.Reader = oldRandReader }()
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := CypherFileData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("CypherFileData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeCypherFileData(t *testing.T) {
	type args struct {
		ctx    context.Context
		data   *storage.FileData
		secret bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{ctx: context.Background(), data: &storage.FileData{ID: 1, FileName: "1111", Path: "22222"}}, wantErr: true},
		{name: "test1", args: args{ctx: context.Background(), secret: true, data: &storage.FileData{ID: 1, FileName: "1111 11111 1111 1111", Path: "22222"}}, wantErr: false},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if tt.args.secret {
				CypherFileData(tt.args.ctx, tt.args.data)
			}
			if err := DeCypherFileData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("DeCypherFileData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
