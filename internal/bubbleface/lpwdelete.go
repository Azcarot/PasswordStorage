package face

import (
	"context"

	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type lpwDeleteModel struct {
	choices  []string
	datas    []storage.LoginData
	cursor   int
	selected map[int]struct{}
}

var lpwDeleteHeader string = "Please select text to delete"

// NewLPWDeleteModel - основная функция для построения и работы с
// меню удаления пары логин/пароль
func NewLPWDeleteModel() lpwDeleteModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var datas []storage.LoginData
	data, err := storage.LPWLiteS.GetAllRecords(ctx)
	if err != nil {
		return lpwDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}
	choices, datas, err = deCypherLPW(ctx, data.([]storage.LoginData))
	if err != nil {
		lpwViewHeader = "Login/password delete error has occured, please try agan"
		return lpwDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}

	return lpwDeleteModel{

		choices:  choices,
		datas:    datas,
		selected: make(map[int]struct{}),
	}
}

func (m lpwDeleteModel) Init() tea.Cmd {
	return nil
}

func (m lpwDeleteModel) View() string {
	lpwDeleteHeader = "Please select text to delete"
	s := buildView(&m, lpwDeleteHeader)

	return s
}

func (m lpwDeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return buildUpdate(&lpwDeleteHeader, msg, &m, NewLPWMenuModel(), updateLPWDeleteModel)
}
