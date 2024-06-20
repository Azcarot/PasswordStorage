package face

import (
	"fmt"

	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	txt = iota
	tcmt
)

var textMenuHeader string = "Main text menu, please chose your action:\n\n"
var textAddHeader string = "Please insert text data:"

type textMenuModel struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}
type textCho struct {
	Add    string
	View   string
	Delete string
}

var textChoices = textCho{
	Add:    "Add New Text",
	View:   "View/Update text data",
	Delete: "Delete text",
}

// NewTextMenuModel - основное меню работы с текстовыми данными
func NewTextMenuModel() textMenuModel {
	return textMenuModel{

		choices: []string{textChoices.Add, textChoices.View, textChoices.Delete},

		selected: make(map[int]struct{}),
	}
}

func (m textMenuModel) Init() tea.Cmd {
	return nil
}

func (m textMenuModel) View() string {

	s := buildView(&m, textMenuHeader)

	return s
}

func (m textMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return buildUpdate(&textMenuHeader, msg, &m, NewMainMenuModel(), updateTextMenuModel)
}

type textModel struct {
	inputs  []textinput.Model
	focused int
	err     error
}

// NewAddTextModel - меню для добавления новых текстовых данных
func NewAddTextModel() textModel {
	var inputs []textinput.Model = make([]textinput.Model, 2)
	inputs[txt] = textinput.New()
	inputs[txt].Placeholder = "Some text **************"
	inputs[txt].Focus()
	inputs[txt].CharLimit = 100
	inputs[txt].Width = 30
	inputs[txt].Prompt = ""

	inputs[tcmt] = textinput.New()
	inputs[tcmt].Placeholder = "Some comment"
	inputs[tcmt].CharLimit = 100
	inputs[tcmt].Width = 40
	inputs[tcmt].Prompt = ""

	return textModel{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m textModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				var req storage.TextData
				req.Text = (m.inputs[txt].Value())
				req.Comment = (m.inputs[tcmt].Value())
				ok, err := requests.AddTextReq(req)

				if err != nil {
					textAddHeader = "Something went wrong, please try again"
					return NewAddTextModel(), nil
				}
				if !ok {
					textAddHeader = "Wrond data, please try again"
					return NewAddTextModel(), nil
				}
				textAddHeader = "Text succsesfully saved!"
				return NewAddTextModel(), tea.ClearScreen
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

func (m textModel) View() string {
	return fmt.Sprintf(
		` %s

 %s
 %s


 %s  
 %s

 %s


 press ctrl+B to go to previous menu
`,
		textAddHeader,
		inputStyle.Width(30).Render("Your Text"),
		m.inputs[txt].View(),
		inputStyle.Width(30).Render("Comment"),
		m.inputs[tcmt].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"
}

func (m *textModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *textModel) prevInput() {
	m.focused--

	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
