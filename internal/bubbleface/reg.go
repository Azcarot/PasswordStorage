package face

import (
	"fmt"
	"time"

	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	login = iota
	pwd
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
	newheader     = "Please enter your credentials"
)

type (
	errMsg error
)

type regmodel struct {
	inputs  []textinput.Model
	header  string
	focused int
	err     error
}

func AuthRegModel() regmodel {
	var inputs []textinput.Model = make([]textinput.Model, 2)
	inputs[login] = textinput.New()
	inputs[login].Placeholder = "Ivan ivanov"
	inputs[login].Focus()
	inputs[login].CharLimit = 20
	inputs[login].Width = 30
	inputs[login].Prompt = ""

	inputs[pwd] = textinput.New()
	inputs[pwd].Placeholder = "123qwe "
	inputs[pwd].CharLimit = 30
	inputs[pwd].Width = 20
	inputs[pwd].Prompt = ""

	return regmodel{
		inputs:  inputs,
		focused: 0,
		err:     nil,
		header:  newheader,
	}
}

func (m regmodel) Init() tea.Cmd {
	return textinput.Blink
}

func (m regmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				var req storage.RegisterRequest
				req.Login = (m.inputs[login].Value())
				req.Password = (m.inputs[pwd].Value())
				ok, err := requests.RegistrationReq(req)
				if err != nil {
					newheader = "Please enter your credentials"
					return AuthRegModel(), nil
				}
				if !ok {
					newheader = "User already exists"
					return AuthRegModel(), tea.ClearScreen
				}
				storage.UserLoginPw.Login = req.Login
				storage.UserLoginPw.Password = req.Password
				ticker := time.NewTicker(2 * time.Second)
				quit := make(chan struct{})
				go func() {
					for {
						select {
						case <-ticker.C:
							requests.SyncCardReq()
						case <-quit:
							ticker.Stop()
							return
						}
					}
				}()
				return MainMenuModel(), nil
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()

		case tea.KeyCtrlB:
			return InitialModel(), tea.ClearScreen

		}

		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m regmodel) View() string {
	return fmt.Sprintf(
		` %s,
 press ctrl+b to go to the previous menu:

 %s
 %s

 %s  
 %s

 %s
`,
		m.header,
		inputStyle.Width(30).Render("Login"),
		m.inputs[login].View(),
		inputStyle.Width(30).Render("Password"),
		m.inputs[pwd].View(),

		continueStyle.Render("Continue"),
	) + "\n"
}

// nextInput focuses the next input field
func (m *regmodel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *regmodel) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
