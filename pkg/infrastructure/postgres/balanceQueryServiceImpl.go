package postgres

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"subscriptions-service/pkg/app"
)

type BalanceQueryServiceImpl struct {
	db *sqlx.DB
}

func (b *BalanceQueryServiceImpl) GetBalanceScoreByUserID(userID uuid.UUID) (int, error) {
	var balanceScore sqlxBalanceScore
	err := b.db.Get(&balanceScore, "SELECT score FROM balance WHERE id_user=$1", userID.String())
	if err != nil {
		return 0, err
	}
	//TODO: добавить обработку ошибки not found (ErrTranslationIsNotFound)
	return balanceScore.Score, err
}

func NewBalanceQueryServiceImpl(db *sqlx.DB) app.BalanceQueryService {
	return &BalanceQueryServiceImpl{db: db}
}

type sqlxBalanceScore struct {
	Score int `db:"score"`
}
