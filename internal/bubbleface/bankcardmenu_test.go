package face

import (
	"reflect"
	"testing"
)

func TestNewCardMenuModel(t *testing.T) {
	tests := []struct {
		name string
		want cardMenuModel
	}{
		{name: "name", want: cardMenuModel{choices: []string{bankChoices.Add, bankChoices.View, bankChoices.Delete}, selected: make(map[int]struct{})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCardMenuModel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCardMenuModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ccnValidator(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "name", args: args{s: "1111 2222"}, wantErr: false},
		{name: "name2", args: args{s: "1111 22asd"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ccnValidator(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ccnValidator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
