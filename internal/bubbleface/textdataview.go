package face

import (
	"context"
	"fmt"

	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type textViewModel struct {
	choices  []string
	cursor   int
	datas    []storage.TextResponse
	selected map[int]struct{}
}

var selectedText storage.TextResponse
var textViewHeader string = "Main text menu, please chose your action:\n\n"

// NewTextViewModel - функция для построения и работы со
// списком сохраненных текстовых данных
func NewTextViewModel() textViewModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var texts []storage.TextResponse
	data, err := storage.TLiteS.GetAllRecords(ctx)
	if err != nil {
		return textViewModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx = context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	choices, texts = deCypherText(ctx, data.([]storage.TextResponse))

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

func deCypherText(ctx context.Context, cards []storage.TextResponse) ([]string, []storage.TextResponse) {
	var err error
	var choices []string
	var datas []storage.TextResponse
	for _, data := range cards {
		data.Text, err = storage.Dechypher(ctx, data.Text)
		if err != nil {
			continue
		}

		data.Comment, err = storage.Dechypher(ctx, data.Comment)
		if err != nil {
			continue
		}

		str := fmt.Sprintf("Your text: %s \nComment: %s", data.Text, data.Comment)
		choices = append(choices, str)
		datas = append(datas, data)
	}
	return choices, datas
}
