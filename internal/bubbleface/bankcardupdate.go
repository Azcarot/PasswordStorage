package face

import (
	"fmt"

	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type cardUpdateModel struct {
	inputs  []textinput.Model
	focused int
	err     error
}

func UpdateCardModel() cardUpdateModel {
	var inputs []textinput.Model = make([]textinput.Model, 5)
	newheader = "Please insert card data:"
	inputs[ccn] = textinput.New()
	inputs[ccn].Placeholder = "4505 **** **** 1234"
	inputs[ccn].Focus()
	inputs[ccn].CharLimit = 20
	inputs[ccn].Width = 30
	inputs[ccn].Prompt = ""
	inputs[ccn].SetValue(selectedCard.CardNumber)
	inputs[ccn].Validate = ccnValidator

	inputs[exp] = textinput.New()
	inputs[exp].Placeholder = "MM/YY "
	inputs[exp].CharLimit = 5
	inputs[exp].Width = 5
	inputs[exp].Prompt = ""
	inputs[exp].SetValue(selectedCard.ExpDate)
	inputs[exp].Validate = expValidator

	inputs[cvv] = textinput.New()
	inputs[cvv].Placeholder = "XXX"
	inputs[cvv].CharLimit = 3
	inputs[cvv].Width = 5
	inputs[cvv].SetValue(selectedCard.Cvc)
	inputs[cvv].Prompt = ""
	inputs[cvv].Validate = cvvValidator

	inputs[fnm] = textinput.New()
	inputs[fnm].Placeholder = "Ivan ivanov"
	inputs[fnm].CharLimit = 40
	inputs[fnm].Width = 40
	inputs[fnm].SetValue(selectedCard.FullName)
	inputs[fnm].Prompt = ""

	inputs[bcmt] = textinput.New()
	inputs[bcmt].Placeholder = "Some text"
	inputs[bcmt].CharLimit = 100
	inputs[bcmt].Width = 40
	inputs[bcmt].SetValue(selectedCard.Comment)
	inputs[bcmt].Prompt = ""

	return cardUpdateModel{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m cardUpdateModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m cardUpdateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				var req storage.BankCardData
				req.CardNumber = (m.inputs[ccn].Value())
				req.ExpDate = (m.inputs[exp].Value())
				req.Cvc = (m.inputs[cvv].Value())
				req.FullName = (m.inputs[fnm].Value())
				req.Comment = (m.inputs[bcmt].Value())
				req.ID = selectedCard.ID
				ok, err := requests.UpdateCardReq(req)

				if err != nil {
					newheader = "Something went wrong, please try again"
					return CardViewModel(), nil
				}
				if !ok {
					newheader = "Wrond card data, try again"
					return CardViewModel(), nil
				}
				newheader = "Card succsesfully updated!"
				return CardViewModel(), tea.ClearScreen
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		case tea.KeyCtrlB:
			return CardMenuModel(), tea.ClearScreen

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

func (m cardUpdateModel) View() string {
	return fmt.Sprintf(
		` %s

 %s
 %s

 %s  %s  %s
 %s  %s  %s

 %s  %s

 %s


 press ctrl+B to go to previous menu
`,
		newheader,
		inputStyle.Width(30).Render("Card Number"),
		m.inputs[ccn].View(),
		inputStyle.Width(6).Render("EXP"),
		inputStyle.Width(6).Render("CVV"),
		inputStyle.Width(30).Render("Full name"),
		m.inputs[exp].View(),
		m.inputs[cvv].View(),
		m.inputs[fnm].View(),
		inputStyle.Width(30).Render("Comment"),
		m.inputs[bcmt].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"
}

func (m *cardUpdateModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *cardUpdateModel) prevInput() {
	m.focused--

	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
