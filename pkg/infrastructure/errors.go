package infrastructure

import (
	"net/http"
	"subscriptions-service/pkg/domain"
)

func handleError(err error, w http.ResponseWriter) {
	switch err {
	case domain.ErrBalanceIsAlreadyExists:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
