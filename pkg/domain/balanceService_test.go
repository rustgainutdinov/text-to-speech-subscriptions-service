package domain

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBalanceService_CreateBalance(t *testing.T) {
	repo := newMockBalanceRepo()
	service := NewBalanceService(repo)
	balanceId := BalanceID(uuid.New())
	err := service.CreateBalance(balanceId)
	assert.Nil(t, err)
	balance, err := repo.FindOne(balanceId)
	assert.Nil(t, err)
	assert.Equal(t, balanceId, balance.ID())
	assert.Equal(t, 0, balance.Score())
}

func TestBalanceService_RemoveBalance(t *testing.T) {
	repo := newMockBalanceRepo()
	service := NewBalanceService(repo)
	balanceId := BalanceID(uuid.New())
	err := service.CreateBalance(balanceId)
	assert.Nil(t, err)
	_, err = repo.FindOne(balanceId)
	assert.Nil(t, err)
	err = service.RemoveBalance(balanceId)
	assert.Nil(t, err)
	_, err = repo.FindOne(balanceId)
	assert.Equal(t, ErrBalanceIsNotFound, err)
}

func TestBalanceService_TopUpBalance(t *testing.T) {
	repo := newMockBalanceRepo()
	service := NewBalanceService(repo)
	balanceId := BalanceID(uuid.New())
	err := service.CreateBalance(balanceId)
	assert.Nil(t, err)
	amountOfSeconds := 27
	err = service.TopUpBalance(balanceId, amountOfSeconds)
	assert.Nil(t, err)
	balance, err := repo.FindOne(balanceId)
	assert.Nil(t, err)
	assert.NotNil(t, balance)
	assert.Equal(t, amountOfSeconds, balance.Score())
}

func TestBalanceService_WriteOffFromBalance(t *testing.T) {
	repo := newMockBalanceRepo()
	service := NewBalanceService(repo)
	balanceId := BalanceID(uuid.New())
	err := service.CreateBalance(balanceId)
	assert.Nil(t, err)
	amountOfSeconds := 88
	err = service.TopUpBalance(balanceId, amountOfSeconds)
	assert.Nil(t, err)
	amountOfSecondsToWriteOff := 34
	err = service.WriteOffFromBalance(balanceId, amountOfSecondsToWriteOff)
	assert.Nil(t, err)
	amountOfSeconds -= amountOfSecondsToWriteOff
	balance, err := repo.FindOne(balanceId)
	assert.Nil(t, err)
	assert.Equal(t, amountOfSeconds, balance.Score())
	err = service.WriteOffFromBalance(balanceId, amountOfSeconds+1)
	assert.Equal(t, err, ErrThereAreNotEnoughSecondsOnTheBalance)

}

type mockBalanceRepo struct {
	balances []Balance
}

func (c *mockBalanceRepo) FindOne(balanceID BalanceID) (Balance, error) {
	for _, balance := range c.balances {
		if balance.ID() == balanceID {
			return balance, nil
		}
	}
	return nil, ErrBalanceIsNotFound
}

func (c *mockBalanceRepo) Remove(balanceID BalanceID) error {
	for i, balance := range c.balances {
		if balance.ID() == balanceID {
			copy(c.balances[i:], c.balances[i+1:])
			c.balances = c.balances[:len(c.balances)-1]
			return nil
		}
	}
	return nil
}

func (c *mockBalanceRepo) Store(balance Balance) error {
	for i, repoBalance := range c.balances {
		if repoBalance.ID() == balance.ID() {
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
