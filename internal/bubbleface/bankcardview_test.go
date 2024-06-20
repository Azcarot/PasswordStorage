package face

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

var (
	mockCardChoices  = []string{"Choice 1", "Choice 2", "Choice 3"}
	mockSelectedCard = map[int]struct{}{0: {}}
)

func TestCardViewModel_Update(t *testing.T) {
	model := &cardViewModel{
		choices:  mockCardChoices,
		cursor:   0,
		selected: mockSelectedCard,
	}

	msg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedModel, _ := model.Update(msg)
	updatedLPWModel, ok := updatedModel.(*cardViewModel)
	assert.True(t, ok, "Expected updated model to be of type *textViewModel")
	assert.Equal(t, 0, updatedLPWModel.cursor, "Expected cursor to be incremented")
}

func TestCardViewModel_Init(t *testing.T) {
	model := &cardViewModel{}
	cmd := model.Init()
	assert.Nil(t, cmd, "Expected Init command to be nil")
}

func TestCardViewModel_View(t *testing.T) {
	model := &cardViewModel{
		choices:  mockCardChoices,
		cursor:   0,
		selected: mockSelectedCard,
	}

	expectedOutput := buildView(model, cardViewHeader)
	output := model.View()

	assert.Equal(t, expectedOutput, output, "Expected View output to match")
}
