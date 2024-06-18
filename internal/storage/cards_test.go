// Package storage - описание и реализация всех методов взаимодействия с хранилищами

// на сервере и клиенте

package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/Azcarot/PasswordStorage/internal/cfg"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "12345",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return -1, err
	}
	defer postgresC.Terminate(ctx)

	host, err := postgresC.Host(ctx)
	if err != nil {
		return -1, err
	}
	port, err := postgresC.MappedPort(ctx, "5432")
	if err != nil {
		return -1, err
	}

	dsn := fmt.Sprintf("postgres://postgres:12345@%s:%s/testdb?sslmode=disable", host, port.Port())
	DB, err = pgx.Connect(ctx, dsn)
	if err != nil {
		return -1, err
	}
	defer DB.Close(ctx)

	ST = MakeConn(DB)

	ST.CreateTablesForGoKeeper()

	sqliteFile := "test.db"
	_, err = os.Stat(sqliteFile)
	if err == nil {
		err = os.Remove(sqliteFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	LiteDB, err = sql.Open("sqlite", "test.db")
	defer func() {
		LiteDB.Close()
		os.Remove(sqliteFile)
	}()

	LiteST = MakeLiteConn(LiteDB)
	LiteST.CreateTablesForGoKeeper()
	BCST = NewBCStorage(BCST, DB)
	BCLiteS = NewBCLiteStorage(BCLiteS, LiteDB)
	FST = NewFSTtorage(FST, DB)
	FLiteS = NewFLiteStorage(FLiteS, LiteDB)
	LPST = NewLPSTtorage(LPST, DB)
	LPWLiteS = NewLPLiteStorage(LPWLiteS, LiteDB)
	TST = NewTSTtorage(TST, DB)
	TLiteS = NewTLiteStorage(TLiteS, LiteDB)
	code = m.Run()

	_, _ = DB.Exec(ctx, "DELETE FROM bank_card")
	_, _ = DB.Exec(ctx, "DELETE FROM file_data")
	_, _ = DB.Exec(ctx, "DELETE FROM text_data")
	_, _ = DB.Exec(ctx, "DELETE FROM login_pw")

	return code, nil
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
		{name: "NoLogin", args: args{data: BankCardData{ID: 1, CardNumber: "11", User: "User", ExpDate: "ExpDate"}, secret: "secret", wantErr: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			BCST.AddData(tt.args.data)
			ctx := context.Background()
			if tt.name != "NoLogin" {
				ctx = context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			}

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

func TestBankCardStorage_GetRecord(t *testing.T) {
	type args struct {
		ctx  context.Context
		data BankCardData
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{name: "no login", args: args{ctx: context.Background()}, want: nil, wantErr: true},
		{name: "with login, no rows", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User")}, want: BankCardData{}, wantErr: true},
		{name: "data", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User"), data: BankCardData{ID: 1, CardNumber: "1111"}}, want: BankCardData{}, wantErr: true},
		{name: "data2", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User"), data: BankCardData{ID: 1, CardNumber: "1111"}}, want: BankCardData{ID: 1, CardNumber: "1111"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "data2" {
				_, _ = DB.Exec(tt.args.ctx, fmt.Sprintf("DELETE FROM %s", "bank_card"))
				ST.CreateTablesForGoKeeper()
				err := BCST.AddData(tt.args.data)
				assert.NoError(t, err)
				err = BCST.CreateNewRecord(tt.args.ctx)
				assert.NoError(t, err)
			}
			got, err := BCST.GetRecord(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("BankCardStorage.GetRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BankCardStorage.GetRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBankCardStorage_HashDatabaseData(t *testing.T) {
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
		{name: "2", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User")}, want: "2b2f27b5391999b0d4a14bcfee24f647ca22b792b0105974ef73c801d0e32bfb", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BCST.HashDatabaseData(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("BankCardStorage.HashDatabaseData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BankCardStorage.HashDatabaseData() = %v, want %v", got, tt.want)
			}
		})
	}
}
