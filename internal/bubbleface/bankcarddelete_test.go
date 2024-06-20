// Package face - модуль взаимодействия с клиентом посредством bubbletea
package face

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	mock_storage "github.com/Azcarot/PasswordStorage/internal/mock"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewCardDeleteModel(t *testing.T) {
	tests := []struct {
		name string
		want cardDeleteModel
	}{
		{name: "typeerr", want: cardDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}},
		{name: "geterr", want: cardDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}},
		{name: "noerr", want: cardDeleteModel{

			choices: []string{},
			datas:   []storage.BankCardData{},

			selected: make(map[int]struct{}),
		}},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock := mock_storage.NewMockPgxStorage(ctrl)
		storage.BCLiteS = mock
		switch tt.name {
		case "geterr":
			mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return(cardDeleteModel{

				choices: []string{},

				selected: make(map[int]struct{}),
			}, fmt.Errorf("error"))

		case "typeerr":
			mock.EXPECT().GetAllRecords(gomock.Any()).Times(1)
		default:
			ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
			mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return([]storage.BankCardData{}, nil)
			choices, datas, err := deCypherBankCard(ctx, []storage.BankCardData{})
			assert.NoError(t, err)
			tt.want.choices = choices
			tt.want.datas = datas
		}

		t.Run(tt.name, func(t *testing.T) {
			if got := NewCardDeleteModel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCardDeleteModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cardDeleteModel_Init(t *testing.T) {
	tests := []struct {
		name string
		m    cardDeleteModel
		want tea.Cmd
	}{
		{name: "name", m: cardDeleteModel{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cardDeleteModel.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cardDeleteModel_View(t *testing.T) {
	var builder strings.Builder

	builder.WriteString(cardDeleteHeader)
	builder.WriteString("\n\n")

	builder.WriteString("")
	builder.WriteString("")
	builder.WriteString("")
	builder.WriteString("\nPress q to quit.\n")

	// Return the complete string
	str := builder.String()
	tests := []struct {
		name string
		m    cardDeleteModel
		want string
	}{
		{name: "name", m: cardDeleteModel{}, want: str},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.View(); got != tt.want {
				t.Errorf("cardDeleteModel.View() = %v, want %v", got, tt.want)
			}
		})
	}
}
