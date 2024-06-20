package face

import (
	"fmt"

	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type textUpdateModel struct {
	inputs  []textinput.Model
	focused int
	err     error
}

var textUpdateHeader string = "Please check text data:"

// NewUpdateTextModel - основная функция для построения и работы с
// c сохраненными текстовыми данными
func NewUpdateTextModel() textUpdateModel {
	var inputs []textinput.Model = make([]textinput.Model, 2)
	inputs[txt] = textinput.New()
	inputs[txt].Placeholder = "Some text **************"
	inputs[txt].Focus()
	inputs[txt].CharLimit = 100
	inputs[txt].Width = 30
	inputs[txt].Prompt = ""
	inputs[txt].SetValue(selectedText.Text)

	inputs[tcmt] = textinput.New()
	inputs[tcmt].Placeholder = "Some comment"
	inputs[tcmt].CharLimit = 100
	inputs[tcmt].Width = 40
	inputs[tcmt].Prompt = ""
	inputs[tcmt].SetValue(selectedText.Comment)

	return textUpdateModel{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m textUpdateModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textUpdateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				var req storage.TextData
				req.Text = (m.inputs[txt].Value())
				req.Comment = (m.inputs[tcmt].Value())
				req.ID = selectedText.ID
				ok, err := requests.UpdateTextReq(req)

				if err != nil {
					textUpdateHeader = "Something went wrong, please try again"
					return NewTextViewModel(), nil
				}
				if !ok {
					textUpdateHeader = "Wrond text data, try again"
					return NewTextViewModel(), nil
				}
				textUpdateHeader = "Text succsesfully updated!"
				return NewTextViewModel(), tea.ClearScreen
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		case tea.KeyCtrlB:
			return NewTextMenuModel(), tea.ClearScreen

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

func (m textUpdateModel) View() string {
	return fmt.Sprintf(
		` %s

 %s
 %s


 %s  
 %s

 %s


 press ctrl+B to go to previous menu
`,
		textUpdateHeader,
		inputStyle.Width(30).Render("Your Text"),
		m.inputs[txt].View(),
		inputStyle.Width(30).Render("Comment"),
		m.inputs[tcmt].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"
}

func (m *textUpdateModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *textUpdateModel) prevInput() {
	m.focused--

	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
