package domain

import (
	"fmt"
	"github.com/google/uuid"
)

type BalanceRepo interface {
	FindOne(userID uuid.UUID) (Balance, error)
	Remove(userID uuid.UUID) error
	Store(balance Balance) error
}

type Balance interface {
	UserID() uuid.UUID
	Score() int
	topUp(amountOfSymbols int) error
	writeOff(amountOfSymbols int) error
}

var ErrAmountOfSymbolsIsInvalid = fmt.Errorf("amount of symbols is invalid")
var ErrThereAreNotEnoughSymbolsOnTheBalance = fmt.Errorf("there are not enough symbols on score")
var ErrBalanceIsNotFound = fmt.Errorf("balance is not found")

type balance struct {
	userID uuid.UUID
	score  int
}

func (c *balance) Score() int {
	return c.score
}

func (c *balance) UserID() uuid.UUID {
	return c.userID
}

func (c *balance) topUp(amountOfSymbols int) error {
	if amountOfSymbols < 1 {
		return ErrAmountOfSymbolsIsInvalid
	}
	c.score += amountOfSymbols
	return nil
}

func (c *balance) writeOff(amountOfSymbols int) error {
	if amountOfSymbols < 1 {
		return ErrAmountOfSymbolsIsInvalid
	}
	if c.score-amountOfSymbols < 0 {
		return ErrThereAreNotEnoughSymbolsOnTheBalance
	}
	c.score -= amountOfSymbols
	return nil
}

func newBalance(userID uuid.UUID, score *int) Balance {
	b := balance{userID: userID}
	if score != nil {
		b.score = *score
	}
	return &b
}

func LoadBalance(userID uuid.UUID, score int) Balance {
	return &balance{userID: userID, score: score}
}
