package face

import (
	tea "github.com/charmbracelet/bubbletea"
)

type authModel struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

type mainMenuModel struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

type MainModel struct {
	currentModel tea.Model
}

type mainMenuCho struct {
	Card  string
	Text  string
	Login string
	File  string
}

var mainChoices = mainMenuCho{
	Card:  "Bank card",
	Text:  "Text data",
	Login: "Login/password",
	File:  "File data",
}

func (m MainModel) Init() tea.Cmd {
	return m.currentModel.Init()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.currentModel, cmd = m.currentModel.Update(msg)
	return m, cmd
}

func (m MainModel) View() string {
	return m.currentModel.View()
}

// NewInitialModel - головное меню регистраци/авторизации
func NewInitialModel() MainModel {
	return MainModel{
		// Список меню авторизации
		currentModel: authModel{choices: []string{"Registration", "Log in"}, selected: make(map[int]struct{})},
	}
}

// MakeTeaProg - создание инстанса баблти
func MakeTeaProg() *tea.Program {
	p := tea.NewProgram(NewInitialModel())
	return p
}

// NewMainMenuModel - пострение главного меню для работы со всеми типами данных
func NewMainMenuModel() mainMenuModel {
	return mainMenuModel{
		// Список основного меню
		choices: []string{mainChoices.Login, mainChoices.Card, mainChoices.Text, mainChoices.File},

		selected: make(map[int]struct{}),
	}
}

func (m authModel) Init() tea.Cmd {
	return nil
}

func (m mainMenuModel) Init() tea.Cmd {
	return nil
}

func (m authModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
				if m.choices[m.cursor] == "Registration" {
					return NewAuthRegModel(), nil
				} else {
					return NewAuthModel(), nil
				}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m authModel) View() string {
	// The header
	newheader := "Authorization\n\n"

	s := buildView(m, newheader)

	return s
}

func (m mainMenuModel) View() string {

	newheader := "Main menu, please select what type of data you want to interact with:\n\n"
	s := buildView(m, newheader)

	return s
}

func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
				if m.choices[m.cursor] == mainChoices.Card {
					return NewCardMenuModel(), nil
				}
				if m.choices[m.cursor] == mainChoices.Text {
					return NewTextMenuModel(), nil
				}
				if m.choices[m.cursor] == mainChoices.Login {
					return NewLPWMenuModel(), nil
				}
				if m.choices[m.cursor] == mainChoices.File {
					return NewFileMenuModel(), nil
				}
			}
		}
	}

	return m, nil
}
