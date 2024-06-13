package face

import (
	"context"
	"fmt"

	"github.com/Azcarot/PasswordStorage/internal/cypher"
	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

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

func updateTextViewModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	data := datas.([]storage.TextData)
	selected[cursor] = struct{}{}
	selectedText = data[cursor]
	return NewUpdateTextModel(), nil
}

func updateTextDeleteModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	data := datas.([]storage.TextData)
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

func deCypherText(ctx context.Context, cards []storage.TextData) ([]string, []storage.TextData, error) {
	var err error
	var choices []string
	var datas []storage.TextData
	for _, data := range cards {
		err = cypher.DeCypherTextData(ctx, &data)
		if err != nil {
			return choices, datas, err
		}
		str := fmt.Sprintf("Your text: %s \nComment: %s", data.Text, data.Comment)
		choices = append(choices, str)
		datas = append(datas, data)
	}
	return choices, datas, nil
}
