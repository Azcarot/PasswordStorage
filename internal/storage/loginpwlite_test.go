// Package storage - описание и реализация всех методов взаимодействия с хранилищами

// на сервере и клиенте

package storage

import (
	"context"
	"testing"
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
		{name: "No secret", args: args{data: LoginData{ID: 1, Login: "11", User: "User", Password: "./Password"}, wantErr: false}},
		{name: "Secret", args: args{data: LoginData{ID: 1, Login: "11", User: "User", Password: "./Password"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LPWLiteS.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			if len(tt.args.secret) != 0 {
				var b [16]byte
				copy(b[:], Secret)
				ctx = context.WithValue(ctx, EncryptionCtxKey, b)
			}
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
		{name: "No secret", args: args{data: LoginData{ID: 1, Login: "11", User: "User", Password: "./Password"}, wantErr: true}},
		{name: "Secret", args: args{data: LoginData{ID: 1, Login: "11", User: "User", Password: "./Password"}, secret: "secret", wantErr: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LPWLiteS.AddData(tt.args.data)
			ctx := context.WithValue(context.Background(), UserLoginCtxKey, tt.args.data.User)
			if len(tt.args.secret) != 0 {
				var b [16]byte
				copy(b[:], Secret)
				ctx = context.WithValue(ctx, EncryptionCtxKey, b)
			}
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
			if len(tt.args.secret) != 0 {
				var b [16]byte
				copy(b[:], Secret)
				ctx = context.WithValue(ctx, EncryptionCtxKey, b)
			}
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
		{name: "No secret", args: args{data: LoginData{ID: 1, Login: "11", User: "User", Password: "./Password"}, wantErr: false}},
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

			if len(tt.args.secret) != 0 {
				var b [16]byte
				copy(b[:], Secret)
				ctx = context.WithValue(ctx, EncryptionCtxKey, b)
			}
			_, err := LPWLiteS.GetAllRecords(ctx)
			if (err != nil) != tt.args.wantErr {
				t.Errorf("GetAllRecordsReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
				return
			}
		})
	}
}
