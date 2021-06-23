package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"subscriptions-service/api"
	"subscriptions-service/pkg/infrastructure"
	server2 "subscriptions-service/pkg/infrastructure/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}

func run() error {
	envConf, err := infrastructure.ParseEnv()
	log.SetFormatter(&log.JSONFormatter{})
	if err != nil {
		return err
	}
	go func() {
		if err := runGRPCService(envConf); err != nil {
			log.Fatal(err.Error())
		}
	}()
	return runHTTPProxy(envConf.GRPCAddress, envConf.HTTPProxyAddress)
}

func runGRPCService(envConf *infrastructure.Config) error {
	rabbitMqChannel, err := getRabbitMqChannel(envConf)
	if err != nil {
		return err
	}
	db, err := getDataBaseConnect(envConf)
	if err != nil {
		return err
	}
	dependencyContainer, err := infrastructure.NewDependencyContainer(db, rabbitMqChannel)
	if err != nil {
		return err
	}
	lis, err := net.Listen("tcp", envConf.GRPCAddress)
	if err != nil {
		return err
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(makeGRPCUnaryInterceptor()))
	api.RegisterTranslationServiceServer(server, &server2.BalanceServer{DependencyContainer: dependencyContainer})
	log.WithFields(log.Fields{"grpc address": envConf.GRPCAddress}).Info("successfully starting grpc transport")
	return server.Serve(lis)
}

func getDataBaseConnect(envConf *infrastructure.Config) (*sqlx.DB, error) {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable", envConf.DBUser, envConf.DBPass, envConf.DBName, envConf.DBPort, envConf.DBHost)
	db, err := sqlx.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getRabbitMqChannel(envConf *infrastructure.Config) (*amqp.Channel, error) {
	rabbitMqInfo := fmt.Sprintf("amqp://%s:%s@%s//", envConf.RabbitMqUser, envConf.RabbitMqPass, envConf.RabbitMqHost)
	conn, err := amqp.Dial(rabbitMqInfo)
	if err != nil {
		return nil, err
	}
	return conn.Channel()
}

func makeGRPCUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		return resp, server2.TranslateError(err)
	}
}

func runHTTPProxy(serviceAddr string, httpProxyPort string) error {
	grpcConn, err := grpc.Dial(serviceAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer grpcConn.Close()
	grpcGWMux := runtime.NewServeMux()
	err = api.RegisterTranslationServiceHandler(context.Background(), grpcGWMux, grpcConn)
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcGWMux)
	log.WithFields(log.Fields{"http port": httpProxyPort}).Info("successfully starting http transport")
	return http.ListenAndServe(httpProxyPort, mux)
}
