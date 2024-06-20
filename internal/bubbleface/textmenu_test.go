package face

import (
	"fmt"
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestNewTextMenuModel(t *testing.T) {
	tests := []struct {
		name string
		want textMenuModel
	}{
		{name: "name", want: textMenuModel{choices: []string{textChoices.Add, textChoices.View, textChoices.Delete}, selected: make(map[int]struct{})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTextMenuModel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTextMenuModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_TextMenuModel_Init(t *testing.T) {
	tests := []struct {
		name string
		m    textMenuModel
		want tea.Cmd
	}{
		{name: "name", m: textMenuModel{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TextMenuModel.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAddTextModel(t *testing.T) {
	cm := NewAddTextModel()

	assert.NotNil(t, cm.inputs, "inputs should not be nil")
	assert.Equal(t, 2, len(cm.inputs), "there should be 2 inputs")

	txtInput := cm.inputs[txt]
	assert.Equal(t, "Some text **************", txtInput.Placeholder)
	assert.Equal(t, 100, txtInput.CharLimit)
	assert.Equal(t, 30, txtInput.Width)
	assert.Empty(t, txtInput.Prompt)
	assert.True(t, txtInput.Focused(), "ccn input should be focused")

	bcmtInput := cm.inputs[tcmt]
	assert.Equal(t, "Some comment", bcmtInput.Placeholder)
	assert.Equal(t, 100, bcmtInput.CharLimit)
	assert.Equal(t, 40, bcmtInput.Width)
	assert.Empty(t, bcmtInput.Prompt)
}

func TestTextModelUpdate(t *testing.T) {
	m := NewAddTextModel()

	m.inputs[txt].SetValue("Some text **************")
	m.inputs[tcmt].SetValue("Test comment")

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := m.Update(msg)

	newModel, ok := updatedModel.(textModel)
	assert.True(t, ok, "updated model should be of type textModel")
	assert.Equal(t, "Some text **************", newModel.inputs[txt].Value(), "ccn input should be empty after submission")
	assert.Equal(t, "Test comment", newModel.inputs[tcmt].Value(), "bcmt input should be empty after submission")

}

func TestTextModelView(t *testing.T) {
	m := NewAddTextModel()

	// Simulate entering data into the inputs
	m.inputs[txt].SetValue("Your Text")
	m.inputs[tcmt].SetValue("Comment")

	expectedView := fmt.Sprintf(
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

	// Call the View method
	actualView := m.View()

	// Assert that the actual view matches the expected view
	assert.Equal(t, expectedView, actualView, "the view should match the expected format")
}

func Test_TextModel_prevInput(t *testing.T) {
	tests := []struct {
		name string
		m    *textModel
	}{
		{name: "name", m: &textModel{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.prevInput()
		})
	}
}
