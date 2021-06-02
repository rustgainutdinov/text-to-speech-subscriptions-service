package domain

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBalanceService_CreateBalance(t *testing.T) {
	repo := newMockBalanceRepo()
	service := NewBalanceService(repo)
	userId := uuid.New()
	err := service.CreateBalance(userId)
	assert.Nil(t, err)
	balance, err := repo.FindOne(userId)
	assert.Nil(t, err)
	assert.Equal(t, userId, balance.UserID())
	assert.Equal(t, 0, balance.Score())
}

func TestBalanceService_RemoveBalance(t *testing.T) {
	repo := newMockBalanceRepo()
	service := NewBalanceService(repo)
	userId := uuid.New()
	err := service.CreateBalance(userId)
	assert.Nil(t, err)
	_, err = repo.FindOne(userId)
	assert.Nil(t, err)
	err = service.RemoveBalance(userId)
	assert.Nil(t, err)
	_, err = repo.FindOne(userId)
	assert.Equal(t, ErrBalanceIsNotFound, err)
}

func TestBalanceService_TopUpBalance(t *testing.T) {
	repo := newMockBalanceRepo()
	service := NewBalanceService(repo)
	userID := uuid.New()
	err := service.CreateBalance(userID)
	assert.Nil(t, err)
	amountOfSymbols := 27
	err = service.TopUpBalance(userID, amountOfSymbols)
	assert.Nil(t, err)
	balance, err := repo.FindOne(userID)
	assert.Nil(t, err)
	assert.NotNil(t, balance)
	assert.Equal(t, amountOfSymbols, balance.Score())
}

func TestBalanceService_WriteOffFromBalance(t *testing.T) {
	repo := newMockBalanceRepo()
	service := NewBalanceService(repo)
	userID := uuid.New()
	err := service.CreateBalance(userID)
	assert.Nil(t, err)
	amountOfSymbols := 88
	err = service.TopUpBalance(userID, amountOfSymbols)
	assert.Nil(t, err)
	amountOfSymbolsToWriteOff := 34
	err = service.WriteOffFromBalance(userID, amountOfSymbolsToWriteOff)
	assert.Nil(t, err)
	amountOfSymbols -= amountOfSymbolsToWriteOff
	balance, err := repo.FindOne(userID)
	assert.Nil(t, err)
	assert.Equal(t, amountOfSymbols, balance.Score())
	err = service.WriteOffFromBalance(userID, amountOfSymbols+1)
	assert.Equal(t, err, ErrThereAreNotEnoughSymbolsOnTheBalance)

}

type mockBalanceRepo struct {
	balances []Balance
}

func (c *mockBalanceRepo) FindOne(userID uuid.UUID) (Balance, error) {
	for _, balance := range c.balances {
		if balance.UserID() == userID {
			return balance, nil
		}
	}
	return nil, ErrBalanceIsNotFound
}

func (c *mockBalanceRepo) Remove(userID uuid.UUID) error {
	for i, balance := range c.balances {
		if balance.UserID() == userID {
			copy(c.balances[i:], c.balances[i+1:])
			c.balances = c.balances[:len(c.balances)-1]
			return nil
		}
	}
	return nil
}

func (c *mockBalanceRepo) Store(balance Balance) error {
	for i, repoBalance := range c.balances {
		if repoBalance.UserID() == balance.UserID() {
			c.balances[i] = balance
			return nil
		}
	}
	c.balances = append(c.balances, balance)
	return nil
}

func newMockBalanceRepo() BalanceRepo {
	return &mockBalanceRepo{}
}
