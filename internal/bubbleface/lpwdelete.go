package face

import (
	"context"

	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type lpwDeleteModel struct {
	choices  []string
	datas    []storage.LoginResponse
	cursor   int
	selected map[int]struct{}
}

var lpwDeleteHeader string = "Please select text to delete"

// NewLPWDeleteModel - основная функция для построения и работы с
// меню удаления пары логин/пароль
func NewLPWDeleteModel() lpwDeleteModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var datas []storage.LoginResponse
	data, err := storage.LPWLiteS.GetAllRecords(ctx)
	if err != nil {
		return lpwDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx = context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	choices, datas = deCypherLPW(ctx, data.([]storage.LoginResponse))

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
	s := buildView(m, lpwDeleteHeader)

	return s
}

func (m lpwDeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return NewLPWMenuModel(), tea.ClearScreen

		case "enter", " ":
			_, ok := m.selected[m.cursor]

			if ok {
				delete(m.selected, m.cursor)
			} else {

				m.selected[m.cursor] = struct{}{}
				selectedLPW = m.datas[m.cursor]
				var lpwData storage.LoginData

				lpwData.Login = selectedLPW.Login
				lpwData.Password = selectedLPW.Password
				lpwData.Comment = selectedLPW.Comment
				lpwData.ID = selectedLPW.ID
				err := storage.LPWLiteS.AddData(lpwData)
				if err != nil {
					lpwDeleteHeader = "Corrupted login/pw data, please try again"
					return NewLPWDeleteModel(), nil
				}

				pauseTicker <- true

				defer func() { resumeTicker <- true }()

				ok, err := requests.DeleteLPWReq(lpwData)

				if err != nil {
					lpwDeleteHeader = "Something went wrong, please try again"
					return NewLPWDeleteModel(), nil
				}
				if !ok {
					lpwDeleteHeader = "Wrond login/pw data, try again"
					return NewLPWDeleteModel(), nil
				}
				lpwDeleteHeader = "login/pw succsesfully deleted!"
				return NewLPWDeleteModel(), tea.ClearScreen

			}

		}
	}

	return m, nil
}
