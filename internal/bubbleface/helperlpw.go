package face

import (
	"context"
	"fmt"

	"github.com/Azcarot/PasswordStorage/internal/cypher"
	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *lpwDeleteModel) GetChoices() []string {
	return m.choices
}

func (m *lpwDeleteModel) GetCursor() int {
	return m.cursor
}

func (m *lpwDeleteModel) SetCursor(c int) {
	m.cursor = c
}

func (m *lpwDeleteModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *lpwDeleteModel) GetData() any {
	return m.datas
}

func (m *lpwMenuModel) GetChoices() []string {
	return m.choices
}

func (m *lpwMenuModel) GetCursor() int {
	return m.cursor
}
func (m *lpwMenuModel) SetCursor(c int) {
	m.cursor = c
}

func (m *lpwMenuModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *lpwMenuModel) GetData() any {
	return nil
}

func (m *lpwViewModel) GetChoices() []string {
	return m.choices
}

func (m *lpwViewModel) GetCursor() int {
	return m.cursor
}

func (m *lpwViewModel) SetCursor(c int) {
	m.cursor = c
}

func (m *lpwViewModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *lpwViewModel) GetData() any {
	return m.datas
}

func updateLPWViewModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	data := datas.([]storage.LoginData)
	selected[cursor] = struct{}{}
	selectedLPW = data[cursor]
	return NewUpdateLPWModel(), nil
}

func updateLPWDeleteModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	selected[cursor] = struct{}{}
	data := datas.([]storage.LoginData)
	selectedLPW = data[cursor]
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

func updateLPWMenuModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	selected[cursor] = struct{}{}
	if choices[cursor] == lpwChoices.Add {
		return NewAddLPWModel(), nil
	}
	if choices[cursor] == lpwChoices.View {
		return NewLPWViewModel(), nil
	}
	if choices[cursor] == lpwChoices.Delete {
		return NewLPWDeleteModel(), nil
	}
	return NewLPWMenuModel(), nil
}

func deCypherLPW(ctx context.Context, cards []storage.LoginData) ([]string, []storage.LoginData, error) {
	var err error
	var choices []string
	var datas []storage.LoginData
	for _, data := range cards {
		err = cypher.DeCypherLPWData(ctx, &data)
		if err != nil {
			return choices, datas, err
		}

		str := fmt.Sprintf("Your login: %s \nYour password: %s \nComment: %s", data.Login, data.Password, data.Comment)
		choices = append(choices, str)
		datas = append(datas, data)
	}
	return choices, datas, nil
}
