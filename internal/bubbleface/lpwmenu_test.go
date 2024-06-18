package face

import (
	"fmt"
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestNewLPWMenuModel(t *testing.T) {
	tests := []struct {
		name string
		want lpwMenuModel
	}{
		{name: "name", want: lpwMenuModel{choices: []string{lpwChoices.Add, lpwChoices.View, lpwChoices.Delete}, selected: make(map[int]struct{})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLPWMenuModel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewlpwMenuModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lpwMenuModel_Init(t *testing.T) {
	tests := []struct {
		name string
		m    lpwMenuModel
		want tea.Cmd
	}{
		{name: "name", m: lpwMenuModel{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lpwMenuModel.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAddLPWModel(t *testing.T) {
	cm := NewAddLPWModel()

	assert.NotNil(t, cm.inputs, "inputs should not be nil")
	assert.Equal(t, 3, len(cm.inputs), "there should be 3 inputs")

	lgnInput := cm.inputs[lgn]
	assert.Equal(t, "Login **************", lgnInput.Placeholder)
	assert.Equal(t, 100, lgnInput.CharLimit)
	assert.Equal(t, 30, lgnInput.Width)
	assert.Empty(t, lgnInput.Prompt)
	assert.True(t, lgnInput.Focused(), "ccn input should be focused")

	// Test exp input
	expInput := cm.inputs[psw]
	assert.Equal(t, "Password **************", expInput.Placeholder)
	assert.Equal(t, 100, expInput.CharLimit)
	assert.Equal(t, 30, expInput.Width)
	assert.Empty(t, expInput.Prompt)

	// Test bcmt input
	bcmtInput := cm.inputs[lcmt]
	assert.Equal(t, "Some comment", bcmtInput.Placeholder)
	assert.Equal(t, 100, bcmtInput.CharLimit)
	assert.Equal(t, 40, bcmtInput.Width)
	assert.Empty(t, bcmtInput.Prompt)
}

func TestLPWModelUpdate(t *testing.T) {
	m := NewAddLPWModel()

	m.inputs[lgn].SetValue("4505123456789012")
	m.inputs[psw].SetValue("12/24")
	m.inputs[lcmt].SetValue("Test comment")

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := m.Update(msg)

	newModel, ok := updatedModel.(lpwModel)
	assert.True(t, ok, "updated model should be of type lpwModel")
	assert.Equal(t, "4505123456789012", newModel.inputs[lgn].Value(), "ccn input should be empty after submission")
	assert.Equal(t, "12/24", newModel.inputs[psw].Value(), "exp input should be empty after submission")
	assert.Equal(t, "Test comment", newModel.inputs[lcmt].Value(), "bcmt input should be empty after submission")

}

func TestLPWModelView(t *testing.T) {
	m := NewAddLPWModel()

	// Simulate entering data into the inputs
	m.inputs[lgn].SetValue("4505123456789012")
	m.inputs[psw].SetValue("12/24")
	m.inputs[lcmt].SetValue("Test comment")

	expectedView := fmt.Sprintf(
		` %s

 %s
 %s

 %s
 %s

 %s
 %s

 %s


 press ctrl+B to go to previous menu
`,
		lpwAddHeader,
		inputStyle.Width(30).Render("Your login"),
		m.inputs[lgn].View(),
		inputStyle.Width(30).Render("Your password"),
		m.inputs[psw].View(),
		inputStyle.Width(30).Render("Comment"),
		m.inputs[lcmt].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"

	// Call the View method
	actualView := m.View()

	// Assert that the actual view matches the expected view
	assert.Equal(t, expectedView, actualView, "the view should match the expected format")
}

func Test_lpwModel_prevInput(t *testing.T) {
	tests := []struct {
		name string
		m    *lpwModel
	}{
		{name: "name", m: &lpwModel{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.prevInput()
		})
	}
}
