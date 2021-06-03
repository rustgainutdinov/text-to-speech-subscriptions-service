package infrastructure

import (
	"github.com/jmoiron/sqlx"
	"subscriptions-service/pkg/domain"
)

type DependencyContainer interface {
	NewBalanceRepo() domain.BalanceRepo
}

type dependencyContainer struct {
	db *sqlx.DB
}

func (d *dependencyContainer) NewBalanceRepo() domain.BalanceRepo {
	return NewBalanceRepo(d.db)
}

func NewDependencyContainer(db *sqlx.DB) DependencyContainer {
	return &dependencyContainer{db: db}
}
