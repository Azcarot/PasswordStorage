package face

import (
	"context"

	"github.com/Azcarot/PasswordStorage/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type cardViewModel struct {
	choices  []string
	datas    []storage.BankCardData
	cursor   int
	selected map[int]struct{}
}

type MsgSendData struct {
	Data string
}

var selectedCard storage.BankCardData
var cardViewHeader string = "Main bank card menu, please chose your action:\n\n"

// NewCardViewModel - основная функция для построения и отбражения
// списка карт для их просмотра/обновления
func NewCardViewModel() cardViewModel {
	ctx := context.WithValue(context.Background(), storage.UserLoginCtxKey, storage.UserLoginPw.Login)
	var choices []string
	var datas []storage.BankCardData
	data, err := storage.BCLiteS.GetAllRecords(ctx)
	if err != nil {
		return cardViewModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}

	choices, datas, err = deCypherBankCard(ctx, data.([]storage.BankCardData))
	if err != nil {
		cardViewHeader = "Bank card view error has occured, please try again"
		return cardViewModel{

			choices: []string{},

			selected: make(map[int]struct{}),
		}
	}

	return cardViewModel{

		choices:  choices,
		datas:    datas,
		selected: make(map[int]struct{}),
	}
}

func (m cardViewModel) Init() tea.Cmd {
	return nil
}

func (m cardViewModel) View() string {

	s := buildView(&m, cardViewHeader)

	return s
}

func (m cardViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return buildUpdate(&cardViewHeader, msg, &m, NewCardMenuModel(), updateViewCardModel)

}
