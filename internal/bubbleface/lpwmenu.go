package face

import (
	"fmt"

	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	lgn = iota
	psw
	lcmt
)

type lpwMenuModel struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}
type lpwCho struct {
	Add    string
	View   string
	Delete string
}

var lpwChoices = lpwCho{
	Add:    "Add New login/password",
	View:   "View/Updatelogin/password data",
	Delete: "Delete login/password",
}

var lpwMenuHeader string = "Main login/password menu, please chose your action:\n\n"
var lpwAddHeader string = "Please insert login/password data:"

// NewLPWMenuModel - основная функция для построения и работы с
// основным меню действий с данными типа логин/пароль
func NewLPWMenuModel() lpwMenuModel {
	return lpwMenuModel{

		choices: []string{lpwChoices.Add, lpwChoices.View, lpwChoices.Delete},

		selected: make(map[int]struct{}),
	}
}

func (m lpwMenuModel) Init() tea.Cmd {
	return nil
}

func (m lpwMenuModel) View() string {

	s := buildView(m, lpwMenuHeader)

	return s
}

func (m lpwMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return NewMainMenuModel(), nil

		case "enter", " ":
			_, ok := m.selected[m.cursor]

			if ok {
				delete(m.selected, m.cursor)
			} else {

				m.selected[m.cursor] = struct{}{}
				if m.choices[m.cursor] == lpwChoices.Add {
					return NewAddLPWModel(), nil
				}
				if m.choices[m.cursor] == lpwChoices.View {
					return NewLPWViewModel(), nil
				}
				if m.choices[m.cursor] == lpwChoices.Delete {
					return NewLPWDeleteModel(), nil
				}
			}
		}
	}

	return m, nil
}

type lpwModel struct {
	inputs  []textinput.Model
	focused int
	err     error
}

func NewAddLPWModel() lpwModel {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	inputs[lgn] = textinput.New()
	inputs[lgn].Placeholder = "Login **************"
	inputs[lgn].Focus()
	inputs[lgn].CharLimit = 100
	inputs[lgn].Width = 30
	inputs[lgn].Prompt = ""

	inputs[psw] = textinput.New()
	inputs[psw].Placeholder = "Password **************"
	inputs[psw].CharLimit = 100
	inputs[psw].Width = 30
	inputs[psw].Prompt = ""

	inputs[lcmt] = textinput.New()
	inputs[lcmt].Placeholder = "Some comment"
	inputs[lcmt].CharLimit = 100
	inputs[lcmt].Width = 40
	inputs[lcmt].Prompt = ""

	return lpwModel{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m lpwModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m lpwModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				var req storage.LoginData
				req.Login = (m.inputs[lgn].Value())
				req.Password = (m.inputs[psw].Value())
				req.Comment = (m.inputs[lcmt].Value())
				ok, err := requests.AddLPWReq(req)

				if err != nil {
					lpwAddHeader = "Something went wrong, please try again"
					return NewAddLPWModel(), nil
				}
				if !ok {
					lpwAddHeader = "Wrond data, please try again"
					return NewAddLPWModel(), nil
				}
				lpwAddHeader = "Text succsesfully saved!"
				return NewAddLPWModel(), tea.ClearScreen
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		case tea.KeyCtrlB:
			return NewLPWMenuModel(), tea.ClearScreen

		}

		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m lpwModel) View() string {
	return fmt.Sprintf(
		` %s

 %s
 %s

 %s
 %s

 %s
 %s

 %s


 press ctrl+B to go to previous menu
`,
		lpwAddHeader,
		inputStyle.Width(30).Render("Your login"),
		m.inputs[lgn].View(),
		inputStyle.Width(30).Render("Your password"),
		m.inputs[psw].View(),
		inputStyle.Width(30).Render("Comment"),
		m.inputs[lcmt].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"
}

func (m *lpwModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *lpwModel) prevInput() {
	m.focused--

	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
