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
	SetCursor(int)
	GetSelected() map[int]struct{}
	GetData() any
}

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

func (m *fileDeleteModel) GetChoices() []string {
	return m.choices
}

func (m *fileDeleteModel) GetCursor() int {
	return m.cursor
}

func (m *fileDeleteModel) SetCursor(c int) {
	m.cursor = c
}

func (m *fileDeleteModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *fileDeleteModel) GetData() any {
	return m.datas
}

func (m *fileMenuModel) GetChoices() []string {
	return m.choices
}

func (m *fileMenuModel) GetCursor() int {
	return m.cursor
}

func (m *fileMenuModel) SetCursor(c int) {
	m.cursor = c
}

func (m *fileMenuModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *fileMenuModel) GetData() any {
	return nil
}

func (m *fileViewModel) GetChoices() []string {
	return m.choices
}

func (m *fileViewModel) GetCursor() int {
	return m.cursor
}

func (m *fileViewModel) SetCursor(c int) {
	m.cursor = c
}

func (m *fileViewModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *fileViewModel) GetData() any {
	return m.datas
}

func (m *textDeleteModel) GetChoices() []string {
	return m.choices
}

func (m *textDeleteModel) GetCursor() int {
	return m.cursor
}

func (m *textDeleteModel) SetCursor(c int) {
	m.cursor = c
}

func (m *textDeleteModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *textDeleteModel) GetData() any {
	return m.datas
}

func (m *textMenuModel) GetChoices() []string {
	return m.choices
}

func (m *textMenuModel) GetCursor() int {
	return m.cursor
}

func (m *textMenuModel) SetCursor(c int) {
	m.cursor = c
}

func (m *textMenuModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *textMenuModel) GetData() any {
	return nil
}

func (m *textViewModel) GetChoices() []string {
	return m.choices
}

func (m *textViewModel) GetCursor() int {
	return m.cursor
}

func (m *textViewModel) SetCursor(c int) {
	m.cursor = c
}

func (m *textViewModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *textViewModel) GetData() any {
	return m.datas
}

func (m *lpwDeleteModel) GetChoices() []string {
	return m.choices
}

func (m *lpwDeleteModel) GetCursor() int {
	return m.cursor
}

func (m *lpwDeleteModel) SetCursor(c int) {
	m.cursor = c
}

func (m *lpwDeleteModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *lpwDeleteModel) GetData() any {
	return m.datas
}

func (m *lpwMenuModel) GetChoices() []string {
	return m.choices
}

func (m *lpwMenuModel) GetCursor() int {
	return m.cursor
}
func (m *lpwMenuModel) SetCursor(c int) {
	m.cursor = c
}

func (m *lpwMenuModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *lpwMenuModel) GetData() any {
	return nil
}

func (m *lpwViewModel) GetChoices() []string {
	return m.choices
}

func (m *lpwViewModel) GetCursor() int {
	return m.cursor
}

func (m *lpwViewModel) SetCursor(c int) {
	m.cursor = c
}

func (m *lpwViewModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *lpwViewModel) GetData() any {
	return m.datas
}

func (m *authModel) GetChoices() []string {
	return m.choices
}

func (m *authModel) GetCursor() int {
	return m.cursor
}

func (m *authModel) SetCursor(c int) {
	m.cursor = c
}

func (m *authModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *authModel) GetData() any {
	return nil
}

func (m *mainMenuModel) GetChoices() []string {
	return m.choices
}

func (m *mainMenuModel) GetCursor() int {
	return m.cursor
}

func (m *mainMenuModel) SetCursor(c int) {
	m.cursor = c
}

func (m *mainMenuModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *mainMenuModel) GetData() any {
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

func buildUpdate(h *string, msg tea.Msg, m cardBasicModel, r tea.Model, f func(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd)) (tea.Model, tea.Cmd) {
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
	data := datas.([]storage.BankCardResponse)
	selected[cursor] = struct{}{}
	selectedCard = data[cursor]
	return NewUpdateCardModel(), nil
}

func updateMainMenu(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	selected[cursor] = struct{}{}
	if choices[cursor] == mainChoices.Card {
		return NewCardMenuModel(), nil
	}
	if choices[cursor] == mainChoices.Text {
		return NewTextMenuModel(), nil
	}
	if choices[cursor] == mainChoices.Login {
		return NewLPWMenuModel(), nil
	}
	if choices[cursor] == mainChoices.File {
		return NewFileMenuModel(), nil
	}
	return NewMainMenuModel(), nil
}

func updateFileDeleteModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	selected[cursor] = struct{}{}
	data := datas.([]storage.FileResponse)
	selectedFile = data[cursor]
	var fileData storage.FileData

	fileData.FileName = selectedFile.FileName
	fileData.Path = selectedFile.Path
	fileData.Comment = selectedFile.Comment
	fileData.ID = selectedFile.ID
	err := storage.FLiteS.AddData(fileData)
	if err != nil {
		deleteFileHeader = "Corrupted file data, please try again"
		return NewFileDeleteModel(), nil
	}

	pauseTicker <- true

	defer func() { resumeTicker <- true }()

	ok, err := requests.DeleteFileReq(fileData)

	if err != nil {
		deleteFileHeader = "Something went wrong, please try again"
		return NewFileDeleteModel(), nil
	}
	if !ok {
		deleteFileHeader = "Wrond file data, try again"
		return NewFileDeleteModel(), nil
	}
	deleteFileHeader = "File succsesfully deleted!"
	return NewFileDeleteModel(), tea.ClearScreen
}

func updateFileMenuModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	selected[cursor] = struct{}{}
	if choices[cursor] == fileChoices.Add {
		return NewAddFileModel(), nil
	}
	if choices[cursor] == fileChoices.View {
		return NewFileViewModel(), nil
	}
	if choices[cursor] == fileChoices.Delete {
		return NewFileDeleteModel(), nil
	}
	return NewFileMenuModel(), nil
}

func updateFileViewModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	selected[cursor] = struct{}{}
	data := datas.([]storage.FileResponse)
	selectedFile = data[cursor]
	return NewUpdateFileModel(), nil
}

func updateLPWViewModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	data := datas.([]storage.LoginResponse)
	selected[cursor] = struct{}{}
	selectedLPW = data[cursor]
	return NewUpdateLPWModel(), nil
}

func updateLPWDeleteModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	selected[cursor] = struct{}{}
	data := datas.([]storage.LoginResponse)
	selectedLPW = data[cursor]
	var lpwData storage.LoginData

	lpwData.Login = selectedLPW.Login
	lpwData.Password = selectedLPW.Password
	lpwData.Comment = selectedLPW.Comment
	lpwData.ID = selectedLPW.ID
	err := storage.LPWLiteS.AddData(lpwData)
	if err != nil {
		lpwDeleteHeader = "Corrupted login/pw data, please try again"
		return NewLPWDeleteModel(), nil
	}

	pauseTicker <- true

	defer func() { resumeTicker <- true }()

	ok, err := requests.DeleteLPWReq(lpwData)

	if err != nil {
		lpwDeleteHeader = "Something went wrong, please try again"
		return NewLPWDeleteModel(), nil
	}
	if !ok {
		lpwDeleteHeader = "Wrond login/pw data, try again"
		return NewLPWDeleteModel(), nil
	}
	lpwDeleteHeader = "login/pw succsesfully deleted!"
	return NewLPWDeleteModel(), tea.ClearScreen
}

func updateLPWMenuModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	selected[cursor] = struct{}{}
	if choices[cursor] == lpwChoices.Add {
		return NewAddLPWModel(), nil
	}
	if choices[cursor] == lpwChoices.View {
		return NewLPWViewModel(), nil
	}
	if choices[cursor] == lpwChoices.Delete {
		return NewLPWDeleteModel(), nil
	}
	return NewLPWMenuModel(), nil
}

func updateTextViewModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	data := datas.([]storage.TextResponse)
	selected[cursor] = struct{}{}
	selectedText = data[cursor]
	return NewUpdateTextModel(), nil
}

func updateTextDeleteModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	data := datas.([]storage.TextResponse)
	selected[cursor] = struct{}{}
	selectedText = data[cursor]
	var textData storage.TextData

	textData.Text = selectedText.Text
	textData.Comment = selectedText.Comment
	textData.ID = selectedText.ID
	err := storage.TLiteS.AddData(textData)
	if err != nil {
		textDeleteHeader = "Corrupted text data, please try again"
		return NewTextDeleteModel(), nil
	}

	pauseTicker <- true

	defer func() { resumeTicker <- true }()

	ok, err := requests.DeleteTextReq(textData)

	if err != nil {
		textDeleteHeader = "Something went wrong, please try again"
		return NewTextDeleteModel(), nil
	}
	if !ok {
		textDeleteHeader = "Wrond text data, try again"
		return NewTextDeleteModel(), nil
	}
	textDeleteHeader = "Text succsesfully deleted!"
	return NewTextDeleteModel(), tea.ClearScreen
}

func updateTextMenuModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	selected[cursor] = struct{}{}
	if choices[cursor] == textChoices.Add {
		return NewAddTextModel(), nil
	}
	if choices[cursor] == textChoices.View {
		return NewTextViewModel(), nil
	}
	if choices[cursor] == textChoices.Delete {
		return NewTextDeleteModel(), nil
	}
	return NewTextMenuModel(), nil
}
