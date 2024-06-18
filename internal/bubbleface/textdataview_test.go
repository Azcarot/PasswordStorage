package face

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

var (
	mockTextChoices  = []string{"Choice 1", "Choice 2", "Choice 3"}
	mockSelectedText = map[int]struct{}{0: {}}
)

func TestTextViewModel_Update(t *testing.T) {
	model := &textViewModel{
		choices:  mockTextChoices,
		cursor:   0,
		selected: mockSelectedText,
	}

	msg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedModel, _ := model.Update(msg)
	updatedTextModel, ok := updatedModel.(*textViewModel)
	assert.True(t, ok, "Expected updated model to be of type *textViewModel")
	assert.Equal(t, 0, updatedTextModel.cursor, "Expected cursor to be incremented")
}

func TestTextViewModel_Init(t *testing.T) {
	model := &textViewModel{}
	cmd := model.Init()
	assert.Nil(t, cmd, "Expected Init command to be nil")
}

func TestTextViewModel_View(t *testing.T) {
	model := &textViewModel{
		choices:  mockTextChoices,
		cursor:   0,
		selected: mockSelectedText,
	}

	expectedOutput := buildView(model, textViewHeader)
	output := model.View()

	assert.Equal(t, expectedOutput, output, "Expected View output to match")
}
