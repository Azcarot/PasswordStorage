package storage

import (
	"reflect"
	"testing"
)

func TestBankCardStorage_GetData(t *testing.T) {
	tests := []struct {
		name string
		want any
	}{
		{name: "1", want: BankCardData{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BCST.AddData(BankCardData{})
			if got := BCST.GetData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BankCardStorage.GetData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBankCardLiteStorage_GetData(t *testing.T) {
	tests := []struct {
		name string
		want any
	}{
		{name: "1", want: BankCardData{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BCLiteS.AddData(BankCardData{})
			if got := BCLiteS.GetData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BankCardStorage.GetData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextStorage_GetData(t *testing.T) {
	tests := []struct {
		name string
		want any
	}{
		{name: "1", want: TextData{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TST.AddData(TextData{})
			if got := TST.GetData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Text.GetData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextLiteStorage_GetData(t *testing.T) {
	tests := []struct {
		name string
		want any
	}{
		{name: "1", want: TextData{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TLiteS.AddData(TextData{})
			if got := TLiteS.GetData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TextStorage.GetData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLPWStorage_GetData(t *testing.T) {
	tests := []struct {
		name string
		want any
	}{
		{name: "1", want: LoginData{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LPST.AddData(LoginData{})
			if got := LPST.GetData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login.GetData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLPWLiteStorage_GetData(t *testing.T) {
	tests := []struct {
		name string
		want any
	}{
		{name: "1", want: LoginData{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LPWLiteS.AddData(LoginData{})
			if got := LPWLiteS.GetData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoginStorage.GetData() = %v, want %v", got, tt.want)
			}
		})
	}
}
