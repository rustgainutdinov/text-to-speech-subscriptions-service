package infrastructure

import (
	"github.com/jmoiron/sqlx"
	"subscriptions-service/pkg/app"
	"subscriptions-service/pkg/domain"
)

type DependencyContainer interface {
	newBalanceRepo() domain.BalanceRepo
	NewBalanceService() app.BalanceService
}

type dependencyContainer struct {
	db *sqlx.DB
}

func (d *dependencyContainer) NewBalanceService() app.BalanceService {
	return app.NewBalanceService(d.newBalanceRepo())
}

func (d *dependencyContainer) newBalanceRepo() domain.BalanceRepo {
	return NewBalanceRepo(d.db)
}

func NewDependencyContainer(db *sqlx.DB) DependencyContainer {
	return &dependencyContainer{db: db}
}
