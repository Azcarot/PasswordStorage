package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/jackc/pgx/v5"
)

type GetReqData struct {
	ID   int    `json:"id"`
	User string `json:"user"`
}

type SerchReqData struct {
	User string `json:"user"`
	Str  string `json:"str"`
}

func GetBankCard(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	_, ok := req.Context().Value(storage.UserLoginCtxKey).(string)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var reqData GetReqData

	data, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	if err = json.Unmarshal(data, &reqData); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	cardData, err := storage.PST.GetRecord(ctx, reqData.ID)

	if errors.Is(err, pgx.ErrNoRows) {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(cardData)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(result)

}
