package face

import (
	"context"

	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type textViewModel struct {
	choices  []string
	cursor   int
	datas    []storage.TextData
	selected map[int]struct{}
}

var selectedText storage.TextData
var textViewHeader string = "Main text menu, please chose your action:\n\n"

// NewTextViewModel - функция для построения и работы со
// списком сохраненных текстовых данных
func NewTextViewModel() textViewModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var texts []storage.TextData
	data, err := storage.TLiteS.GetAllRecords(ctx)
	if err != nil {
		return textViewModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}
	choices, texts, err = deCypherText(ctx, data.([]storage.TextData))
	if err != nil {
		textViewHeader = "Text view error has occured, please try agan"
		return textViewModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}

	return textViewModel{

		choices:  choices,
		datas:    texts,
		selected: make(map[int]struct{}),
	}
}

func (m textViewModel) Init() tea.Cmd {
	return nil
}

func (m textViewModel) View() string {

	s := buildView(&m, textViewHeader)

	return s
}

func (m textViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return buildUpdate(&textViewHeader, msg, &m, NewTextMenuModel(), updateTextViewModel)
}
