// Package storage - описание и реализация всех методов взаимодействия с хранилищами

// на сервере и клиенте

package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Azcarot/PasswordStorage/internal/cfg"
	"github.com/jackc/pgx/v5"
)

func TestMain(m *testing.M) {
	// os.Exit skips defer calls
	// so we need to call another function
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	var f cfg.Flags
	f.FlagDBAddr = "host='localhost' user='postgres' password='12345' sslmode=disable"
	DB, err = pgx.Connect(context.Background(), f.FlagDBAddr)
	if err != nil {
		//handle the error
		log.Fatal(err)
	}
	ST = MakeConn(DB)
	dbName := "testdb"
	ctx := context.Background()
	_, err = DB.Exec(ctx, "DROP DATABASE IF EXISTS "+dbName)
	if err != nil {
		//handle the error
		log.Fatal(err)
	}
	_, err = DB.Exec(ctx, "create database "+dbName)
	if err != nil {
		//handle the error
		log.Fatal(err)
	}
	ST.CreateTablesForGoKeeper()
	_, err = os.Stat("test.db")
	if err == nil {
		err = os.Remove("test.db")
		if err != nil {
			log.Fatal(err)
		}
	}
	LiteDB, err = sql.Open("sqlite", "test.db")
	LiteST = MakeLiteConn(LiteDB)
	LiteST.CreateTablesForGoKeeper()

	// truncates all test data after the tests are run
	defer func() {

		_, _ = DB.Exec(ctx, fmt.Sprintf("DELETE FROM %s", "bank_card"))
		_, _ = DB.Exec(ctx, fmt.Sprintf("DELETE FROM %s", "file_data"))
		_, _ = DB.Exec(ctx, fmt.Sprintf("DELETE FROM %s", "text_data"))
		_, _ = DB.Exec(ctx, fmt.Sprintf("DELETE FROM %s", "login_pw"))
		LiteDB.Close()
		os.Remove("test.db")

		DB.Close(ctx)
	}()
	BCST = NewBCStorage(BCST, DB)
	BCLiteS = NewBCLiteStorage(BCLiteS, LiteDB)
	FST = NewFSTtorage(FST, DB)
	FLiteS = NewFLiteStorage(FLiteS, LiteDB)
	LPST = NewLPSTtorage(LPST, DB)
	LPWLiteS = NewLPLiteStorage(LPWLiteS, LiteDB)
	TST = NewTSTtorage(TST, DB)
	TLiteS = NewTLiteStorage(TLiteS, LiteDB)
	return m.Run(), nil
}

func TestCardSQL_CreateNewRecord(t *testing.T) {
	type args struct {
		data    BankCardData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "Secret", args: args{data: BankCardData{ID: 1, CardNumber: "11", User: "User", ExpDate: "ExpDate"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			BCST.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)

			err := BCST.CreateNewRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("AddCardReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestCardSQL_UpdateRecord(t *testing.T) {
	type args struct {
		data    BankCardData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "Secret", args: args{data: BankCardData{ID: 1, CardNumber: "11", User: "User", ExpDate: "ExpDate"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BCST.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)

			err := BCST.UpdateRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("UpdateReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestCardSQL_DeleteRecord(t *testing.T) {
	type args struct {
		data    BankCardData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "With ID", args: args{data: BankCardData{ID: 1, CardNumber: "11", User: "User", ExpDate: "ExpDate"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BCST.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)

			err := BCST.DeleteRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("DeleteReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestCardSQL_GetAllRecords(t *testing.T) {
	type args struct {
		data    BankCardData
		wantErr bool
		secret  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "No login", args: args{data: BankCardData{ID: 1, CardNumber: "11", ExpDate: "ExpDate"}, wantErr: true}},
		{name: "Secret", args: args{data: BankCardData{ID: 1, CardNumber: "11", User: "User", ExpDate: "ExpDate"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BCST.AddData(tt.args.data)
			var ctx context.Context
			if len(tt.args.data.User) != 0 {
				ctx = context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			} else {
				ctx = context.Background()
			}

			_, err := BCST.GetAllRecords(ctx)

			if (err != nil) != tt.args.wantErr {
				t.Errorf("GetAllRecordsReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestNewConn(t *testing.T) {
	type args struct {
		f cfg.Flags
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "normal data",
			args: args{f: cfg.Flags{FlagDBAddr: "host='localhost' user='postgres' password='12345' sslmode=disable"}}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewConn(tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("NewConn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
