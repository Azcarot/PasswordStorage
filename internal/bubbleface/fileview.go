package face

import (
	"context"

	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type fileViewModel struct {
	choices  []string
	cursor   int
	datas    []storage.FileData
	selected map[int]struct{}
}

var selectedFile storage.FileData
var fileViewHeader string = "Main file menu, please chose your file:\n\n"

// NewFileViewModel - основная функция для построения и просмотра
// списка сохраненных файлов
func NewFileViewModel() fileViewModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var files []storage.FileData
	data, err := storage.FLiteS.GetAllRecords(ctx)
	if err != nil {
		return fileViewModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}

	choices, files, err = deCypherFile(ctx, data.([]storage.FileData))
	if err != nil {
		fileViewHeader = "File view error occured, please try again"
		return fileViewModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}
	return fileViewModel{

		choices:  choices,
		datas:    files,
		selected: make(map[int]struct{}),
	}
}

func (m fileViewModel) Init() tea.Cmd {
	return nil
}

func (m fileViewModel) View() string {

	s := buildView(&m, fileViewHeader)

	return s
}

func (m fileViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return buildUpdate(&fileViewHeader, msg, &m, NewFileMenuModel(), updateFileViewModel)
}
