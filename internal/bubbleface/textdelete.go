package face

import (
	"context"

	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type textDeleteModel struct {
	choices  []string
	datas    []storage.TextData
	cursor   int
	selected map[int]struct{}
}

var textDeleteHeader string = "Please select text to delete"

// NewTextDeleteModel - функция для построения и работы с
// списком сохраненных текстовых данных и их удалением
func NewTextDeleteModel() textDeleteModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var datas []storage.TextData
	data, err := storage.TLiteS.GetAllRecords(ctx)
	if err != nil {
		return textDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}
	choices, datas, err = deCypherText(ctx, data.([]storage.TextData))
	if err != nil {
		textViewHeader = "Text delete error has occured, please try agan"
		return textDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}

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
