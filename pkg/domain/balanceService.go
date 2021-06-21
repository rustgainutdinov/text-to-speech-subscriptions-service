package domain

import (
	"fmt"
	"github.com/google/uuid"
)

type BalanceService interface {
	CreateBalance(userID uuid.UUID) error
	RemoveBalance(userID uuid.UUID) error
	TopUpBalance(userID uuid.UUID, score int) error
	WriteOffFromBalance(userID uuid.UUID, score int) error
}

var ErrBalanceIsAlreadyExists = fmt.Errorf("balance is already exists")

type balanceService struct {
	repo BalanceRepo
}

func (c *balanceService) CreateBalance(userID uuid.UUID) error {
	foundedBalance, err := c.repo.FindOne(userID)
	if err != nil && err != ErrBalanceIsNotFound {
		return err
	}
	if foundedBalance != nil {
		return ErrBalanceIsAlreadyExists
	}
	balance := newBalance(userID, nil)
	return c.repo.Store(balance)
}

func (c *balanceService) RemoveBalance(userID uuid.UUID) error {
	foundedBalance, err := c.repo.FindOne(userID)
	if err != nil {
		return err
	}
	return c.repo.Remove(foundedBalance.UserID())
}

func (c *balanceService) TopUpBalance(userID uuid.UUID, score int) error {
	balance, err := c.repo.FindOne(userID)
	if err != nil {
		return err
	}
	err = balance.topUp(score)
	if err != nil {
		return err
	}
	return c.repo.Store(balance)
}

func (c *balanceService) WriteOffFromBalance(userID uuid.UUID, score int) error {
	balance, err := c.repo.FindOne(userID)
	if err != nil {
		return err
	}
	err = balance.writeOff(score)
	if err != nil {
		return err
	}
	return c.repo.Store(balance)
}

func NewBalanceService(repo BalanceRepo) BalanceService {
	return &balanceService{repo: repo}
}
