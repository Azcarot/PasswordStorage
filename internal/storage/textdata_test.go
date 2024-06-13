// Package storage - описание и реализация всех методов взаимодействия с хранилищами

// на сервере и клиенте

package storage

import (
	"context"
	"testing"
)

func TestTextSQL_CreateNewRecord(t *testing.T) {
	type args struct {
		data    TextData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "Secret", args: args{data: TextData{ID: 1, Text: "11", User: "User"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TST.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			err := TST.CreateNewRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("AddCardReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestTextSQL_UpdateRecord(t *testing.T) {
	type args struct {
		data    TextData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "Secret", args: args{data: TextData{ID: 1, Text: "11", User: "User"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TST.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)

			err := TST.UpdateRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("UpdateReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestTextSQL_DeleteRecord(t *testing.T) {
	type args struct {
		data    TextData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "With ID", args: args{data: TextData{ID: 1, Text: "11", User: "User"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TST.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)

			err := TST.DeleteRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("DeleteReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestTextSQL_GetAllRecords(t *testing.T) {
	type args struct {
		data    TextData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "No login", args: args{data: TextData{ID: 1, Text: "11"}, wantErr: true}},
		{name: "Secret", args: args{data: TextData{ID: 1, Text: "11", User: "User"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TST.AddData(tt.args.data)
			var ctx context.Context
			if len(tt.args.data.User) != 0 {
				ctx = context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			} else {
				ctx = context.Background()
			}

			_, err := TST.GetAllRecords(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("GetAllRecordsReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}
