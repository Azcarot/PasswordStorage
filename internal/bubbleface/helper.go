package face

import (
	tea "github.com/charmbracelet/bubbletea"
)

type basicModel interface {
	GetChoices() []string
	GetCursor() int
	SetCursor(int)
	GetSelected() map[int]struct{}
	GetData() any
}

func (m *authModel) GetChoices() []string {
	return m.choices
}

func (m *authModel) GetCursor() int {
	return m.cursor
}

func (m *authModel) SetCursor(c int) {
	m.cursor = c
}

func (m *authModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *authModel) GetData() any {
	return nil
}

func (m *mainMenuModel) GetChoices() []string {
	return m.choices
}

func (m *mainMenuModel) GetCursor() int {
	return m.cursor
}

func (m *mainMenuModel) SetCursor(c int) {
	m.cursor = c
}

func (m *mainMenuModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *mainMenuModel) GetData() any {
	return nil
}

func updateMainMenu(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	selected[cursor] = struct{}{}
	if choices[cursor] == mainChoices.Card {
		return NewCardMenuModel(), nil
	}
	if choices[cursor] == mainChoices.Text {
		return NewTextMenuModel(), nil
	}
	if choices[cursor] == mainChoices.Login {
		return NewLPWMenuModel(), nil
	}
	if choices[cursor] == mainChoices.File {
		return NewFileMenuModel(), nil
	}
	return NewMainMenuModel(), nil
}
