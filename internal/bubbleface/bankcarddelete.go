// Package face - модуль взаимодействия с клиентом посредством bubbletea
package face

import (
	"context"
	"fmt"

	"github.com/Azcarot/PasswordStorage/internal/requests"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type cardDeleteModel struct {
	choices  []string
	cards    []storage.BankCardResponse
	cursor   int
	selected map[int]struct{}
}

// CardDeleteModel - основная функция для построения и работы с
// моделью просмотра списка/удаления конкретной карты
func CardDeleteModel() cardDeleteModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var cards []storage.BankCardResponse
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
	choices, cards = deCypherBankCard(ctx, data.([]storage.BankCardResponse))

	return cardDeleteModel{

		choices:  choices,
		cards:    cards,
		selected: make(map[int]struct{}),
	}
}

// Init - инициализация модели для просмотра списка/удаления конкретной карты
func (m cardDeleteModel) Init() tea.Cmd {
	return nil
}

// View - настройка структуры отображения списка/удаления конкретной карты
func (m cardDeleteModel) View() string {
	newheader = "Please select card to delete"
	s := newheader + "\n\n"

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

// Update - функция c описанием реакций на разные события (нажатие клавиш, сообщения)
func (m cardDeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return CardMenuModel(), tea.ClearScreen

		case "enter", " ":
			_, ok := m.selected[m.cursor]

			if ok {
				delete(m.selected, m.cursor)
			} else {

				m.selected[m.cursor] = struct{}{}
				selectedCard = m.cards[m.cursor]
				var bankData storage.BankCardData
				bankData.CardNumber = selectedCard.CardNumber
				bankData.Cvc = selectedCard.Cvc
				bankData.Comment = selectedCard.Comment
				bankData.ID = selectedCard.ID
				bankData.ExpDate = selectedCard.ExpDate
				bankData.FullName = selectedCard.FullName
				err := storage.BCLiteS.AddData(bankData)
				if err != nil {
					return CardDeleteModel(), nil
				}

				pauseTicker <- true

				defer func() { resumeTicker <- true }()

				ok, err := requests.DeleteCardReq(bankData)

				if err != nil {
					newheader = "Something went wrong, please try again"
					return CardDeleteModel(), nil
				}
				if !ok {
					newheader = "Wrond card data, try again"
					return CardDeleteModel(), nil
				}
				newheader = "Card succsesfully deleted!"
				return CardDeleteModel(), tea.ClearScreen

			}

		}
	}

	return m, nil
}
