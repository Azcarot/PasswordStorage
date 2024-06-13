package face

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Azcarot/PasswordStorage/internal/cypher"
	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *fileDeleteModel) GetChoices() []string {
	return m.choices
}

func (m *fileDeleteModel) GetCursor() int {
	return m.cursor
}

func (m *fileDeleteModel) SetCursor(c int) {
	m.cursor = c
}

func (m *fileDeleteModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *fileDeleteModel) GetData() any {
	return m.datas
}

func (m *fileMenuModel) GetChoices() []string {
	return m.choices
}

func (m *fileMenuModel) GetCursor() int {
	return m.cursor
}

func (m *fileMenuModel) SetCursor(c int) {
	m.cursor = c
}

func (m *fileMenuModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *fileMenuModel) GetData() any {
	return nil
}

func (m *fileViewModel) GetChoices() []string {
	return m.choices
}

func (m *fileViewModel) GetCursor() int {
	return m.cursor
}

func (m *fileViewModel) SetCursor(c int) {
	m.cursor = c
}

func (m *fileViewModel) GetSelected() map[int]struct{} {
	return m.selected
}

func (m *fileViewModel) GetData() any {
	return m.datas
}

func updateFileDeleteModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	selected[cursor] = struct{}{}
	data := datas.([]storage.FileData)
	selectedFile = data[cursor]
	var fileData storage.FileData

	fileData.FileName = selectedFile.FileName
	fileData.Path = selectedFile.Path
	fileData.Comment = selectedFile.Comment
	fileData.ID = selectedFile.ID
	err := storage.FLiteS.AddData(fileData)
	if err != nil {
		deleteFileHeader = "Corrupted file data, please try again"
		return NewFileDeleteModel(), nil
	}

	pauseTicker <- true

	defer func() { resumeTicker <- true }()

	ok, err := requests.DeleteFileReq(fileData)

	if err != nil {
		deleteFileHeader = "Something went wrong, please try again"
		return NewFileDeleteModel(), nil
	}
	if !ok {
		deleteFileHeader = "Wrond file data, try again"
		return NewFileDeleteModel(), nil
	}
	deleteFileHeader = "File succsesfully deleted!"
	return NewFileDeleteModel(), tea.ClearScreen
}

func updateFileMenuModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	selected[cursor] = struct{}{}
	if choices[cursor] == fileChoices.Add {
		return NewAddFileModel(), nil
	}
	if choices[cursor] == fileChoices.View {
		return NewFileViewModel(), nil
	}
	if choices[cursor] == fileChoices.Delete {
		return NewFileDeleteModel(), nil
	}
	return NewFileMenuModel(), nil
}

func updateFileViewModel(selected map[int]struct{}, cursor int, datas any, choices []string) (tea.Model, tea.Cmd) {
	selected[cursor] = struct{}{}
	data := datas.([]storage.FileData)
	selectedFile = data[cursor]
	return NewUpdateFileModel(), nil
}

func deCypherFile(ctx context.Context, cards []storage.FileData) ([]string, []storage.FileData, error) {
	var err error
	var choices []string
	var datas []storage.FileData
	for _, data := range cards {
		err = cypher.DeCypherFileData(ctx, &data)
		if err != nil {
			return choices, datas, err
		}
		str := fmt.Sprintf("Your file name: %s \nYour file path: %s \nComment: %s", data.FileName, data.Path, data.Comment)
		choices = append(choices, str)
		datas = append(datas, data)
	}
	return choices, datas, nil
}

func downloadFileFromDB(fdata storage.FileData, filePath string) error {
	absPath, err := filepath.Abs(filePath)

	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Dir(absPath), 0777); err != nil {
			return err
		}
	}
	file, err := os.Create(fdata.FileName)
	if err != nil {
		return err
	}
	// Ensure the file is closed when the function exits
	defer file.Close()

	// Write data to the file
	_, err = file.WriteString(fdata.Data)
	if err != nil {
		return err
	}

	return nil
}
