package infrastructure

import (
	"github.com/jmoiron/sqlx"
	"subscriptions-service/pkg/app"
	"subscriptions-service/pkg/domain"
)

type DependencyContainer interface {
	newBalanceService() app.BalanceService
}

type dependencyContainer struct {
	db *sqlx.DB
}

func (d *dependencyContainer) newBalanceService() app.BalanceService {
	return app.NewBalanceService(d.newBalanceRepo(), d.newBalanceQueryService())
}

func (d *dependencyContainer) newBalanceRepo() domain.BalanceRepo {
	return NewBalanceRepo(d.db)
}

func (d *dependencyContainer) newBalanceQueryService() app.BalanceQueryService {
	return NewBalanceQueryServiceImpl(d.db)
}

func NewDependencyContainer(db *sqlx.DB) DependencyContainer {
	return &dependencyContainer{db: db}
}
