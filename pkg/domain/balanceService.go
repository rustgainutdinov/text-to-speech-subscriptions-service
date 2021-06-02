package domain

import "fmt"

type BalanceService interface {
	CreateBalance(balanceID BalanceID) error
	RemoveBalance(balanceID BalanceID) error
	TopUpBalance(balanceID BalanceID, amountOfSeconds int) error
	WriteOffFromBalance(balanceID BalanceID, amountOfSeconds int) error
}

var ErrBalanceIsAlreadyExists = fmt.Errorf("balance is already exists")

type balanceService struct {
	repo BalanceRepo
}

func (c *balanceService) CreateBalance(balanceID BalanceID) error {
	foundedBalance, err := c.repo.FindOne(balanceID)
	if err != nil && err != ErrBalanceIsNotFound {
		return err
	}
	if foundedBalance != nil {
		return ErrBalanceIsAlreadyExists
	}
	balance := newBalance(balanceID, nil)
	return c.repo.Store(balance)
}

func (c *balanceService) RemoveBalance(balanceID BalanceID) error {
	foundedBalance, err := c.repo.FindOne(balanceID)
	if err != nil {
		return err
	}
	if foundedBalance == nil {
		return ErrBalanceIsNotFound
	}
	return c.repo.Remove(balanceID)
}

func (c *balanceService) TopUpBalance(balanceID BalanceID, amountOfSeconds int) error {
	balance, err := c.repo.FindOne(balanceID)
	if err != nil {
		return err
	}
	err = balance.topUp(amountOfSeconds)
	if err != nil {
		return err
	}
	return c.repo.Store(balance)
}

func (c *balanceService) WriteOffFromBalance(balanceID BalanceID, amountOfSeconds int) error {
	balance, err := c.repo.FindOne(balanceID)
	if err != nil {
		return err
	}
	err = balance.writeOff(amountOfSeconds)
	if err != nil {
		return err
	}
	return c.repo.Store(balance)
}

func NewBalanceService(repo BalanceRepo) BalanceService {
	return &balanceService{repo: repo}
}
