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

func TestLPWLiteSQL_CreateNewRecord(t *testing.T) {
	type args struct {
		data    LoginData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "Secret", args: args{data: LoginData{ID: 1, Login: "11", User: "User", Password: "./Password"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LPWLiteS.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			err := LPWLiteS.CreateNewRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("AddLPWReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestLPWLiteSQL_UpdateRecord(t *testing.T) {
	type args struct {
		data    LoginData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "Secret", args: args{data: LoginData{ID: 1, Login: "11", User: "User", Password: "./Password"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LPWLiteS.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)

			err := LPWLiteS.UpdateRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("UpdateReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestLPWLiteSQL_DeleteRecord(t *testing.T) {
	type args struct {
		data    LoginData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "With ID", args: args{data: LoginData{ID: 1, Login: "11", User: "User", Password: "./Password"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LPWLiteS.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			err := LPWLiteS.DeleteRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("DeleteReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestLPWLiteSQL_GetAllRecords(t *testing.T) {
	type args struct {
		data    LoginData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "No login", args: args{data: LoginData{ID: 1, Login: "11", Password: "./Password"}, wantErr: true}},
		{name: "Secret", args: args{data: LoginData{ID: 1, Login: "11", User: "User", Password: "./Password"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LPWLiteS.AddData(tt.args.data)
			var ctx context.Context
			if len(tt.args.data.User) != 0 {
				ctx = context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			} else {
				ctx = context.Background()
			}

			_, err := LPWLiteS.GetAllRecords(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("GetAllRecordsReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestLPWLiteStorage_GetRecord(t *testing.T) {
	type args struct {
		ctx  context.Context
		data LoginData
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{name: "no login", args: args{ctx: context.Background()}, want: nil, wantErr: true},
		{name: "with login, no rows", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User")}, want: LoginData{}, wantErr: true},
		{name: "data", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User"), data: LoginData{ID: 1, Login: "1111"}}, want: LoginData{}, wantErr: true},
		{name: "data2", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User"), data: LoginData{ID: 1, Login: "1111"}}, want: LoginData{ID: 0, Login: "1111"}, wantErr: false},
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
				err := LPWLiteS.AddData(tt.args.data)
				assert.NoError(t, err)
				err = LPWLiteS.CreateNewRecord(tt.args.ctx)
				assert.NoError(t, err)
			}
			got, err := LPWLiteS.GetRecord(tt.args.ctx)
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

func TestLPWLiteStorage_HashDatabaseData(t *testing.T) {
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
		{name: "2", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User")}, want: "460b73975b97e180f75306c05978ed9d2107e7cd8108f975461f7e87356cb1ff", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LPWLiteS.HashDatabaseData(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("LPWLiteStorage.HashDatabaseData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LPWLiteStorage.HashDatabaseData() = %v, want %v", got, tt.want)
			}
		})
	}
}
