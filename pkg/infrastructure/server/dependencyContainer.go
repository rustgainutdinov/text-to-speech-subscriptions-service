package server

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"
	"subscriptions-service/pkg/app"
	"subscriptions-service/pkg/domain"
	"subscriptions-service/pkg/infrastructure/externalService"
	"subscriptions-service/pkg/infrastructure/postgres"
)

type DependencyContainer interface {
	newBalanceService() app.BalanceService
}

type dependencyContainer struct {
	db                    *sqlx.DB
	externalEventListener *externalService.ExternalEventListener
}

func (d *dependencyContainer) newBalanceService() app.BalanceService {
	return app.NewBalanceService(d.newBalanceRepo(), d.newBalanceQueryService())
}

func (d *dependencyContainer) newBalanceRepo() domain.BalanceRepo {
	return postgres.NewBalanceRepo(d.db)
}

func (d *dependencyContainer) newBalanceQueryService() app.BalanceQueryService {
	return postgres.NewBalanceQueryServiceImpl(d.db)
}

func NewDependencyContainer(db *sqlx.DB, rabbitMqChannel *amqp.Channel) DependencyContainer {
	dp := &dependencyContainer{db: db}
	externalEventListener, err := externalService.NewExternalEventListener(rabbitMqChannel, dp.newBalanceRepo())
	if err != nil {
		fmt.Println(err.Error())
	}
	if externalEventListener != nil {
		externalEventListener.ActivateTextTranslatedHandler()
		dp.externalEventListener = externalEventListener
	}
	return dp
}
