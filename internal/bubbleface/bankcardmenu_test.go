package face

import (
	"fmt"
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestNewCardMenuModel(t *testing.T) {
	tests := []struct {
		name string
		want cardMenuModel
	}{
		{name: "name", want: cardMenuModel{choices: []string{bankChoices.Add, bankChoices.View, bankChoices.Delete}, selected: make(map[int]struct{})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCardMenuModel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCardMenuModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ccnValidator(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "name", args: args{s: "1111 2222"}, wantErr: false},
		{name: "name2", args: args{s: "1111 22asd"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ccnValidator(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ccnValidator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_expValidator(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "name", args: args{s: "12/21"}, wantErr: false},
		{name: "name2", args: args{s: "1111 22asd"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := expValidator(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("expValidator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cvvValidator(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "name", args: args{s: "111"}, wantErr: false},
		{name: "name2", args: args{s: "1aasd"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cvvValidator(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("cvvValidator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cardMenuModel_Init(t *testing.T) {
	tests := []struct {
		name string
		m    cardMenuModel
		want tea.Cmd
	}{
		{name: "name", m: cardMenuModel{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cardMenuModel.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAddCardModel(t *testing.T) {
	cm := NewAddCardModel()

	assert.NotNil(t, cm.inputs, "inputs should not be nil")
	assert.Equal(t, 5, len(cm.inputs), "there should be 5 inputs")

	// Test ccn input
	ccnInput := cm.inputs[ccn]
	assert.Equal(t, "4505 **** **** 1234", ccnInput.Placeholder)
	assert.Equal(t, 20, ccnInput.CharLimit)
	assert.Equal(t, 30, ccnInput.Width)
	assert.Empty(t, ccnInput.Prompt)
	assert.NotNil(t, ccnInput.Validate)
	assert.True(t, ccnInput.Focused(), "ccn input should be focused")

	// Test exp input
	expInput := cm.inputs[exp]
	assert.Equal(t, "MM/YY ", expInput.Placeholder)
	assert.Equal(t, 5, expInput.CharLimit)
	assert.Equal(t, 5, expInput.Width)
	assert.Empty(t, expInput.Prompt)
	assert.NotNil(t, expInput.Validate)

	// Test cvv input
	cvvInput := cm.inputs[cvv]
	assert.Equal(t, "XXX", cvvInput.Placeholder)
	assert.Equal(t, 3, cvvInput.CharLimit)
	assert.Equal(t, 5, cvvInput.Width)
	assert.Empty(t, cvvInput.Prompt)
	assert.NotNil(t, cvvInput.Validate)

	// Test fnm input
	fnmInput := cm.inputs[fnm]
	assert.Equal(t, "Ivan ivanov", fnmInput.Placeholder)
	assert.Equal(t, 40, fnmInput.CharLimit)
	assert.Equal(t, 40, fnmInput.Width)
	assert.Empty(t, fnmInput.Prompt)
	assert.Nil(t, fnmInput.Validate)

	// Test bcmt input
	bcmtInput := cm.inputs[bcmt]
	assert.Equal(t, "Some text", bcmtInput.Placeholder)
	assert.Equal(t, 100, bcmtInput.CharLimit)
	assert.Equal(t, 40, bcmtInput.Width)
	assert.Empty(t, bcmtInput.Prompt)
	assert.Nil(t, bcmtInput.Validate)
}

func TestCardModelUpdate(t *testing.T) {
	m := NewAddCardModel()

	m.inputs[ccn].SetValue("4505 1234 5678 9012")
	m.inputs[exp].SetValue("12/24")
	m.inputs[cvv].SetValue("123")
	m.inputs[fnm].SetValue("Ivan Ivanov")
	m.inputs[bcmt].SetValue("Test comment")

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := m.Update(msg)

	newModel, ok := updatedModel.(cardModel)
	assert.True(t, ok, "updated model should be of type cardModel")
	assert.Equal(t, "4505 1234 5678 9012", newModel.inputs[ccn].Value(), "ccn input should be empty after submission")
	assert.Equal(t, "12/24", newModel.inputs[exp].Value(), "exp input should be empty after submission")
	assert.Equal(t, "123", newModel.inputs[cvv].Value(), "cvv input should be empty after submission")
	assert.Equal(t, "Ivan Ivanov", newModel.inputs[fnm].Value(), "fnm input should be empty after submission")
	assert.Equal(t, "Test comment", newModel.inputs[bcmt].Value(), "bcmt input should be empty after submission")

}

func TestCardModelView(t *testing.T) {
	m := NewAddCardModel()

	// Simulate entering data into the inputs
	m.inputs[ccn].SetValue("4505 1234 5678 9012")
	m.inputs[exp].SetValue("12/24")
	m.inputs[cvv].SetValue("123")
	m.inputs[fnm].SetValue("Ivan Ivanov")
	m.inputs[bcmt].SetValue("Test comment")

	expectedView := fmt.Sprintf(
		` %s

 %s
 %s

 %s  %s  %s
 %s  %s  %s

 %s  %s

 %s


 press ctrl+B to go to previous menu
`,
		cardAddHeader,
		inputStyle.Width(30).Render("Card Number"),
		m.inputs[ccn].View(),
		inputStyle.Width(6).Render("EXP"),
		inputStyle.Width(6).Render("CVV"),
		inputStyle.Width(30).Render("Full name"),
		m.inputs[exp].View(),
		m.inputs[cvv].View(),
		m.inputs[fnm].View(),
		inputStyle.Width(30).Render("Comment"),
		m.inputs[bcmt].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"

	// Call the View method
	actualView := m.View()

	// Assert that the actual view matches the expected view
	assert.Equal(t, expectedView, actualView, "the view should match the expected format")
}

func Test_cardModel_prevInput(t *testing.T) {
	tests := []struct {
		name string
		m    *cardModel
	}{
		{name: "name", m: &cardModel{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.prevInput()
		})
	}
}