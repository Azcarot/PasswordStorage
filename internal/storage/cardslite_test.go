package storage

import (
	"context"
	"database/sql"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardLiteSQL_CreateNewRecord(t *testing.T) {
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
			BCLiteS.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			err := BCLiteS.CreateNewRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("AddCardReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestCardLiteSQL_UpdateRecord(t *testing.T) {
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
			BCLiteS.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			err := BCLiteS.UpdateRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("UpdateReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestCardLiteSQL_DeleteRecord(t *testing.T) {
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
			BCLiteS.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			err := BCLiteS.DeleteRecord(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("DeleteReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestCardLiteSQL_GetAllRecords(t *testing.T) {
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
			BCLiteS.AddData(tt.args.data)
			var ctx context.Context
			if len(tt.args.data.User) != 0 {
				ctx = context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			} else {
				ctx = context.Background()
			}

			_, err := BCLiteS.GetAllRecords(ctx)

			if (err != nil) != tt.args.wantErr {
				t.Errorf("GetAllRecordsReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}

func TestBankCardLiteStorage_GetRecord(t *testing.T) {
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
		{name: "data2", args: args{ctx: context.WithValue(context.Background(), UserLoginCtxKey, "User"), data: BankCardData{ID: 1, CardNumber: "1111"}}, want: BankCardData{ID: 0, CardNumber: "1111"}, wantErr: false},
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
				err := BCLiteS.AddData(tt.args.data)
				assert.NoError(t, err)
				err = BCLiteS.CreateNewRecord(tt.args.ctx)
				assert.NoError(t, err)
			}
			got, err := BCLiteS.GetRecord(tt.args.ctx)
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

func TestBankCardLiteStorage_HashDatabaseData(t *testing.T) {
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
			got, err := BCLiteS.HashDatabaseData(tt.args.ctx)
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
