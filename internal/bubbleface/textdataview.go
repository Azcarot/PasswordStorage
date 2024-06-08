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
	texts    []storage.TextResponse
	selected map[int]struct{}
}

var selectedText storage.TextResponse

// TextViewModel - функция для построения и работы со
// списком сохраненных текстовых данных
func TextViewModel() textViewModel {
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
		texts:    texts,
		selected: make(map[int]struct{}),
	}
}

func (m textViewModel) Init() tea.Cmd {
	return nil
}

func (m textViewModel) View() string {
	s := "Main text menu, please chose your action:\n\n"

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

func (m textViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return TextMenuModel(), nil
		case "enter", " ":
			_, ok := m.selected[m.cursor]

			if ok {
				delete(m.selected, m.cursor)
			} else {

				m.selected[m.cursor] = struct{}{}
				selectedText = m.texts[m.cursor]
				return UpdateTextModel(), nil

			}
		}
	}

	return m, nil
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
