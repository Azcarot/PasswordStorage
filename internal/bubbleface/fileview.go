package face

import (
	"context"
	"fmt"

	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type fileViewModel struct {
	choices  []string
	cursor   int
	files    []storage.FileResponse
	selected map[int]struct{}
}

var selectedFile storage.FileResponse

func FileViewModel() fileViewModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var files []storage.FileResponse
	data, err := storage.FLiteS.GetAllRecords(ctx)
	if err != nil {
		return fileViewModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx = context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	choices, files = deCypherFile(ctx, data.([]storage.FileResponse))

	return fileViewModel{

		choices:  choices,
		files:    files,
		selected: make(map[int]struct{}),
	}
}

func (m fileViewModel) Init() tea.Cmd {
	return nil
}

func (m fileViewModel) View() string {
	s := "Main file menu, please chose your file:\n\n"

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

func (m fileViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return FileMenuModel(), nil
		case "enter", " ":
			_, ok := m.selected[m.cursor]

			if ok {
				delete(m.selected, m.cursor)
			} else {

				m.selected[m.cursor] = struct{}{}
				selectedFile = m.files[m.cursor]
				return UpdateFileModel(), nil

			}
		}
	}

	return m, nil
}

func deCypherFile(ctx context.Context, cards []storage.FileResponse) ([]string, []storage.FileResponse) {
	var err error
	var choices []string
	var datas []storage.FileResponse
	for _, data := range cards {
		data.FileName, err = storage.Dechypher(ctx, data.FileName)
		if err != nil {
			continue
		}
		data.Path, err = storage.Dechypher(ctx, data.Path)
		if err != nil {
			continue
		}
		data.Data, err = storage.Dechypher(ctx, data.Data)
		if err != nil {
			continue
		}
		data.Comment, err = storage.Dechypher(ctx, data.Comment)
		if err != nil {
			continue
		}

		str := fmt.Sprintf("Your file name: %s \nYour file path: %s \nComment: %s", data.FileName, data.Path, data.Comment)
		choices = append(choices, str)
		datas = append(datas, data)
	}
	return choices, datas
}
