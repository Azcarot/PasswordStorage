// Package face - модуль взаимодействия с клиентом посредством bubbletea
package face

import (
	"reflect"
	"testing"

	mock_storage "github.com/Azcarot/PasswordStorage/internal/mock"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/golang/mock/gomock"
)

func TestNewCardDeleteModel(t *testing.T) {
	tests := []struct {
		name string
		want cardDeleteModel
	}{
		{name: "name", want: cardDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock := mock_storage.NewMockPgxStorage(ctrl)
		storage.BCLiteS = mock
		mock.EXPECT().GetAllRecords(gomock.Any()).Times(1)
		mock.EXPECT().AddData(gomock.Any()).Times(1)
		storage.BCLiteS.AddData(storage.BankCardData{ID: 1, CardNumber: "12121212"})
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCardDeleteModel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCardDeleteModel() = %v, want %v", got, tt.want)
			}
		})
	}
}
