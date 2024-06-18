package face

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

var (
	mockFileChoices  = []string{"Choice 1", "Choice 2", "Choice 3"}
	mockSelectedFile = map[int]struct{}{0: {}}
)

func TestFileViewModel_Update(t *testing.T) {
	model := &fileViewModel{
		choices:  mockFileChoices,
		cursor:   0,
		selected: mockSelectedFile,
	}

	msg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedModel, _ := model.Update(msg)
	updatedLPWModel, ok := updatedModel.(*fileViewModel)
	assert.True(t, ok, "Expected updated model to be of type *textViewModel")
	assert.Equal(t, 0, updatedLPWModel.cursor, "Expected cursor to be incremented")
}

func TestFileViewModel_Init(t *testing.T) {
	model := &fileViewModel{}
	cmd := model.Init()
	assert.Nil(t, cmd, "Expected Init command to be nil")
}

func TestFileViewModel_View(t *testing.T) {
	model := &fileViewModel{
		choices:  mockFileChoices,
		cursor:   0,
		selected: mockSelectedFile,
	}

	expectedOutput := buildView(model, fileViewHeader)
	output := model.View()

	assert.Equal(t, expectedOutput, output, "Expected View output to match")
}
