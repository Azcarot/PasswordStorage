package storage

import (
	"context"
	"testing"
)

func TestSQLStore_CreateNewUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		data UserData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "1", args: args{ctx: context.Background(), data: UserData{Login: "111", Password: "222"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ST.CreateNewUser(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("SQLStore.CreateNewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLStore_CheckUserExists(t *testing.T) {
	type args struct {
		ctx  context.Context
		data UserData
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "1", args: args{ctx: context.Background(), data: UserData{Login: "111", Password: "222"}}, wantErr: false, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ST.CreateNewUser(context.Background(), tt.args.data)
			got, err := ST.CheckUserExists(tt.args.ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLStore.CheckUserExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SQLStore.CheckUserExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLStore_CheckUserPassword(t *testing.T) {
	type args struct {
		ctx  context.Context
		data UserData
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "1", args: args{ctx: context.Background(), data: UserData{Login: "111", Password: "222"}}, wantErr: false, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ST.CreateNewUser(context.Background(), tt.args.data)
			got, err := ST.CheckUserPassword(tt.args.ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLStore.CheckUserPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SQLStore.CheckUserPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
