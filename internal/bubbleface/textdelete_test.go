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

func TestNewTextDeleteModel(t *testing.T) {
	tests := []struct {
		name string
		want textDeleteModel
	}{
		{name: "typeerr", want: textDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}},
		{name: "geterr", want: textDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}},
		{name: "noerr", want: textDeleteModel{

			choices: []string{},
			datas:   []storage.TextData{},

			selected: make(map[int]struct{}),
		}},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock := mock_storage.NewMockPgxStorage(ctrl)
		storage.TLiteS = mock
		switch tt.name {
		case "geterr":
			mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return(textDeleteModel{

				choices: []string{},

				selected: make(map[int]struct{}),
			}, fmt.Errorf("error"))
		default:
			ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
			mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return([]storage.TextData{}, nil)
			choices, datas, err := deCypherText(ctx, []storage.TextData{})
			assert.NoError(t, err)
			tt.want.choices = choices
			tt.want.datas = datas
		}

		t.Run(tt.name, func(t *testing.T) {
			if got := NewTextDeleteModel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewtextDeleteModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_textDeleteModel_Init(t *testing.T) {
	tests := []struct {
		name string
		m    textDeleteModel
		want tea.Cmd
	}{
		{name: "name", m: textDeleteModel{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("textDeleteModel.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_textDeleteModel_View(t *testing.T) {
	var builder strings.Builder

	builder.WriteString(textDeleteHeader)
	builder.WriteString("\n\n")

	builder.WriteString("")
	builder.WriteString("")
	builder.WriteString("")
	builder.WriteString("\nPress q to quit.\n")

	// Return the complete string
	str := builder.String()
	tests := []struct {
		name string
		m    textDeleteModel
		want string
	}{
		{name: "name", m: textDeleteModel{}, want: str},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.View(); got != tt.want {
				t.Errorf("textDeleteModel.View() = %v, want %v", got, tt.want)
			}
		})
	}
}
