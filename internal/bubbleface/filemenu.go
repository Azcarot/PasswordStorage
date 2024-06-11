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

var menuFileHeader string = "Main file menu, please chose your action:\n\n"
var addFileHeader string = "Please insert file name/path:"

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

// NewFileMenuModel - основная функция для построения и работы с
// основным меню работы с файлами
func NewFileMenuModel() fileMenuModel {
	return fileMenuModel{

		choices: []string{fileChoices.Add, fileChoices.View, fileChoices.Delete},

		selected: make(map[int]struct{}),
	}
}

func (m fileMenuModel) Init() tea.Cmd {
	return nil
}

func (m fileMenuModel) View() string {

	s := buildView(m, menuFileHeader)

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
			return NewMainMenuModel(), nil

		case "enter", " ":
			_, ok := m.selected[m.cursor]

			if ok {
				delete(m.selected, m.cursor)
			} else {

				m.selected[m.cursor] = struct{}{}
				if m.choices[m.cursor] == fileChoices.Add {
					return NewAddFileModel(), nil
				}
				if m.choices[m.cursor] == fileChoices.View {
					return NewFileViewModel(), nil
				}
				if m.choices[m.cursor] == fileChoices.Delete {
					return NewFileDeleteModel(), nil
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

// NewAddFileModel - основная функция для построения и работы с
// меню добавления нового файла
func NewAddFileModel() fileModel {
	var inputs []textinput.Model = make([]textinput.Model, 3)
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
					addFileHeader = "Failed to get file data, please try again"
					return NewAddFileModel(), tea.ClearScreen
				}
				req.Data = data
				ok, err := requests.AddFileReq(req)

				if err != nil {
					addFileHeader = "Something went wrong, please try again"
					return NewAddFileModel(), nil
				}
				if !ok {
					addFileHeader = "Wrond data, please try again"
					return NewAddFileModel(), nil
				}
				addFileHeader = "File succsesfully saved!"
				return NewAddFileModel(), tea.ClearScreen
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
		addFileHeader,
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
