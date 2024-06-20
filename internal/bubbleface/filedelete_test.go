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

func TestNewFileDeleteModel(t *testing.T) {
	tests := []struct {
		name string
		want fileDeleteModel
	}{
		{name: "typeerr", want: fileDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}},
		{name: "geterr", want: fileDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}},
		{name: "noerr", want: fileDeleteModel{

			choices: []string{},
			datas:   []storage.FileData{},

			selected: make(map[int]struct{}),
		}},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock := mock_storage.NewMockPgxStorage(ctrl)
		storage.FLiteS = mock
		switch tt.name {
		case "geterr":
			mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return(fileDeleteModel{

				choices: []string{},

				selected: make(map[int]struct{}),
			}, fmt.Errorf("error"))
		default:
			ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
			mock.EXPECT().GetAllRecords(gomock.Any()).Times(1).Return([]storage.FileData{}, nil)
			choices, datas, err := deCypherFile(ctx, []storage.FileData{})
			assert.NoError(t, err)
			tt.want.choices = choices
			tt.want.datas = datas
		}

		t.Run(tt.name, func(t *testing.T) {
			if got := NewFileDeleteModel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileDeleteModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileDeleteModel_Init(t *testing.T) {
	tests := []struct {
		name string
		m    fileDeleteModel
		want tea.Cmd
	}{
		{name: "name", m: fileDeleteModel{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fileDeleteModel.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileDeleteModel_View(t *testing.T) {
	var builder strings.Builder

	builder.WriteString(deleteFileHeader)
	builder.WriteString("\n\n")

	builder.WriteString("")
	builder.WriteString("")
	builder.WriteString("")
	builder.WriteString("\nPress q to quit.\n")

	// Return the complete string
	str := builder.String()
	tests := []struct {
		name string
		m    fileDeleteModel
		want string
	}{
		{name: "name", m: fileDeleteModel{}, want: str},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.View(); got != tt.want {
				t.Errorf("fileDeleteModel.View() = %v, want %v", got, tt.want)
			}
		})
	}
}
