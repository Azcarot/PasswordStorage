package face

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/stretchr/testify/assert"
)

func TestRegmodel_View(t *testing.T) {
	inputs := []textinput.Model{
		textinput.New(),
		textinput.New(),
	}
	inputs[login].SetValue("testlogin")
	inputs[pwd].SetValue("testpassword")

	model := regmodel{
		header:  "Test Header",
		inputs:  inputs,
		focused: 0,
	}

	expectedOutput := fmt.Sprintf(
		` %s,
 press ctrl+b to go to the previous menu:

 %s
 %s

 %s  
 %s

 %s
`,
		"Test Header",
		inputStyle.Width(30).Render("Login"),
		model.inputs[login].View(),
		inputStyle.Width(30).Render("Password"),
		model.inputs[pwd].View(),
		continueStyle.Render("Continue"),
	) + "\n"

	assert.Equal(t, expectedOutput, model.View(), "Expected View output to match")
}

func TestRegmodel_NextInput(t *testing.T) {
	inputs := []textinput.Model{
		textinput.New(),
		textinput.New(),
	}

	model := regmodel{
		header:  "Test Header",
		inputs:  inputs,
		focused: 0,
	}

	model.nextInput()
	assert.Equal(t, 1, model.focused, "Expected focused input to be 1")

	model.nextInput()
	assert.Equal(t, 0, model.focused, "Expected focused input to wrap around to 0")
}

func TestRegmodel_PrevInput(t *testing.T) {
	inputs := []textinput.Model{
		textinput.New(),
		textinput.New(),
	}

	model := regmodel{
		header:  "Test Header",
		inputs:  inputs,
		focused: 0,
	}

	model.prevInput()
	assert.Equal(t, 1, model.focused, "Expected focused input to wrap around to 1")

	model.prevInput()
	assert.Equal(t, 0, model.focused, "Expected focused input to be 0")
}
