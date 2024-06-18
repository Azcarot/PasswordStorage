// Package storage - описание и реализация всех методов взаимодействия с хранилищами

// на сервере и клиенте

package storage

import (
	"context"
	"database/sql"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextLiteSQL_CreateNewRecord(t *testing.T) {
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
			TLiteS.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)

			err := TLiteS.CreateNewRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("AddText() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestTextLiteSQL_UpdateRecord(t *testing.T) {
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
			TLiteS.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)

			err := TLiteS.UpdateRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("UpdateReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestTextLiteSQL_DeleteRecord(t *testing.T) {
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
			TLiteS.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)

			err := TLiteS.DeleteRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("DeleteReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestTextLiteSQL_GetAllRecords(t *testing.T) {
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
			TLiteS.AddData(tt.args.data)
			var ctx context.Context
			if len(tt.args.data.User) != 0 {
				ctx = context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			} else {
				ctx = context.Background()
			}

			_, err := TLiteS.GetAllRecords(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("GetAllRecordsReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestTextLiteStorage_GetRecord(t *testing.T) {
	type args struct {
		ctx  context.Context
		data TextData
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{name: "no login", args: args{ctx: context.Background()}, want: nil, wantErr: true},
		{name: "with login, no rows", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User")}, want: TextData{}, wantErr: true},
		{name: "data", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User"), data: TextData{ID: 1, Text: "1111"}}, want: TextData{}, wantErr: true},
		{name: "data2", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User"), data: TextData{ID: 1, Text: "1111"}}, want: TextData{ID: 0, Text: "1111"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "data2" {
				LiteDB.Close()
				os.Remove("test.db")
				LiteDB, _ = sql.Open("sqlite", "test.db")
				LiteST = MakeLiteConn(LiteDB)
				LiteST.CreateTablesForGoKeeper()
				BCLiteS = NewBCLiteStorage(BCLiteS, LiteDB)
				FLiteS = NewFLiteStorage(FLiteS, LiteDB)
				LPWLiteS = NewLPLiteStorage(LPWLiteS, LiteDB)
				TLiteS = NewTLiteStorage(TLiteS, LiteDB)
				err := TLiteS.AddData(tt.args.data)
				assert.NoError(t, err)
				err = TLiteS.CreateNewRecord(tt.args.ctx)
				assert.NoError(t, err)
			}
			got, err := TLiteS.GetRecord(tt.args.ctx)
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

func TestTextLiteStorage_HashDatabaseData(t *testing.T) {
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
		{name: "2", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User")}, want: "3b6927ff4a11ac33052eab23863faa068fa6888b87c18c711faffe41a9a6c13d", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TLiteS.HashDatabaseData(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("TextLiteStorage.HashDatabaseData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TextLiteStorage.HashDatabaseData() = %v, want %v", got, tt.want)
			}
		})
	}
}
