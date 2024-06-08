package face

import (
	"context"
	"fmt"

	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type fileDeleteModel struct {
	choices  []string
	files    []storage.FileResponse
	cursor   int
	selected map[int]struct{}
}

// FileDeleteModel - основная функция для построения списка и удалением
// файла
func FileDeleteModel() fileDeleteModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var files []storage.FileResponse
	data, err := storage.FLiteS.GetAllRecords(ctx)
	if err != nil {
		return fileDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx = context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	choices, files = deCypherFile(ctx, data.([]storage.FileResponse))

	return fileDeleteModel{

		choices:  choices,
		files:    files,
		selected: make(map[int]struct{}),
	}
}

func (m fileDeleteModel) Init() tea.Cmd {
	return nil
}

func (m fileDeleteModel) View() string {
	newheader = "Please select file to delete"
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

func (m fileDeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return FileMenuModel(), tea.ClearScreen

		case "enter", " ":
			_, ok := m.selected[m.cursor]

			if ok {
				delete(m.selected, m.cursor)
			} else {

				m.selected[m.cursor] = struct{}{}
				selectedFile = m.files[m.cursor]
				var fileData storage.FileData

				fileData.FileName = selectedFile.FileName
				fileData.Path = selectedFile.Path
				fileData.Comment = selectedFile.Comment
				fileData.ID = selectedFile.ID
				err := storage.FLiteS.AddData(fileData)
				if err != nil {
					return FileDeleteModel(), nil
				}

				pauseTicker <- true

				defer func() { resumeTicker <- true }()

				ok, err := requests.DeleteFileReq(fileData)

				if err != nil {
					newheader = "Something went wrong, please try again"
					return FileDeleteModel(), nil
				}
				if !ok {
					newheader = "Wrond file data, try again"
					return FileDeleteModel(), nil
				}
				newheader = "File succsesfully deleted!"
				return FileDeleteModel(), tea.ClearScreen

			}

		}
	}

	return m, nil
}
