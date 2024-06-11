package face

import (
	"context"
	"fmt"

	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type lpwViewModel struct {
	choices  []string
	cursor   int
	datas    []storage.LoginResponse
	selected map[int]struct{}
}

var selectedLPW storage.LoginResponse
var lpwViewHeader string = "Main login/pw menu, please chose your action:\n\n"

// NewLPWViewModel - основная функция для построения и просмотра списка
// сохраненных данных типа логин/пароль
func NewLPWViewModel() lpwViewModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var lpws []storage.LoginResponse
	data, err := storage.LPWLiteS.GetAllRecords(ctx)
	if err != nil {
		return lpwViewModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx = context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	choices, lpws = deCypherLPW(ctx, data.([]storage.LoginResponse))

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

	s := buildView(m, lpwViewHeader)

	return s
}

func (m lpwViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return NewLPWMenuModel(), nil
		case "enter", " ":
			_, ok := m.selected[m.cursor]

			if ok {
				delete(m.selected, m.cursor)
			} else {

				m.selected[m.cursor] = struct{}{}
				selectedLPW = m.datas[m.cursor]
				return NewUpdateLPWModel(), nil

			}
		}
	}

	return m, nil
}

func deCypherLPW(ctx context.Context, cards []storage.LoginResponse) ([]string, []storage.LoginResponse) {
	var err error
	var choices []string
	var datas []storage.LoginResponse
	for _, data := range cards {
		data.Login, err = storage.Dechypher(ctx, data.Login)
		if err != nil {
			continue
		}
		data.Password, err = storage.Dechypher(ctx, data.Password)
		if err != nil {
			continue
		}
		data.Comment, err = storage.Dechypher(ctx, data.Comment)
		if err != nil {
			continue
		}

		str := fmt.Sprintf("Your login: %s \nYour password: %s \nComment: %s", data.Login, data.Password, data.Comment)
		choices = append(choices, str)
		datas = append(datas, data)
	}
	return choices, datas
}
