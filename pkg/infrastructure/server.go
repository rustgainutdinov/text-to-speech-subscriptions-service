package infrastructure

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
)

type Server interface {
	createBalance(w http.ResponseWriter, r *http.Request)
	topUpBalance(w http.ResponseWriter, r *http.Request)
	writeOffFromBalance(w http.ResponseWriter, r *http.Request)
	removeBalance(w http.ResponseWriter, r *http.Request)
}

var ErrInvalidRequest = fmt.Errorf("invalid request")
var ErrBodyParsing = fmt.Errorf("body parsing error")

type server struct {
	dependencyContainer DependencyContainer
}

func (s *server) topUpBalance(w http.ResponseWriter, r *http.Request) {
	userID, score, err := getBalanceMovementDataFromReq(r)
	if err != nil {
		handleError(err, w)
		return
	}
	balanceService := s.dependencyContainer.NewBalanceService()
	err = balanceService.TopUpBalance(userID, score)
	if err != nil {
		handleError(err, w)
		return
	}
	setOkResponse(w)
}

func (s *server) writeOffFromBalance(w http.ResponseWriter, r *http.Request) {
	userID, score, err := getBalanceMovementDataFromReq(r)
	if err != nil {
		handleError(err, w)
		return
	}
	balanceService := s.dependencyContainer.NewBalanceService()
	err = balanceService.WriteOffFromBalance(userID, score)
	if err != nil {
		handleError(err, w)
		return
	}
	setOkResponse(w)
}

func (s *server) createBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromReq(r)
	if err != nil {
		handleError(err, w)
		return
	}
	balanceService := s.dependencyContainer.NewBalanceService()
	err = balanceService.CreateBalance(userID)
	if err != nil {
		handleError(err, w)
		return
	}
	setOkResponse(w)
}

func (s *server) removeBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromReq(r)
	if err != nil {
		handleError(err, w)
		return
	}
	balanceService := s.dependencyContainer.NewBalanceService()
	err = balanceService.RemoveBalance(userID)
	if err != nil {
		handleError(err, w)
		return
	}
	setOkResponse(w)
}

type createBalanceDTO struct {
	UserID string `json:"userID"`
}

type balanceMovementDTO struct {
	UserID string `json:"userID"`
	Score  int    `json:"score"`
}

type msgResponse struct {
	Msg string `json:"msg"`
}

func setOkResponse(w http.ResponseWriter) {
	b, err := json.Marshal(msgResponse{Msg: "Ok"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = io.WriteString(w, string(b))
}

func getBalanceMovementDataFromReq(r *http.Request) (uuid.UUID, int, error) {
	b, err := getBodyFromReq(r)
	if err != nil {
		return uuid.UUID{}, 0, err
	}
	defer r.Body.Close()
	var dto balanceMovementDTO
	err = json.Unmarshal(b, &dto)
	if err != nil {
		return uuid.UUID{}, 0, err
	}
	userID, err := parseUUID(dto.UserID)
	if err != nil {
		return uuid.UUID{}, 0, err
	}
	return userID, dto.Score, nil
}

func getUserIDFromReq(r *http.Request) (uuid.UUID, error) {
	b, err := getBodyFromReq(r)
	if err != nil {
		return uuid.UUID{}, err
	}
	defer r.Body.Close()
	var dto createBalanceDTO
	err = json.Unmarshal(b, &dto)
	if err != nil {
		return uuid.UUID{}, err
	}
	userID, err := parseUUID(dto.UserID)
	if err != nil {
		return uuid.UUID{}, err
	}
	return userID, nil
}

func getBodyFromReq(r *http.Request) ([]byte, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, ErrBodyParsing
	}
	return b, nil
}

func parseUUID(str string) (uuid.UUID, error) {
	uid, err := uuid.Parse(str)
	if err != nil {
		return uuid.UUID{}, ErrInvalidRequest
	}
	return uid, nil
}

func NewServer(dependencyContainer DependencyContainer) Server {
	return &server{dependencyContainer: dependencyContainer}
}
