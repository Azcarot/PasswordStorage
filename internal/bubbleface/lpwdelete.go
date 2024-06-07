package face

import (
	"context"
	"fmt"

	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type lpwDeleteModel struct {
	choices  []string
	lpws     []storage.LoginResponse
	cursor   int
	selected map[int]struct{}
}

func LPWDeleteModel() lpwDeleteModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var lpws []storage.LoginResponse
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
	choices, lpws = deCypherLPW(ctx, data.([]storage.LoginResponse))

	return lpwDeleteModel{

		choices:  choices,
		lpws:     lpws,
		selected: make(map[int]struct{}),
	}
}

func (m lpwDeleteModel) Init() tea.Cmd {
	return nil
}

func (m lpwDeleteModel) View() string {
	newheader = "Please select text to delete"
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
			return TextMenuModel(), tea.ClearScreen

		case "enter", " ":
			_, ok := m.selected[m.cursor]

			if ok {
				delete(m.selected, m.cursor)
			} else {

				m.selected[m.cursor] = struct{}{}
				selectedLPW = m.lpws[m.cursor]
				var lpwData storage.LoginData

				lpwData.Login = selectedLPW.Login
				lpwData.Password = selectedLPW.Password
				lpwData.Comment = selectedLPW.Comment
				lpwData.ID = selectedLPW.ID
				err := storage.LPWLiteS.AddData(lpwData)
				if err != nil {
					return LPWDeleteModel(), nil
				}

				pauseTicker <- true

				defer func() { resumeTicker <- true }()

				ok, err := requests.DeleteLPWReq(lpwData)

				if err != nil {
					newheader = "Something went wrong, please try again"
					return LPWDeleteModel(), nil
				}
				if !ok {
					newheader = "Wrond login/pw data, try again"
					return LPWDeleteModel(), nil
				}
				newheader = "login/pw succsesfully deleted!"
				return LPWDeleteModel(), tea.ClearScreen

			}

		}
	}

	return m, nil
}
