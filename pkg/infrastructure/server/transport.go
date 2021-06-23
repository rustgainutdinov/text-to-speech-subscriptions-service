package server

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"subscriptions-service/api"
	"subscriptions-service/pkg/infrastructure"
)

type BalanceServer struct {
	DependencyContainer infrastructure.DependencyContainer
}

var ErrInvalidUUID = fmt.Errorf("invalid uuid")

func (t *BalanceServer) CreateBalance(_ context.Context, req *api.UserID) (*empty.Empty, error) {
	userID, err := parseUUID(req.UserID)
	if err != nil {
		return &empty.Empty{}, err
	}
	balanceService := t.DependencyContainer.NewBalanceService()
	err = balanceService.CreateBalance(userID)
	return &empty.Empty{}, err
}

func (t *BalanceServer) TopUpBalance(_ context.Context, req *api.BalanceMovementData) (*empty.Empty, error) {
	userID, err := parseUUID(req.UserID)
	if err != nil {
		return &empty.Empty{}, err
	}
	balanceService := t.DependencyContainer.NewBalanceService()
	err = balanceService.TopUpBalance(userID, int(req.Score))
	return &empty.Empty{}, err
}

func (t *BalanceServer) WriteOffFromBalance(_ context.Context, req *api.BalanceMovementData) (*empty.Empty, error) {
	userID, err := parseUUID(req.UserID)
	if err != nil {
		return &empty.Empty{}, err
	}
	balanceService := t.DependencyContainer.NewBalanceService()
	err = balanceService.WriteOffFromBalance(userID, int(req.Score))
	return &empty.Empty{}, err

}

func (t *BalanceServer) RemoveBalance(_ context.Context, req *api.UserID) (*empty.Empty, error) {
	userID, err := parseUUID(req.UserID)
	if err != nil {
		return &empty.Empty{}, err
	}
	balanceService := t.DependencyContainer.NewBalanceService()
	err = balanceService.RemoveBalance(userID)
	return &empty.Empty{}, err
}

func (t *BalanceServer) CanWriteOffFromBalance(_ context.Context, req *api.BalanceMovementData) (*api.CanWriteOffFromBalanceResponse, error) {
	userID, err := parseUUID(req.UserID)
	if err != nil {
		return &api.CanWriteOffFromBalanceResponse{}, err
	}
	balanceService := t.DependencyContainer.NewBalanceService()
	canWriteOff, err := balanceService.CanWriteOffFromBalance(userID, int(req.Score))
	if err != nil {
		return &api.CanWriteOffFromBalanceResponse{}, err
	}
	log.Info(canWriteOff)
	return &api.CanWriteOffFromBalanceResponse{Result: &wrappers.BoolValue{Value: canWriteOff}}, nil
}

func parseUUID(str string) (uuid.UUID, error) {
	uid, err := uuid.Parse(str)
	if err != nil {
		return uuid.UUID{}, ErrInvalidUUID
	}
	return uid, nil
}
