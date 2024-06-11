package face

import (
	"context"

	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type fileDeleteModel struct {
	choices  []string
	datas    []storage.FileResponse
	cursor   int
	selected map[int]struct{}
}

var deleteFileHeader string = "Please select file to delete"

// NewFileDeleteModel - основная функция для построения списка и удалением
// файла
func NewFileDeleteModel() fileDeleteModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)

	var choices []string
	var datas []storage.FileResponse
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
	choices, datas = deCypherFile(ctx, data.([]storage.FileResponse))

	return fileDeleteModel{

		choices:  choices,
		datas:    datas,
		selected: make(map[int]struct{}),
	}
}

func (m fileDeleteModel) Init() tea.Cmd {
	return nil
}

func (m fileDeleteModel) View() string {

	s := buildView(m, deleteFileHeader)

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
			return NewFileMenuModel(), tea.ClearScreen

		case "enter", " ":
			_, ok := m.selected[m.cursor]

			if ok {
				delete(m.selected, m.cursor)
			} else {

				m.selected[m.cursor] = struct{}{}
				selectedFile = m.datas[m.cursor]
				var fileData storage.FileData

				fileData.FileName = selectedFile.FileName
				fileData.Path = selectedFile.Path
				fileData.Comment = selectedFile.Comment
				fileData.ID = selectedFile.ID
				err := storage.FLiteS.AddData(fileData)
				if err != nil {
					deleteFileHeader = "Corrupted file data, please try again"
					return NewFileDeleteModel(), nil
				}

				pauseTicker <- true

				defer func() { resumeTicker <- true }()

				ok, err := requests.DeleteFileReq(fileData)

				if err != nil {
					deleteFileHeader = "Something went wrong, please try again"
					return NewFileDeleteModel(), nil
				}
				if !ok {
					deleteFileHeader = "Wrond file data, try again"
					return NewFileDeleteModel(), nil
				}
				deleteFileHeader = "File succsesfully deleted!"
				return NewFileDeleteModel(), tea.ClearScreen

			}

		}
	}

	return m, nil
}
