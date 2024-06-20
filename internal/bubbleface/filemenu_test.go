package face

import (
	"fmt"
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestNewFileMenuModel(t *testing.T) {
	tests := []struct {
		name string
		want fileMenuModel
	}{
		{name: "name", want: fileMenuModel{choices: []string{fileChoices.Add, fileChoices.View, fileChoices.Delete}, selected: make(map[int]struct{})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFileMenuModel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileMenuModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FileMenuModel_Init(t *testing.T) {
	tests := []struct {
		name string
		m    fileMenuModel
		want tea.Cmd
	}{
		{name: "name", m: fileMenuModel{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileMenuModel.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAddFileModel(t *testing.T) {
	cm := NewAddFileModel()

	assert.NotNil(t, cm.inputs, "inputs should not be nil")
	assert.Equal(t, 3, len(cm.inputs), "there should be 3 inputs")

	flnInput := cm.inputs[finm]
	assert.Equal(t, "Some file", flnInput.Placeholder)
	assert.Equal(t, 100, flnInput.CharLimit)
	assert.Equal(t, 30, flnInput.Width)
	assert.Empty(t, flnInput.Prompt)
	assert.True(t, flnInput.Focused(), "ccn input should be focused")

	// Test exp input
	expInput := cm.inputs[fpth]
	assert.Equal(t, "C/File/file.txt", expInput.Placeholder)
	assert.Equal(t, 300, expInput.CharLimit)
	assert.Equal(t, 300, expInput.Width)
	assert.Equal(t, "> ", expInput.Prompt)

	// Test bcmt input
	bcmtInput := cm.inputs[fcmt]
	assert.Equal(t, "Some comment", bcmtInput.Placeholder)
	assert.Equal(t, 100, bcmtInput.CharLimit)
	assert.Equal(t, 40, bcmtInput.Width)
	assert.Empty(t, bcmtInput.Prompt)
}

func TestFileModelUpdate(t *testing.T) {
	m := NewAddFileModel()

	m.inputs[finm].SetValue("Some file")
	m.inputs[fpth].SetValue("12/24")
	m.inputs[fcmt].SetValue("Test comment")

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := m.Update(msg)

	newModel, ok := updatedModel.(fileModel)
	assert.True(t, ok, "updated model should be of type lpwModel")
	assert.Equal(t, "Some file", newModel.inputs[finm].Value(), "ccn input should be empty after submission")
	assert.Equal(t, "12/24", newModel.inputs[fpth].Value(), "exp input should be empty after submission")
	assert.Equal(t, "Test comment", newModel.inputs[fcmt].Value(), "bcmt input should be empty after submission")

}

func TestFileModelView(t *testing.T) {
	m := NewAddFileModel()

	// Simulate entering data into the inputs
	m.inputs[finm].SetValue("File name")
	m.inputs[fpth].SetValue("File path")
	m.inputs[fcmt].SetValue("Comment")

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
		addFileHeader,
		inputStyle.Width(30).Render("File name"),
		m.inputs[finm].View(),
		inputStyle.Width(100).Render("File path"),
		m.inputs[fpth].View(),
		inputStyle.Width(30).Render("Comment"),
		m.inputs[fcmt].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"

	// Call the View method
	actualView := m.View()

	// Assert that the actual view matches the expected view
	assert.Equal(t, expectedView, actualView, "the view should match the expected format")
}

func Test_FileModel_prevInput(t *testing.T) {
	tests := []struct {
		name string
		m    *fileModel
	}{
		{name: "name", m: &fileModel{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.prevInput()
		})
	}
}
