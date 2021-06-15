package app

import (
	"fmt"
	"github.com/google/uuid"
)

var ErrBalanceIsNotFound = fmt.Errorf("balance is not found")

type BalanceQueryService interface {
	GetBalanceScoreByUserID(userID uuid.UUID) (int, error)
}
