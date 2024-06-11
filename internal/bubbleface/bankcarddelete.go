// Package face - модуль взаимодействия с клиентом посредством bubbletea
package face

import (
	"context"

	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type cardDeleteModel struct {
	choices  []string
	datas    []storage.BankCardResponse
	cursor   int
	selected map[int]struct{}
}

var cardDeleteHeader string = "Please select card to delete"

// NewCardDeleteModel - основная функция для построения и работы с
// моделью просмотра списка/удаления конкретной карты
func NewCardDeleteModel() cardDeleteModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var datas []storage.BankCardResponse
	data, err := storage.BCLiteS.GetAllRecords(ctx)
	if err != nil {
		return cardDeleteModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}

	var b [16]byte
	copy(b[:], storage.Secret)
	ctx = context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	choices, datas = deCypherBankCard(ctx, data.([]storage.BankCardResponse))

	return cardDeleteModel{

		choices:  choices,
		datas:    datas,
		selected: make(map[int]struct{}),
	}
}

// Init - инициализация модели для просмотра списка/удаления конкретной карты
func (m cardDeleteModel) Init() tea.Cmd {
	return nil
}

// View - настройка структуры отображения списка/удаления конкретной карты
func (m cardDeleteModel) View() string {

	s := buildView(&m, cardDeleteHeader)

	return s
}

// Update - функция c описанием реакций на разные события (нажатие клавиш, сообщения)
func (m cardDeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return buildUpdate(&cardDeleteHeader, msg, &m, NewCardMenuModel(), updateDeleteCardModel)

}
