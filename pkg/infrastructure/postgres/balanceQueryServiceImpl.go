package postgres

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"subscriptions-service/pkg/app"
)

type BalanceQueryServiceImpl struct {
	db *sqlx.DB
}

func (b *BalanceQueryServiceImpl) GetBalanceScoreByUserID(userID uuid.UUID) (int, error) {
	var balanceScore sqlxBalanceScore
	err := b.db.Get(&balanceScore, "SELECT score FROM balance WHERE id_user=$1", userID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, app.ErrBalanceIsNotFound
		}
		return 0, err
	}
	return balanceScore.Score, err
}

func NewBalanceQueryServiceImpl(db *sqlx.DB) app.BalanceQueryService {
	return &BalanceQueryServiceImpl{db: db}
}

type sqlxBalanceScore struct {
	Score int `db:"score"`
}
