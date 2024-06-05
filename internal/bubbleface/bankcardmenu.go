package face

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	ccn = iota
	exp
	cvv
	fnm
	bcmt
)

type cardMenuModel struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}
type bankCardCho struct {
	Add    string
	View   string
	Delete string
}

var bankChoices = bankCardCho{
	Add:    "Add New Card",
	View:   "View/Update card data",
	Delete: "Delete card",
}

func CardMenuModel() cardMenuModel {
	return cardMenuModel{

		choices: []string{bankChoices.Add, bankChoices.View, bankChoices.Delete},

		selected: make(map[int]struct{}),
	}
}

func (m cardMenuModel) Init() tea.Cmd {
	return nil
}

func (m cardMenuModel) View() string {
	s := "Main bank card menu, please chose your action:\n\n"

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

func (m cardMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				if m.choices[m.cursor] == bankChoices.Add {
					return AddCardModel(), nil
				}
				if m.choices[m.cursor] == bankChoices.View {
					return CardViewModel(), nil
				}
			}
		}
	}

	return m, nil
}

type cardModel struct {
	inputs  []textinput.Model
	focused int
	err     error
}

func ccnValidator(s string) error {

	if len(s) > 16+3 {
		return fmt.Errorf("CCN is too long")
	}

	if len(s) == 0 || len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
		return fmt.Errorf("CCN is invalid")
	}

	if len(s)%5 == 0 && s[len(s)-1] != ' ' {
		return fmt.Errorf("CCN must separate groups with spaces")
	}

	c := strings.ReplaceAll(s, " ", "")
	_, err := strconv.ParseInt(c, 10, 64)

	return err
}

func expValidator(s string) error {

	e := strings.ReplaceAll(s, "/", "")
	_, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		return fmt.Errorf("EXP is invalid")
	}

	if len(s) >= 3 && (strings.Index(s, "/") != 2 || strings.LastIndex(s, "/") != 2) {
		return fmt.Errorf("EXP is invalid")
	}

	return nil
}

func cvvValidator(s string) error {

	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

func AddCardModel() cardModel {
	var inputs []textinput.Model = make([]textinput.Model, 5)
	newheader = "Please insert card data:"
	inputs[ccn] = textinput.New()
	inputs[ccn].Placeholder = "4505 **** **** 1234"
	inputs[ccn].Focus()
	inputs[ccn].CharLimit = 20
	inputs[ccn].Width = 30
	inputs[ccn].Prompt = ""
	inputs[ccn].Validate = ccnValidator

	inputs[exp] = textinput.New()
	inputs[exp].Placeholder = "MM/YY "
	inputs[exp].CharLimit = 5
	inputs[exp].Width = 5
	inputs[exp].Prompt = ""
	inputs[exp].Validate = expValidator

	inputs[cvv] = textinput.New()
	inputs[cvv].Placeholder = "XXX"
	inputs[cvv].CharLimit = 3
	inputs[cvv].Width = 5
	inputs[cvv].Prompt = ""
	inputs[cvv].Validate = cvvValidator

	inputs[fnm] = textinput.New()
	inputs[fnm].Placeholder = "Ivan ivanov"
	inputs[fnm].CharLimit = 40
	inputs[fnm].Width = 40
	inputs[fnm].Prompt = ""

	inputs[bcmt] = textinput.New()
	inputs[bcmt].Placeholder = "Some text"
	inputs[bcmt].CharLimit = 100
	inputs[bcmt].Width = 40
	inputs[bcmt].Prompt = ""

	return cardModel{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m cardModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m cardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				ok, err := requests.AddCardReq(req)

				if err != nil {
					newheader = "Something went wrong, please try again"
					return AddCardModel(), nil
				}
				if !ok {
					newheader = "Wrond card data, try again"
					return AddCardModel(), nil
				}
				newheader = "Card succsesfully saved!"
				return AddCardModel(), tea.ClearScreen
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

func (m cardModel) View() string {
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

func (m *cardModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *cardModel) prevInput() {
	m.focused--

	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
