package infrastructure

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"subscriptions-service/pkg/app"
	"subscriptions-service/pkg/domain"
)

type BalanceQueryServiceImpl struct {
	db *sqlx.DB
}

func (b *BalanceQueryServiceImpl) GetBalanceScoreByUserID(userID uuid.UUID) (int, error) {
	var balanceScores []sqlxBalanceScore
	err := b.db.Select(&balanceScores, "SELECT score FROM balance WHERE id_user=$1 LIMIT 1", userID.String())
	if err != nil {
		return 0, err
	}
	if len(balanceScores) == 0 {
		return 0, domain.ErrBalanceIsNotFound
	}
	return balanceScores[0].Score, err
}

func NewBalanceQueryServiceImpl(db *sqlx.DB) app.BalanceQueryService {
	return &BalanceQueryServiceImpl{db: db}
}

type sqlxBalanceScore struct {
	Score int `db:"score"`
}
