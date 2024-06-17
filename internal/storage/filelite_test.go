// Package storage - описание и реализация всех методов взаимодействия с хранилищами

// на сервере и клиенте

package storage

import (
	"context"
	"reflect"
	"testing"
)

func TestFileLiteSQL_CreateNewRecord(t *testing.T) {
	type args struct {
		data    FileData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "Secret", args: args{data: FileData{ID: 1, FileName: "11", User: "User", Path: "./Path"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FLiteS.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)

			err := FLiteS.CreateNewRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("AddCardReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestFileLiteSQL_UpdateRecord(t *testing.T) {
	type args struct {
		data    FileData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "Secret", args: args{data: FileData{ID: 1, FileName: "11", User: "User", Path: "./Path"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FLiteS.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)

			err := FLiteS.UpdateRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("UpdateReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestFileLiteSQL_DeleteRecord(t *testing.T) {
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
			FLiteS.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)

			err := FLiteS.DeleteRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("DeleteReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestFileLiteSQL_GetAllRecords(t *testing.T) {
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
		{name: "Secret", args: args{data: FileData{ID: 1, FileName: "11", User: "User", Path: "./Path"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FLiteS.AddData(tt.args.data)
			var ctx context.Context
			if len(tt.args.data.User) != 0 {
				ctx = context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			} else {
				ctx = context.Background()
			}

			_, err := FLiteS.GetAllRecords(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("GetAllRecordsReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestFileLiteStorage_GetRecord(t *testing.T) {
	type args struct {
		ctx  context.Context
		data FileData
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{name: "no login", args: args{ctx: context.Background()}, want: nil, wantErr: true},
		{name: "with login, no rows", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User")}, want: FileData{}, wantErr: true},
		{name: "data", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User"), data: FileData{ID: 1, FileName: "1111"}}, want: FileData{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := FLiteS.GetRecord(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileLiteStorage.GetRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileLiteStorage.GetRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileLiteStorage_HashDatabaseData(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "1", args: args{ctx: context.Background()}, want: "", wantErr: true},
		{name: "2", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User")}, want: "4f53cda18c2baa0c0354bb5f9a3ecbe5ed12ab4d8e11ba873c2f11161202b945", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FLiteS.HashDatabaseData(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileLiteStorage.HashDatabaseData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FileLiteStorage.HashDatabaseData() = %v, want %v", got, tt.want)
			}
		})
	}
}
