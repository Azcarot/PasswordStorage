// Package storage - описание и реализация всех методов взаимодействия с хранилищами

// на сервере и клиенте

package storage

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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
		{name: "Secret", args: args{data: FileData{ID: 1, FileName: "11", User: "User", Path: "./Path"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FST.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)

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
		{name: "Secret", args: args{data: FileData{ID: 1, FileName: "11", User: "User", Path: "./Path"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FST.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)

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
		{name: "Secret", args: args{data: FileData{ID: 1, FileName: "11", User: "User", Path: "./Path"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FST.AddData(tt.args.data)
			FST.CreateNewRecord(context.Background())
			FST.AddData(FileData{ID: 2, FileName: "333", Path: "path"})
			FST.CreateNewRecord(context.Background())

			var ctx context.Context
			if len(tt.args.data.User) != 0 {
				ctx = context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			} else {
				ctx = context.Background()
			}

			_, err := FST.GetAllRecords(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("GetAllRecordsReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}
func TestFileStorage_GetRecord(t *testing.T) {
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
		{name: "data2", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User"), data: FileData{ID: 1, FileName: "1111"}}, want: FileData{ID: 0, FileName: "1111"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "data2" {
				_, _ = DB.Exec(tt.args.ctx, fmt.Sprintf("DELETE FROM %s", "bank_card"))
				ST.CreateTablesForGoKeeper()
				err := FST.AddData(tt.args.data)
				assert.NoError(t, err)
				err = FST.CreateNewRecord(tt.args.ctx)
				assert.NoError(t, err)
			}
			got, err := FST.GetRecord(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStorage.GetRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileStorage.GetRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileStorage_HashDatabaseData(t *testing.T) {
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
		{name: "2", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User")}, want: "58bd7bd63b18cda75e8ede6f2842699f5eca3f2bfec5b0bcbefc06f3e794b448", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FST.HashDatabaseData(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStorage.HashDatabaseData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FileStorage.HashDatabaseData() = %v, want %v", got, tt.want)
			}
		})
	}
}
