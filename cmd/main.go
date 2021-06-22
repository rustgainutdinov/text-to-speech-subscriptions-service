package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"net/http"
	"subscriptions-service/pkg/infrastructure"
	server2 "subscriptions-service/pkg/infrastructure/server"
)

func main() {
	envConf, err := infrastructure.ParseEnv()
	if err != nil {
		log.Fatal(err.Error())
	}
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable", envConf.DBUser, envConf.DBPass, envConf.DBName, envConf.DBPort, envConf.DBHost)
	db, err := sqlx.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.SetFormatter(&log.JSONFormatter{})
	rabbitMqChannel, err := getRabbitMqChannel(envConf)
	if err != nil {
		log.Fatal(err.Error())
	}
	dependencyContainer, err := infrastructure.NewDependencyContainer(db, rabbitMqChannel)
	if err != nil {
		log.Fatal(err.Error())
	}
	server := server2.NewServer(dependencyContainer)
	handler := server2.Router(server)
	srv := &http.Server{Addr: envConf.ServeRESTAddress, Handler: handler}
	log.WithFields(log.Fields{"port": envConf.ServeRESTAddress}).Info("Successful starting")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
	_ = srv.Shutdown(context.Background())
}

func getRabbitMqChannel(envConf *infrastructure.Config) (*amqp.Channel, error) {
	rabbitMqInfo := fmt.Sprintf("amqp://%s:%s@%s//", envConf.RabbitMqUser, envConf.RabbitMqPass, envConf.RabbitMqHost)
	conn, err := amqp.Dial(rabbitMqInfo)
	if err != nil {
		return nil, err
	}
	return conn.Channel()
}
