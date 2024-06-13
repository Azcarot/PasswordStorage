// Package storage - описание и реализация всех методов взаимодействия с хранилищами

// на сервере и клиенте

package storage

import (
	"context"
	"testing"
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
				t.Errorf("AddCardReq() error = %v, wantErr %v, test %v", err, tt.args.wantErr, tt.name)
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
