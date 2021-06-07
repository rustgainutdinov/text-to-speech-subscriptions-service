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
	var balance []sqlxBalance
	err := c.db.Select(&balance, "SELECT * FROM balance WHERE id_user=$1 LIMIT 1", userID.String())
	if err != nil {
		return nil, err
	}
	return domain.LoadBalance(userID, balance[0].Score), err
}

func (c *balanceRepo) Remove(userID uuid.UUID) error {
	return nil
}

func (c *balanceRepo) Store(balance domain.Balance) error {
	_, err := c.db.Query(
		`INSERT INTO balance (id_user, score)
				VALUES ($1, $2)
				ON CONFLICT (id_user)
					DO UPDATE SET score = $2;`,
		balance.UserID().String(), balance.Score())
	return err
}

func NewBalanceRepo(db *sqlx.DB) domain.BalanceRepo {
	return &balanceRepo{db: db}
}

type sqlxBalance struct {
	UserID string `db:"id_user"`
	Score  int    `db:"score"`
}
