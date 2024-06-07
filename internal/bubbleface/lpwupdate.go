package face

import (
	"fmt"

	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type lpwUpdateModel struct {
	inputs  []textinput.Model
	focused int
	err     error
}

func UpdateLPWModel() lpwUpdateModel {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	newheader = "Please check login/password data:"
	inputs[lgn] = textinput.New()
	inputs[lgn].Placeholder = "Login **************"
	inputs[lgn].Focus()
	inputs[lgn].CharLimit = 100
	inputs[lgn].Width = 30
	inputs[lgn].Prompt = ""
	inputs[lgn].SetValue(selectedLPW.Login)

	inputs[psw] = textinput.New()
	inputs[psw].Placeholder = "Password **************"
	inputs[psw].CharLimit = 100
	inputs[psw].Width = 30
	inputs[psw].Prompt = ""
	inputs[psw].SetValue(selectedLPW.Password)

	inputs[lcmt] = textinput.New()
	inputs[lcmt].Placeholder = "Some comment"
	inputs[lcmt].CharLimit = 100
	inputs[lcmt].Width = 40
	inputs[lcmt].Prompt = ""
	inputs[lcmt].SetValue(selectedLPW.Comment)

	return lpwUpdateModel{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m lpwUpdateModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m lpwUpdateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				req.ID = selectedLPW.ID
				ok, err := requests.UpdateLPWReq(req)

				if err != nil {
					newheader = "Something went wrong, please try again"
					return LPWViewModel(), nil
				}
				if !ok {
					newheader = "Wrond text data, try again"
					return LPWViewModel(), nil
				}
				newheader = "Login/password succsesfully updated!"
				return LPWViewModel(), tea.ClearScreen
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		case tea.KeyCtrlB:
			return LPWMenuModel(), tea.ClearScreen

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

func (m lpwUpdateModel) View() string {
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
		newheader,
		inputStyle.Width(30).Render("Your Login"),
		m.inputs[lgn].View(),
		inputStyle.Width(30).Render("Your Password"),
		m.inputs[psw].View(),
		inputStyle.Width(30).Render("Comment"),
		m.inputs[lcmt].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"
}

func (m *lpwUpdateModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *lpwUpdateModel) prevInput() {
	m.focused--

	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
