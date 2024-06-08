package face

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	finm = iota
	fpth
	fcmt
)

const maxAllowedFileSize int64 = 2e+6

type fileMenuModel struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}
type fileCho struct {
	Add    string
	View   string
	Delete string
}

var fileChoices = fileCho{
	Add:    "Add New File",
	View:   "View/Download File",
	Delete: "Delete File",
}

// FileMenuModel - основная функция для построения и работы с
// основным меню работы с файлами
func FileMenuModel() fileMenuModel {
	return fileMenuModel{

		choices: []string{fileChoices.Add, fileChoices.View, fileChoices.Delete},

		selected: make(map[int]struct{}),
	}
}

func (m fileMenuModel) Init() tea.Cmd {
	return nil
}

func (m fileMenuModel) View() string {
	s := "Main file menu, please chose your action:\n\n"

	for i, choice := range m.choices {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}

func (m fileMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "ctrl+b":
			return MainMenuModel(), nil

		case "enter", " ":
			_, ok := m.selected[m.cursor]

			if ok {
				delete(m.selected, m.cursor)
			} else {

				m.selected[m.cursor] = struct{}{}
				if m.choices[m.cursor] == fileChoices.Add {
					return AddFileModel(), nil
				}
				if m.choices[m.cursor] == fileChoices.View {
					return FileViewModel(), nil
				}
				if m.choices[m.cursor] == fileChoices.Delete {
					return FileDeleteModel(), nil
				}
			}
		}
	}

	return m, nil
}

type fileModel struct {
	inputs  []textinput.Model
	focused int
	err     error
}

// AddFileModel - основная функция для построения и работы с
// меню добавления нового файла
func AddFileModel() fileModel {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	newheader = "Please insert file name/path:"
	inputs[finm] = textinput.New()
	inputs[finm].Placeholder = "Some file"
	inputs[finm].Focus()
	inputs[finm].CharLimit = 100
	inputs[finm].Width = 30
	inputs[finm].Prompt = ""

	inputs[fpth] = textinput.New()
	inputs[fpth].Placeholder = "C/File/file.txt"
	inputs[fpth].CharLimit = 300
	inputs[fpth].Width = 300

	inputs[fcmt] = textinput.New()
	inputs[fcmt].Placeholder = "Some comment"
	inputs[fcmt].CharLimit = 100
	inputs[fcmt].Width = 40
	inputs[fcmt].Prompt = ""

	return fileModel{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m fileModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m fileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				req.Path = strings.TrimSpace(req.Path)
				data, err := checkFileSizeAndRead(req.Path)
				if err != nil {
					newheader = "Failed to get file data, please try again"
					return AddFileModel(), tea.ClearScreen
				}
				req.Data = data
				ok, err := requests.AddFileReq(req)

				if err != nil {
					newheader = "Something went wrong, please try again"
					return AddFileModel(), nil
				}
				if !ok {
					newheader = "Wrond data, please try again"
					return AddFileModel(), nil
				}
				newheader = "File succsesfully saved!"
				return AddFileModel(), tea.ClearScreen
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		case tea.KeyCtrlB:
			return FileMenuModel(), tea.ClearScreen

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

func (m fileModel) View() string {
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
`,
		newheader,
		inputStyle.Width(30).Render("File name"),
		m.inputs[finm].View(),
		inputStyle.Width(100).Render("File path"),
		m.inputs[fpth].View(),
		inputStyle.Width(30).Render("Comment"),
		m.inputs[fcmt].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"
}

func (m *fileModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *fileModel) prevInput() {
	m.focused--

	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func checkFileSizeAndRead(filePath string) (string, error) {

	absPath, err := filepath.Abs(filePath)

	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	if fileInfo.Size() > maxAllowedFileSize {
		return "", nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(data), nil
}
