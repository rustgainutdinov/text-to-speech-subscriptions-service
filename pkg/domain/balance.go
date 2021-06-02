package domain

import (
	"fmt"
	"github.com/google/uuid"
)

type BalanceID uuid.UUID

type BalanceRepo interface {
	FindOne(balanceID BalanceID) (Balance, error)
	Remove(balanceID BalanceID) error
	Store(balance Balance) error
}

type Balance interface {
	ID() BalanceID
	Score() int
	topUp(amountOfSeconds int) error
	writeOff(amountOfSeconds int) error
}

var ErrAmountOfSecondsIsInvalid = fmt.Errorf("amount of seconds is invalid")
var ErrThereAreNotEnoughSecondsOnTheBalance = fmt.Errorf("there are not enough seconds on score")
var ErrBalanceIsNotFound = fmt.Errorf("balance is not found")

type balance struct {
	id    BalanceID
	score int
}

func (c balance) ID() BalanceID {
	return c.id
}

func (c *balance) Score() int {
	return c.score
}

func (c *balance) topUp(amountOfSeconds int) error {
	if amountOfSeconds < 1 {
		return ErrAmountOfSecondsIsInvalid
	}
	c.score += amountOfSeconds
	return nil
}

func (c *balance) writeOff(amountOfSeconds int) error {
	if amountOfSeconds < 1 {
		return ErrAmountOfSecondsIsInvalid
	}
	if c.score-amountOfSeconds < 0 {
		return ErrThereAreNotEnoughSecondsOnTheBalance
	}
	c.score -= amountOfSeconds
	return nil
}

func newBalance(balanceId BalanceID, score *int) Balance {
	b := balance{id: balanceId}
	if score != nil {
		b.score = *score
	}
	return &b
}

func LoadBalance(balanceId BalanceID, score int) Balance {
	return &balance{id: balanceId, score: score}
}
