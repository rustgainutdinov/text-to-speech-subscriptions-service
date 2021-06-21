package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"net/http"
	server2 "subscriptions-service/pkg/infrastructure/server"
)

func main() {
	envConf, err := server2.ParseEnv()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable", envConf.DBUser, envConf.DBPass, envConf.DBName, envConf.DBPort, envConf.DBHost)
	db, err := sqlx.Open("postgres", dbInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	rabbitMqChannel, err := getRabbitMqChannel(envConf)
	if err != nil {
		fmt.Println(err.Error())
	}
	dependencyContainer := server2.NewDependencyContainer(db, rabbitMqChannel)
	server := server2.NewServer(dependencyContainer)
	handler := server2.Router(server)
	srv := &http.Server{Addr: envConf.ServeRESTAddress, Handler: handler}
	fmt.Println(srv.ListenAndServe())
	_ = srv.Shutdown(context.Background())
}

func getRabbitMqChannel(envConf *server2.Config) (*amqp.Channel, error) {
	rabbitMqInfo := fmt.Sprintf("amqp://%s:%s@%s//", envConf.RabbitMqUser, envConf.RabbitMqPass, envConf.RabbitMqHost)
	conn, err := amqp.Dial(rabbitMqInfo)
	if err != nil {
		return nil, err
	}
	return conn.Channel()
}
