package infrastructure

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
)

type CreateBalanceDTO struct {
	UserID string `json:"userID"`
}

type Server interface {
	createBalance(w http.ResponseWriter, r *http.Request)
}

var ErrInvalidRequest = fmt.Errorf("balance is already exists")

type server struct {
	dependencyContainer DependencyContainer
}

func (s *server) createBalance(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var balance CreateBalanceDTO
	err = json.Unmarshal(b, &balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	balanceService := s.dependencyContainer.NewBalanceService()
	userID, err := uuid.Parse(balance.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = balanceService.CreateBalance(userID)
	if err != nil {
		handleError(err, w)
	}
}

func NewServer(dependencyContainer DependencyContainer) Server {
	return &server{dependencyContainer: dependencyContainer}
}
