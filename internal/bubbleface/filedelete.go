package face

import (
	"context"

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

	s := buildView(&m, deleteFileHeader)

	return s
}

func (m fileDeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return buildUpdate(&deleteFileHeader, msg, &m, NewFileMenuModel(), updateFileDeleteModel)
}
