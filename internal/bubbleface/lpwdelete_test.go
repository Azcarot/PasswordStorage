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

func TestNewLPWDeleteModel(t *testing.T) {
	tests := []struct {
		name string
		want lpwDeleteModel
	}{
		{name: "typeerr", want: lpwDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}},
		{name: "geterr", want: lpwDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}},
		{name: "noerr", want: lpwDeleteModel{

			choices: []string{},
			datas:   []storage.LoginData{},

			selected: make(map[int]struct{}),
		}},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock := mock_storage.NewMockPgxStorage(ctrl)
		storage.LPWLiteS = mock
		switch tt.name {
		case "geterr":
			mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return(lpwDeleteModel{

				choices: []string{},

				selected: make(map[int]struct{}),
			}, fmt.Errorf("error"))
		default:
			ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
			mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return([]storage.LoginData{}, nil)
			choices, datas, err := deCypherLPW(ctx, []storage.LoginData{})
			assert.NoError(t, err)
			tt.want.choices = choices
			tt.want.datas = datas
		}

		t.Run(tt.name, func(t *testing.T) {
			if got := NewLPWDeleteModel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewlpwDeleteModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lpwDeleteModel_Init(t *testing.T) {
	tests := []struct {
		name string
		m    lpwDeleteModel
		want tea.Cmd
	}{
		{name: "name", m: lpwDeleteModel{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lpwDeleteModel.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lpwDeleteModel_View(t *testing.T) {
	var builder strings.Builder

	builder.WriteString(lpwDeleteHeader)
	builder.WriteString("\n\n")

	builder.WriteString("")
	builder.WriteString("")
	builder.WriteString("")
	builder.WriteString("\nPress q to quit.\n")

	// Return the complete string
	str := builder.String()
	tests := []struct {
		name string
		m    lpwDeleteModel
		want string
	}{
		{name: "name", m: lpwDeleteModel{}, want: str},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.View(); got != tt.want {
				t.Errorf("lpwDeleteModel.View() = %v, want %v", got, tt.want)
			}
		})
	}
}
