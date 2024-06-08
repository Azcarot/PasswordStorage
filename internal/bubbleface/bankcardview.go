package face

import (
	"context"
	"fmt"

	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type cardViewModel struct {
	choices  []string
	cards    []storage.BankCardResponse
	cursor   int
	selected map[int]struct{}
}

type MsgSendData struct {
	Data string
}

var selectedCard storage.BankCardResponse

// CardViewModel - основная функция для построения и отбражения
// списка карт для их просмотра/обновления
func CardViewModel() cardViewModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var cards []storage.BankCardResponse
	data, err := storage.BCLiteS.GetAllRecords(ctx)
	if err != nil {
		return cardViewModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}
	var b [16]byte
	copy(b[:], storage.Secret)
	ctx = context.WithValue(context.Background(), storage.EncryptionCtxKey, b)
	choices, cards = deCypherBankCard(ctx, data.([]storage.BankCardResponse))

	return cardViewModel{

		choices:  choices,
		cards:    cards,
		selected: make(map[int]struct{}),
	}
}

func (m cardViewModel) Init() tea.Cmd {
	return nil
}

func (m cardViewModel) View() string {
	s := "Main bank card menu, please chose your action:\n\n"

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

func (m cardViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				return UpdateCardModel(), nil

			}

		}
	}

	return m, nil
}

func deCypherBankCard(ctx context.Context, cards []storage.BankCardResponse) ([]string, []storage.BankCardResponse) {
	var err error
	var choices []string
	var datas []storage.BankCardResponse
	for _, data := range cards {

		data.CardNumber, err = storage.Dechypher(ctx, data.CardNumber)
		if err != nil {
			continue
		}
		data.Cvc, err = storage.Dechypher(ctx, data.Cvc)
		if err != nil {
			continue
		}
		data.ExpDate, err = storage.Dechypher(ctx, data.ExpDate)
		if err != nil {
			continue
		}
		data.FullName, err = storage.Dechypher(ctx, data.FullName)
		if err != nil {
			continue
		}
		data.Comment, err = storage.Dechypher(ctx, data.Comment)
		if err != nil {
			continue
		}
		str := fmt.Sprintf("CCN: %s \nCVC: %s   ExpDate: %s   FillName: %s\nComment: %s", data.CardNumber, data.Cvc, data.ExpDate, data.FullName, data.Comment)
		choices = append(choices, str)
		datas = append(datas, data)

	}

	return choices, datas
}
