package face

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

var (
	mockLPWChoices  = []string{"Choice 1", "Choice 2", "Choice 3"}
	mockSelectedLPW = map[int]struct{}{0: {}}
)

func TestLPWViewModel_Update(t *testing.T) {
	model := &lpwViewModel{
		choices:  mockLPWChoices,
		cursor:   0,
		selected: mockSelectedLPW,
	}

	msg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedModel, _ := model.Update(msg)
	updatedLPWModel, ok := updatedModel.(*lpwViewModel)
	assert.True(t, ok, "Expected updated model to be of type *textViewModel")
	assert.Equal(t, 0, updatedLPWModel.cursor, "Expected cursor to be incremented")
}

func TestLPWViewModel_Init(t *testing.T) {
	model := &lpwViewModel{}
	cmd := model.Init()
	assert.Nil(t, cmd, "Expected Init command to be nil")
}

func TestLPWViewModel_View(t *testing.T) {
	model := &lpwViewModel{
		choices:  mockLPWChoices,
		cursor:   0,
		selected: mockSelectedLPW,
	}

	expectedOutput := buildView(model, lpwViewHeader)
	output := model.View()

	assert.Equal(t, expectedOutput, output, "Expected View output to match")
}
