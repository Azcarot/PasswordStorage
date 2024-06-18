package face

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthModel_GetChoices(t *testing.T) {
	model := &authModel{choices: []string{"choice1", "choice2"}}
	assert.Equal(t, []string{"choice1", "choice2"}, model.GetChoices())
}

func TestAuthModel_GetCursor(t *testing.T) {
	model := &authModel{cursor: 1}
	assert.Equal(t, 1, model.GetCursor())
}

func TestAuthModel_SetCursor(t *testing.T) {
	model := &authModel{}
	model.SetCursor(2)
	assert.Equal(t, 2, model.GetCursor())
}

func TestAuthModel_GetSelected(t *testing.T) {
	model := &authModel{selected: map[int]struct{}{1: {}}}
	assert.Equal(t, map[int]struct{}{1: {}}, model.GetSelected())
}

func TestAuthModel_GetData(t *testing.T) {
	model := &authModel{}
	assert.Nil(t, model.GetData())
}

func TestMainMenuModel_GetChoices(t *testing.T) {
	model := &mainMenuModel{choices: []string{"choice1", "choice2"}}
	assert.Equal(t, []string{"choice1", "choice2"}, model.GetChoices())
}

func TestMainMenuModel_GetCursor(t *testing.T) {
	model := &mainMenuModel{cursor: 1}
	assert.Equal(t, 1, model.GetCursor())
}

func TestMainMenuModel_SetCursor(t *testing.T) {
	model := &mainMenuModel{}
	model.SetCursor(2)
	assert.Equal(t, 2, model.GetCursor())
}

func TestMainMenuModel_GetSelected(t *testing.T) {
	model := &mainMenuModel{selected: map[int]struct{}{1: {}}}
	assert.Equal(t, map[int]struct{}{1: {}}, model.GetSelected())
}

func TestMainMenuModel_GetData(t *testing.T) {
	model := &mainMenuModel{}
	assert.Nil(t, model.GetData())
}

func TestUpdateMainMenu(t *testing.T) {
	selected := make(map[int]struct{})
	cursor := 0
	datas := any(nil)
	choices := []string{mainChoices.Card, mainChoices.Text, mainChoices.Login, mainChoices.File}

	model, cmd := updateMainMenu(selected, cursor, datas, choices)
	assert.IsType(t, cardMenuModel{}, model)
	assert.Nil(t, cmd)

	cursor = 1
	model, cmd = updateMainMenu(selected, cursor, datas, choices)
	assert.IsType(t, textMenuModel{}, model)
	assert.Nil(t, cmd)

	cursor = 2
	model, cmd = updateMainMenu(selected, cursor, datas, choices)
	assert.IsType(t, lpwMenuModel{}, model)
	assert.Nil(t, cmd)

	cursor = 3
	model, cmd = updateMainMenu(selected, cursor, datas, choices)
	assert.IsType(t, fileMenuModel{}, model)
	assert.Nil(t, cmd)

}
