package face

import (
	"testing"

	mock_storage "github.com/Azcarot/PasswordStorage/internal/mock"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewUpdateTextModel(t *testing.T) {

	model := NewUpdateTextModel()

	assert.Len(t, model.inputs, 2, "Expected inputs length to be 2")

	txtInput := model.inputs[txt]
	assert.Equal(t, "Some text **************", txtInput.Placeholder, "Expected placeholder to match")
	assert.Equal(t, 100, txtInput.CharLimit, "Expected char limit to be 100")
	assert.Equal(t, 30, txtInput.Width, "Expected width to be 30")
	assert.Equal(t, "", txtInput.Prompt, "Expected prompt to be empty")
	assert.Equal(t, "", txtInput.Value(), "Expected value to match ")

	tcmtInput := model.inputs[tcmt]
	assert.Equal(t, "Some comment", tcmtInput.Placeholder, "Expected placeholder to match")
	assert.Equal(t, 100, tcmtInput.CharLimit, "Expected char limit to be 100")
	assert.Equal(t, 40, tcmtInput.Width, "Expected width to be 40")
	assert.Equal(t, "", tcmtInput.Prompt, "Expected prompt to be empty")
	assert.Equal(t, "", tcmtInput.Value(), "Expected value to match selectedText.Comment")

	assert.Equal(t, 0, model.focused, "Expected focused to be 0")
	assert.Nil(t, model.err, "Expected err to be nil")
}

func TestTextUpdateModel_Update(t *testing.T) {
	model := NewUpdateTextModel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_storage.NewMockPgxStorage(ctrl)
	storage.TLiteS = mock
	mock.EXPECT().GetAllRecords(gomock.Any()).Times(2).Return([]storage.TextData{}, nil)
	model.focused = 1 // Focus on last input
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	model.inputs[txt].SetValue("Test text")
	model.inputs[tcmt].SetValue("Test comment")
	selectedText.ID = 1

	newModel, cmd := model.Update(msg)

	assert.IsType(t, textViewModel{}, newModel, "Expected newModel to be of type textViewModel")
	assert.Nil(t, cmd, "Expected cmd to be nil")
	assert.Equal(t, "Something went wrong, please try again", textUpdateHeader, "Expected textUpdateHeader to be updated")

	newModel, cmd = model.Update(msg)
	assert.IsType(t, textViewModel{}, newModel, "Expected newModel to be of type textUpdateModel")
	assert.Nil(t, cmd, "Expected cmd to be nil")
	assert.Equal(t, "Something went wrong, please try again", textUpdateHeader, "Expected textUpdateHeader to be updated with error message")

	// Test case 4: Tab key, should focus next input
	model.focused = 0 // Focus on first input
	msg = tea.KeyMsg{Type: tea.KeyTab}
	newModel, _ = model.Update(msg)

	assert.Equal(t, 1, newModel.(textUpdateModel).focused, "Expected focus to be on the second input")

	// Test case 5: Shift + Tab key, should focus previous input
	msg = tea.KeyMsg{Type: tea.KeyShiftTab}
	newModel, _ = model.Update(msg)
	assert.Equal(t, 1, newModel.(textUpdateModel).focused, "Expected focus to be on the first input")
}
