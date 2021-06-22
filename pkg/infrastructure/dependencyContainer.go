package infrastructure

import (
	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"
	"subscriptions-service/pkg/app"
	"subscriptions-service/pkg/domain"
	"subscriptions-service/pkg/infrastructure/eventHandler"
	"subscriptions-service/pkg/infrastructure/postgres"
)

type DependencyContainer interface {
	NewBalanceService() app.BalanceService
}

type dependencyContainer struct {
	db                    *sqlx.DB
	externalEventListener *eventHandler.ExternalEventListener
	balanceService        app.BalanceService
}

func (d *dependencyContainer) NewBalanceService() app.BalanceService {
	if d.balanceService == nil {
		d.balanceService = app.NewBalanceService(d.newBalanceRepo(), d.newBalanceQueryService())
	}
	return d.balanceService
}

func (d *dependencyContainer) newBalanceRepo() domain.BalanceRepo {
	return postgres.NewBalanceRepo(d.db)
}

func (d *dependencyContainer) newBalanceQueryService() app.BalanceQueryService {
	return postgres.NewBalanceQueryServiceImpl(d.db)
}

func NewDependencyContainer(db *sqlx.DB, rabbitMqChannel *amqp.Channel) (DependencyContainer, error) {
	dp := &dependencyContainer{db: db}
	externalEventListener, err := eventHandler.NewExternalEventListener(rabbitMqChannel, dp.newBalanceRepo())
	if err != nil {
		return nil, err
	}
	if externalEventListener != nil {
		externalEventListener.ActivateExternalEventListener()
		dp.externalEventListener = externalEventListener
	}
	return dp, nil
}
