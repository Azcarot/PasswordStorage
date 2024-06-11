package face

import (
	"strings"

	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type cardBasicModel interface {
	GetChoices() []string
	GetCursor() int
	GetSelected() map[int]struct{}
	GetData() any
}

func (m cardDeleteModel) GetChoices() []string {
	return m.choices
}

func (m cardDeleteModel) GetCursor() int {
	return m.cursor
}

func (m cardDeleteModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m cardDeleteModel) GetData() any {
	return m.datas
}

func (m cardMenuModel) GetChoices() []string {
	return m.choices
}

func (m cardMenuModel) GetCursor() int {
	return m.cursor
}

func (m cardMenuModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m cardMenuModel) GetData() any {
	return nil
}

func (m cardViewModel) GetChoices() []string {
	return m.choices
}

func (m cardViewModel) GetCursor() int {
	return m.cursor
}

func (m cardViewModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m cardViewModel) GetData() any {
	return m.datas
}

func (m fileDeleteModel) GetChoices() []string {
	return m.choices
}

func (m fileDeleteModel) GetCursor() int {
	return m.cursor
}

func (m fileDeleteModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m fileDeleteModel) GetData() any {
	return m.datas
}

func (m fileMenuModel) GetChoices() []string {
	return m.choices
}

func (m fileMenuModel) GetCursor() int {
	return m.cursor
}

func (m fileMenuModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m fileMenuModel) GetData() any {
	return nil
}

func (m fileViewModel) GetChoices() []string {
	return m.choices
}

func (m fileViewModel) GetCursor() int {
	return m.cursor
}

func (m fileViewModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m fileViewModel) GetData() any {
	return m.datas
}

func (m textDeleteModel) GetChoices() []string {
	return m.choices
}

func (m textDeleteModel) GetCursor() int {
	return m.cursor
}

func (m textDeleteModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m textDeleteModel) GetData() any {
	return m.datas
}

func (m textMenuModel) GetChoices() []string {
	return m.choices
}

func (m textMenuModel) GetCursor() int {
	return m.cursor
}

func (m textMenuModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m textMenuModel) GetData() any {
	return nil
}

func (m textViewModel) GetChoices() []string {
	return m.choices
}

func (m textViewModel) GetCursor() int {
	return m.cursor
}

func (m textViewModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m textViewModel) GetData() any {
	return m.datas
}

func (m lpwDeleteModel) GetChoices() []string {
	return m.choices
}

func (m lpwDeleteModel) GetCursor() int {
	return m.cursor
}

func (m lpwDeleteModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m lpwDeleteModel) GetData() any {
	return m.datas
}

func (m lpwMenuModel) GetChoices() []string {
	return m.choices
}

func (m lpwMenuModel) GetCursor() int {
	return m.cursor
}

func (m lpwMenuModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m lpwMenuModel) GetData() any {
	return nil
}

func (m lpwViewModel) GetChoices() []string {
	return m.choices
}

func (m lpwViewModel) GetCursor() int {
	return m.cursor
}

func (m lpwViewModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m lpwViewModel) GetData() any {
	return m.datas
}

func (m authModel) GetChoices() []string {
	return m.choices
}

func (m authModel) GetCursor() int {
	return m.cursor
}

func (m authModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m authModel) GetData() any {
	return nil
}

func (m mainMenuModel) GetChoices() []string {
	return m.choices
}

func (m mainMenuModel) GetCursor() int {
	return m.cursor
}

func (m mainMenuModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m mainMenuModel) GetData() any {
	return nil
}

func buildView(m cardBasicModel, header string) string {
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

func buildUpdate(h *string, msg tea.Msg, m cardBasicModel, r tea.Model, f func(selected map[int]struct{}, cursor int, datas any) (tea.Model, tea.Cmd)) (tea.Model, tea.Cmd) {

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
		return r, tea.Quit

	case "up", "k":
		if cursor > 0 {
			cursor--
		}

	case "down", "j":
		if cursor < len(choices)-1 {
			cursor++
		}

	case "ctrl+b":
		return NewCardMenuModel(), tea.ClearScreen

	case "enter", " ":
		_, ok := selected[cursor]

		if ok {
			delete(selected, cursor)
		} else {
			return f(selected, cursor, datas)
		}
	}
	return r, nil
}

func updateDeleteCardModel(selected map[int]struct{}, cursor int, datas any) (tea.Model, tea.Cmd) {
	selected[cursor] = struct{}{}
	cards := datas.([]storage.BankCardResponse)
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
