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
	datas    []storage.FileResponse
	selected map[int]struct{}
}

var selectedFile storage.FileResponse
var fileViewHeader string = "Main file menu, please chose your file:\n\n"

// NewFileViewModel - основная функция для построения и просмотра
// списка сохраненных файлов
func NewFileViewModel() fileViewModel {
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
