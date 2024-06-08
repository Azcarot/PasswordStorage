package face

import (
	"context"
	"fmt"

	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type textDeleteModel struct {
	choices  []string
	texts    []storage.TextResponse
	cursor   int
	selected map[int]struct{}
}

// TextDeleteModel - функция для построения и работы с
// списком сохраненных текстовых данных и их удалением
func TextDeleteModel() textDeleteModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var texts []storage.TextResponse
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
	choices, texts = deCypherText(ctx, data.([]storage.TextResponse))

	return textDeleteModel{

		choices:  choices,
		texts:    texts,
		selected: make(map[int]struct{}),
	}
}

func (m textDeleteModel) Init() tea.Cmd {
	return nil
}

func (m textDeleteModel) View() string {
	newheader = "Please select text to delete"
	s := newheader + "\n\n"

	for i, choice := range m.choices {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

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
			return TextMenuModel(), tea.ClearScreen

		case "enter", " ":
			_, ok := m.selected[m.cursor]

			if ok {
				delete(m.selected, m.cursor)
			} else {

				m.selected[m.cursor] = struct{}{}
				selectedText = m.texts[m.cursor]
				var textData storage.TextData

				textData.Text = selectedText.Text
				textData.Comment = selectedText.Comment
				textData.ID = selectedText.ID
				err := storage.TLiteS.AddData(textData)
				if err != nil {
					return TextDeleteModel(), nil
				}

				pauseTicker <- true

				defer func() { resumeTicker <- true }()

				ok, err := requests.DeleteTextReq(textData)

				if err != nil {
					newheader = "Something went wrong, please try again"
					return TextDeleteModel(), nil
				}
				if !ok {
					newheader = "Wrond text data, try again"
					return TextDeleteModel(), nil
				}
				newheader = "Text succsesfully deleted!"
				return TextDeleteModel(), tea.ClearScreen

			}

		}
	}

	return m, nil
}
