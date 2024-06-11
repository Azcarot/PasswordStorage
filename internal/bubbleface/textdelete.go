package face

import (
	"context"

	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type textDeleteModel struct {
	choices  []string
	datas    []storage.TextResponse
	cursor   int
	selected map[int]struct{}
}

var textDeleteHeader string = "Please select text to delete"

// NewTextDeleteModel - функция для построения и работы с
// списком сохраненных текстовых данных и их удалением
func NewTextDeleteModel() textDeleteModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var datas []storage.TextResponse
	data, err := storage.TLiteS.GetAllRecords(ctx)
	if err != nil {
		return textDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx = context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	choices, datas = deCypherText(ctx, data.([]storage.TextResponse))

	return textDeleteModel{

		choices:  choices,
		datas:    datas,
		selected: make(map[int]struct{}),
	}
}

func (m textDeleteModel) Init() tea.Cmd {
	return nil
}

func (m textDeleteModel) View() string {

	s := buildView(&m, textDeleteHeader)

	return s
}

func (m textDeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return buildUpdate(&textDeleteHeader, msg, &m, NewTextMenuModel(), updateTextDeleteModel)
}
