// Package storage - описание и реализация всех методов взаимодействия с хранилищами

// на сервере и клиенте

package storage

import (
	"context"
	"testing"
)

func TestFileSQL_CreateNewRecord(t *testing.T) {
	type args struct {
		data    FileData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "No secret", args: args{data: FileData{ID: 1, FileName: "11", User: "User", Path: "./Path"}, wantErr: true}},
		{name: "Secret", args: args{data: FileData{ID: 1, FileName: "11", User: "User", Path: "./Path"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FST.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			if len(tt.args.secret) != 0 {
				var b [16]byte
				copy(b[:], Secret)
				ctx = context.WithValue(ctx, EncryptionCtxKey, b)
			}
			err := FST.CreateNewRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("AddCardReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestFileSQL_UpdateRecord(t *testing.T) {
	type args struct {
		data    FileData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "No secret", args: args{data: FileData{ID: 1, FileName: "11", User: "User", Path: "./Path"}, wantErr: true}},
		{name: "Secret", args: args{data: FileData{ID: 1, FileName: "11", User: "User", Path: "./Path"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FST.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			if len(tt.args.secret) != 0 {
				var b [16]byte
				copy(b[:], Secret)
				ctx = context.WithValue(ctx, EncryptionCtxKey, b)
			}
			err := FST.UpdateRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("UpdateReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestFileSQL_DeleteRecord(t *testing.T) {
	type args struct {
		data    FileData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "With ID", args: args{data: FileData{ID: 1, FileName: "11", User: "User", Path: "./Path"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FST.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			if len(tt.args.secret) != 0 {
				var b [16]byte
				copy(b[:], Secret)
				ctx = context.WithValue(ctx, EncryptionCtxKey, b)
			}
			err := FST.DeleteRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("DeleteReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestFileSQL_GetAllRecords(t *testing.T) {
	type args struct {
		data    FileData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "No login", args: args{data: FileData{ID: 1, FileName: "11", Path: "./Path"}, wantErr: true}},
		{name: "No secret", args: args{data: FileData{ID: 1, FileName: "11", User: "User", Path: "./Path"}, wantErr: false}},
		{name: "Secret", args: args{data: FileData{ID: 1, FileName: "11", User: "User", Path: "./Path"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FST.AddData(tt.args.data)
			var ctx context.Context
			if len(tt.args.data.User) != 0 {
				ctx = context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			} else {
				ctx = context.Background()
			}

			if len(tt.args.secret) != 0 {
				var b [16]byte
				copy(b[:], Secret)
				ctx = context.WithValue(ctx, EncryptionCtxKey, b)
			}
			_, err := FST.GetAllRecords(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("GetAllRecordsReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}
