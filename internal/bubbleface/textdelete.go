package face

import (
	"context"

	"github.com/Azcarot/PasswordStorage/internal/requests"
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

	s := buildView(m, textDeleteHeader)

	return s
}

func (m textDeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "ctrl+b":
			return NewTextMenuModel(), tea.ClearScreen

		case "enter", " ":
			_, ok := m.selected[m.cursor]

			if ok {
				delete(m.selected, m.cursor)
			} else {

				m.selected[m.cursor] = struct{}{}
				selectedText = m.datas[m.cursor]
				var textData storage.TextData

				textData.Text = selectedText.Text
				textData.Comment = selectedText.Comment
				textData.ID = selectedText.ID
				err := storage.TLiteS.AddData(textData)
				if err != nil {
					textDeleteHeader = "Corrupted text data, please try again"
					return NewTextDeleteModel(), nil
				}

				pauseTicker <- true

				defer func() { resumeTicker <- true }()

				ok, err := requests.DeleteTextReq(textData)

				if err != nil {
					textDeleteHeader = "Something went wrong, please try again"
					return NewTextDeleteModel(), nil
				}
				if !ok {
					textDeleteHeader = "Wrond text data, try again"
					return NewTextDeleteModel(), nil
				}
				textDeleteHeader = "Text succsesfully deleted!"
				return NewTextDeleteModel(), tea.ClearScreen

			}

		}
	}

	return m, nil
}
