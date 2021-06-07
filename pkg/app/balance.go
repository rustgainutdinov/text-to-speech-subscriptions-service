package app

import (
	"github.com/google/uuid"
	"subscriptions-service/pkg/domain"
)

type BalanceService interface {
	CreateBalance(userID uuid.UUID) error
}

type balanceService struct {
	balanceRepo domain.BalanceRepo
}

func (b *balanceService) CreateBalance(userID uuid.UUID) error {
	return domain.NewBalanceService(b.balanceRepo).CreateBalance(userID)
}

func NewBalanceService(balanceRepo domain.BalanceRepo) BalanceService {
	return &balanceService{balanceRepo: balanceRepo}
}
