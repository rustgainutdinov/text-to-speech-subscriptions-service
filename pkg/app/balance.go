package app

import (
	"github.com/google/uuid"
	"subscriptions-service/pkg/domain"
)

type BalanceService interface {
	CreateBalance(userID uuid.UUID) error
	RemoveBalance(userID uuid.UUID) error
	TopUpBalance(userID uuid.UUID, score int) error
	WriteOffFromBalance(userID uuid.UUID, score int) error
}

type balanceService struct {
	balanceRepo domain.BalanceRepo
}

func (b *balanceService) CreateBalance(userID uuid.UUID) error {
	return domain.NewBalanceService(b.balanceRepo).CreateBalance(userID)
}

func (b *balanceService) RemoveBalance(userID uuid.UUID) error {
	return domain.NewBalanceService(b.balanceRepo).RemoveBalance(userID)
}

func (b *balanceService) TopUpBalance(userID uuid.UUID, score int) error {
	return domain.NewBalanceService(b.balanceRepo).TopUpBalance(userID, score)
}

func (b *balanceService) WriteOffFromBalance(userID uuid.UUID, score int) error {
	return domain.NewBalanceService(b.balanceRepo).WriteOffFromBalance(userID, score)
}

func NewBalanceService(balanceRepo domain.BalanceRepo) BalanceService {
	return &balanceService{balanceRepo: balanceRepo}
}
