package face

import (
	"testing"

	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestCardDeleteModel_GetChoices(t *testing.T) {
	model := &cardDeleteModel{choices: []string{"choice1", "choice2"}}
	assert.Equal(t, []string{"choice1", "choice2"}, model.GetChoices())
}

func TestCardDeleteModel_GetCursor(t *testing.T) {
	model := &cardDeleteModel{cursor: 1}
	assert.Equal(t, 1, model.GetCursor())
}

func TestCardDeleteModel_SetCursor(t *testing.T) {
	model := &cardDeleteModel{}
	model.SetCursor(2)
	assert.Equal(t, 2, model.GetCursor())
}

func TestCardDeleteModel_GetSelected(t *testing.T) {
	model := &cardDeleteModel{selected: map[int]struct{}{1: {}}}
	assert.Equal(t, map[int]struct{}{1: {}}, model.GetSelected())
}

func TestCardDeleteModel_GetData(t *testing.T) {
	expectedData := []storage.BankCardData{}
	model := &cardDeleteModel{datas: expectedData}
	assert.Equal(t, expectedData, model.GetData())
}

func TestCardMenuModel_GetChoices(t *testing.T) {
	model := &cardMenuModel{choices: []string{"choice1", "choice2"}}
	assert.Equal(t, []string{"choice1", "choice2"}, model.GetChoices())
}

func TestCardMenuModel_GetCursor(t *testing.T) {
	model := &cardMenuModel{cursor: 1}
	assert.Equal(t, 1, model.GetCursor())
}

func TestCardMenuModel_SetCursor(t *testing.T) {
	model := &cardMenuModel{}
	model.SetCursor(2)
	assert.Equal(t, 2, model.GetCursor())
}

func TestCardMenuModel_GetSelected(t *testing.T) {
	model := &cardMenuModel{selected: map[int]struct{}{1: {}}}
	assert.Equal(t, map[int]struct{}{1: {}}, model.GetSelected())
}

func TestCardMenuModel_GetData(t *testing.T) {
	model := &cardMenuModel{}
	assert.Nil(t, model.GetData())
}

func TestCardViewModel_GetChoices(t *testing.T) {
	model := &cardViewModel{choices: []string{"choice1", "choice2"}}
	assert.Equal(t, []string{"choice1", "choice2"}, model.GetChoices())
}

func TestCardViewModel_GetCursor(t *testing.T) {
	model := &cardViewModel{cursor: 1}
	assert.Equal(t, 1, model.GetCursor())
}

func TestCardViewModel_SetCursor(t *testing.T) {
	model := &cardViewModel{}
	model.SetCursor(2)
	assert.Equal(t, 2, model.GetCursor())
}

func TestCardViewModel_GetSelected(t *testing.T) {
	model := &cardViewModel{selected: map[int]struct{}{1: {}}}
	assert.Equal(t, map[int]struct{}{1: {}}, model.GetSelected())
}

func TestCardViewModel_GetData(t *testing.T) {
	expectedData := []storage.BankCardData{}
	model := &cardViewModel{datas: expectedData}
	assert.Equal(t, expectedData, model.GetData())
}
