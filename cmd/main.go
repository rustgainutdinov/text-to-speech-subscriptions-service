package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
	"subscriptions-service/pkg/infrastructure"
)

func main() {
	envConf, err := infrastructure.ParseEnv()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", envConf.DBUser, envConf.DBPass, envConf.DBName, envConf.DBPort)
	db, err := sqlx.Open("postgres", dbInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	dependencyContainer := infrastructure.NewDependencyContainer(db)
	server := infrastructure.NewServer(dependencyContainer)
	handler := infrastructure.Router(server)
	srv := &http.Server{Addr: envConf.ServeRESTAddress, Handler: handler}
	fmt.Println(srv.ListenAndServe())
	_ = srv.Shutdown(context.Background())
}
