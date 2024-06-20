package face

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azcarot/PasswordStorage/internal/cypher"
	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *cardDeleteModel) GetChoices() []string {
	return m.choices
}

func (m *cardDeleteModel) GetCursor() int {
	return m.cursor
}

func (m *cardDeleteModel) SetCursor(c int) {
	m.cursor = c
}

func (m *cardDeleteModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *cardDeleteModel) GetData() any {
	return m.datas
}

func (m *cardMenuModel) GetChoices() []string {
	return m.choices
}

func (m *cardMenuModel) GetCursor() int {
	return m.cursor
}

func (m *cardMenuModel) SetCursor(c int) {
	m.cursor = c
}

func (m *cardMenuModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *cardMenuModel) GetData() any {
	return nil
}

func (m *cardViewModel) GetChoices() []string {
	return m.choices
}

func (m *cardViewModel) GetCursor() int {
	return m.cursor
}

func (m *cardViewModel) SetCursor(c int) {
	m.cursor = c
}

func (m *cardViewModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *cardViewModel) GetData() any {
	return m.datas
}

func buildView(m basicModel, header string) string {
	var builder strings.Builder

	// Write the header
	builder.WriteString(header)
	builder.WriteString("\n\n")

	// Get choices, cursor, and selected from the interface
	choices := m.GetChoices()
	cursor := m.GetCursor()
	selected := m.GetSelected()

	// Iterate through the choices
	for i, choice := range choices {
		cursorSymbol := " "
		if cursor == i {
			cursorSymbol = ">"
		}

		checked := " "
		if _, ok := selected[i]; ok {
			checked = "x"
		}

		// Write the formatted string to the builder
		builder.WriteString(cursorSymbol)
		builder.WriteString(checked)
		builder.WriteString(choice)
		builder.WriteString("\n")

	}

	// Write the quit instruction
	builder.WriteString("\nPress q to quit.\n")

	// Return the complete string
	return builder.String()
}

func buildUpdate(h *string, msg tea.Msg, m basicModel, r tea.Model, f func(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd)) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		ans := "buildUpdate error, please try again"
		*h = ans
		return r, nil
	}

	choices := m.GetChoices()
	cursor := m.GetCursor()
	selected := m.GetSelected()
	datas := m.GetData()

	switch key.String() {

	case "ctrl+c", "q":
		return m.(tea.Model), tea.Quit

	case "up", "k":
		if cursor > 0 {
			cursor--
			m.SetCursor(cursor)
		}

	case "down", "j":
		if cursor < len(choices)-1 {
			cursor++
			m.SetCursor(cursor)
		}

	case "ctrl+b":
		return r, nil

	case "enter", " ":
		_, ok := selected[cursor]

		if ok {
			delete(selected, cursor)
		} else {
			return f(selected, cursor, datas, choices)
		}
	}
	return m.(tea.Model), nil
}

func updateDeleteCardModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	selected[cursor] = struct{}{}
	cards := datas.([]storage.BankCardData)
	selectedCard = cards[cursor]
	var bankData storage.BankCardData
	bankData.CardNumber = selectedCard.CardNumber
	bankData.Cvc = selectedCard.Cvc
	bankData.Comment = selectedCard.Comment
	bankData.ID = selectedCard.ID
	bankData.ExpDate = selectedCard.ExpDate
	bankData.FullName = selectedCard.FullName
	err := storage.BCLiteS.AddData(bankData)
	if err != nil {
		cardDeleteHeader = "Corrupted card data, please try again"
		return NewCardDeleteModel(), nil
	}

	pauseTicker <- true

	defer func() { resumeTicker <- true }()

	ok, err := requests.DeleteCardReq(bankData)

	if err != nil {
		cardDeleteHeader = "Something went wrong, please try again"
		return NewCardDeleteModel(), nil
	}
	if !ok {
		cardDeleteHeader = "Wrond card data, try again"
		return NewCardDeleteModel(), nil
	}
	cardDeleteHeader = "Card succsesfully deleted!"
	return NewCardDeleteModel(), nil
}

func updateMenuCardModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	selected[cursor] = struct{}{}
	if choices[cursor] == bankChoices.Add {
		return NewAddCardModel(), nil
	}

	if choices[cursor] == bankChoices.View {
		return NewCardViewModel(), nil
	}
	if choices[cursor] == bankChoices.Delete {
		return NewCardDeleteModel(), nil
	}
	return NewCardMenuModel(), nil
}

func updateViewCardModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	data := datas.([]storage.BankCardData)
	selected[cursor] = struct{}{}
	selectedCard = data[cursor]
	return NewUpdateCardModel(), nil
}

func deCypherBankCard(ctx context.Context, cards []storage.BankCardData) ([]string, []storage.BankCardData, error) {
	var err error
	var choices []string
	var datas []storage.BankCardData
	for _, data := range cards {

		err = cypher.DeCypherBankData(ctx, &data)
		if err != nil {
			return choices, datas, err
		}
		str := fmt.Sprintf("CCN: %s \nCVC: %s   ExpDate: %s   FillName: %s\nComment: %s", data.CardNumber, data.Cvc, data.ExpDate, data.FullName, data.Comment)
		choices = append(choices, str)
		datas = append(datas, data)

	}

	return choices, datas, nil
}
