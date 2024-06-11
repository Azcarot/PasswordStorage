// Package handlers - все обработчики запросов на сервере

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_storage "github.com/Azcarot/PasswordStorage/internal/mock"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddNewCard(t *testing.T) {
	// Create test data
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_storage.NewMockPgxStorage(ctrl)
	storage.BCST = mock
	bankCardData := storage.BankCardData{
		CardNumber: "1234567890123456",
		Cvc:        "123",
		ExpDate:    "12/24",
		FullName:   "John Doe",
		Comment:    "Test Card",
	}

	dataLogin := "testUser"

	body, err := json.Marshal(bankCardData)
	assert.NoError(t, err)
	mock.EXPECT().AddData(gomock.Any()).Times(1)
	mock.EXPECT().CreateNewRecord(gomock.Any()).Times(1)

	req := httptest.NewRequest(http.MethodPost, storage.ServURL+"/card/add", bytes.NewBuffer(body))
	dataLogin = "Login"

	ctx := context.WithValue(req.Context(), storage.UserLoginCtxKey, dataLogin)
	req = req.WithContext(ctx)

	res := httptest.NewRecorder()

	AddNewCard(res, req)

	assert.Equal(t, http.StatusAccepted, res.Code)

}
