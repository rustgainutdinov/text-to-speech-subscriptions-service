package infrastructure

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"subscriptions-service/pkg/domain"
)

type balanceRepo struct {
	db *sqlx.DB
}

func (c *balanceRepo) FindOne(userID uuid.UUID) (domain.Balance, error) {
	return nil, nil
}

func (c *balanceRepo) Remove(userID uuid.UUID) error {
	return nil
}

func (c *balanceRepo) Store(balance domain.Balance) error {
	return nil
}

func NewBalanceRepo(db *sqlx.DB) domain.BalanceRepo {
	return &balanceRepo{db: db}
}
