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
	CanWriteOffFromBalance(userID uuid.UUID, score int) (bool, error)
}

type balanceService struct {
	balanceRepo  domain.BalanceRepo
	queryService BalanceQueryService
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

func (b *balanceService) CanWriteOffFromBalance(userID uuid.UUID, score int) (bool, error) {
	currScore, err := b.queryService.GetBalanceScoreByUserID(userID)
	if err != nil {
		return false, err
	}
	return currScore-score >= 0, err
}

func NewBalanceService(balanceRepo domain.BalanceRepo, queryService BalanceQueryService) BalanceService {
	return &balanceService{balanceRepo: balanceRepo, queryService: queryService}
}
