package cypher

import (
	"context"
	"testing"

	"github.com/Azcarot/PasswordStorage/internal/storage"
)

func TestCypherTextData(t *testing.T) {
	type args struct {
		ctx  context.Context
		data *storage.TextData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{ctx: context.Background(), data: &storage.TextData{ID: 1, Text: "1111", Comment: "22222"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CypherTextData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("CypherTextData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeCypherTextData(t *testing.T) {
	type args struct {
		ctx    context.Context
		data   *storage.TextData
		secret bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{ctx: context.Background(), data: &storage.TextData{ID: 1, Text: "1111", Comment: "22222"}}, wantErr: true},
		{name: "test1", args: args{ctx: context.Background(), secret: true, data: &storage.TextData{ID: 1, Text: "1111 11111 1111 1111", Comment: "22222"}}, wantErr: false},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if tt.args.secret {
				CypherTextData(tt.args.ctx, tt.args.data)
			}
			if err := DeCypherTextData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("DeCypherTextData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
