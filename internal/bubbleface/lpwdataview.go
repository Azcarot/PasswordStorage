package face

import (
	"context"

	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type lpwViewModel struct {
	choices  []string
	cursor   int
	datas    []storage.LoginData
	selected map[int]struct{}
}

var selectedLPW storage.LoginData
var lpwViewHeader string = "Main login/pw menu, please chose your action:\n\n"

// NewLPWViewModel - основная функция для построения и просмотра списка
// сохраненных данных типа логин/пароль
func NewLPWViewModel() lpwViewModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var lpws []storage.LoginData
	data, err := storage.LPWLiteS.GetAllRecords(ctx)
	if err != nil {
		return lpwViewModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}

	choices, lpws, err = deCypherLPW(ctx, data.([]storage.LoginData))
	if err != nil {
		lpwViewHeader = "Login/password view error has occured, please try agan"
		return lpwViewModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}

	return lpwViewModel{

		choices:  choices,
		datas:    lpws,
		selected: make(map[int]struct{}),
	}
}

func (m lpwViewModel) Init() tea.Cmd {
	return nil
}

func (m lpwViewModel) View() string {

	s := buildView(&m, lpwViewHeader)

	return s
}

func (m lpwViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return buildUpdate(&lpwViewHeader, msg, &m, NewLPWMenuModel(), updateLPWViewModel)
}
