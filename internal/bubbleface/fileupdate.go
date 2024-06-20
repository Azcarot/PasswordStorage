package face

import (
	"fmt"

	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type fileUpdateModel struct {
	inputs  []textinput.Model
	focused int
	err     error
}

var updateFileHeader string = "Please check file data:"

// NewUpdateFileModel - основная функция для работы с
// меню обновления и сохранения на диск файла
func NewUpdateFileModel() fileUpdateModel {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	inputs[finm] = textinput.New()
	inputs[finm].Placeholder = "File name **************"
	inputs[finm].Focus()
	inputs[finm].CharLimit = 100
	inputs[finm].Width = 30
	inputs[finm].Prompt = ""
	inputs[finm].SetValue(selectedFile.FileName)

	inputs[fpth] = textinput.New()
	inputs[fpth].Placeholder = "File path **************"
	inputs[fpth].CharLimit = 100
	inputs[fpth].Width = 100
	inputs[fpth].Prompt = ""
	inputs[fpth].SetValue(selectedFile.Path)

	inputs[fcmt] = textinput.New()
	inputs[fcmt].Placeholder = "Some comment"
	inputs[fcmt].CharLimit = 100
	inputs[fcmt].Width = 40
	inputs[fcmt].Prompt = ""
	inputs[fcmt].SetValue(selectedFile.Comment)

	return fileUpdateModel{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m fileUpdateModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m fileUpdateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				var req storage.FileData
				req.FileName = (m.inputs[finm].Value())
				req.Path = (m.inputs[fpth].Value())
				req.Comment = (m.inputs[fcmt].Value())
				req.ID = selectedFile.ID
				data, err := checkFileSizeAndRead(req.Path)
				if err != nil {
					updateFileHeader = "Failed to get file data, please try again"
					return NewUpdateFileModel(), nil
				}
				req.Data = data
				ok, err := requests.UpdateFileReq(req)

				if err != nil {
					updateFileHeader = "Something went wrong, please try again"
					return NewFileViewModel(), nil
				}
				if !ok {
					updateFileHeader = "Wrond file data, try again"
					return NewFileViewModel(), nil
				}
				updateFileHeader = "File succsesfully updated!"
				return NewFileViewModel(), tea.ClearScreen
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		case tea.KeyCtrlB:
			return NewFileMenuModel(), tea.ClearScreen

		case tea.KeyCtrlS:
			filePath := "/Downloads/"

			err := downloadFileFromDB(selectedFile, filePath)
			if err != nil {
				updateFileHeader = "Error saving file, try again"
				return NewUpdateFileModel(), nil
			}

			updateFileHeader = "File saved!"
			return NewFileViewModel(), nil
		}

		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m fileUpdateModel) View() string {
	return fmt.Sprintf(
		` %s

 %s
 %s

 %s
 %s

 %s
 %s

 %s


 press ctrl+B to go to previous menu
 press ctrl+S to save file to its path
`,
		updateFileHeader,
		inputStyle.Width(30).Render("File name"),
		m.inputs[finm].View(),
		inputStyle.Width(100).Render("File path"),
		m.inputs[fpth].View(),
		inputStyle.Width(30).Render("Comment"),
		m.inputs[fcmt].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"
}

func (m *fileUpdateModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *fileUpdateModel) prevInput() {
	m.focused--

	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
