package infrastructure

import (
	"encoding/json"
	"io"
	"net/http"
	"subscriptions-service/pkg/domain"
)

func handleError(err error, w http.ResponseWriter) {
	if err == nil {
		return
	}
	b, err := json.Marshal(msgResponse{Msg: err.Error()})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	switch err {
	case domain.ErrAmountOfSymbolsIsInvalid,
		domain.ErrBalanceIsNotFound,
		domain.ErrThereAreNotEnoughSymbolsOnTheBalance,
		ErrBodyParsing,
		ErrInvalidRequest:
		w.WriteHeader(http.StatusBadRequest)
	case domain.ErrBalanceIsAlreadyExists:
		w.WriteHeader(http.StatusConflict)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, _ = io.WriteString(w, string(b))
}
