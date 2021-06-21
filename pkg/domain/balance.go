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
	topUp(score int) error
	writeOff(score int) error
}

var ErrScoreIsInvalid = fmt.Errorf("amount of symbols is invalid")
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

func (c *balance) topUp(score int) error {
	if score < 1 {
		return ErrScoreIsInvalid
	}
	c.score += score
	return nil
}

func (c *balance) writeOff(score int) error {
	if score < 1 {
		return ErrScoreIsInvalid
	}
	if c.score-score < 0 {
		return ErrThereAreNotEnoughSymbolsOnTheBalance
	}
	c.score -= score
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
