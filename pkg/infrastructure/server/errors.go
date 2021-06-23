package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"subscriptions-service/pkg/app"
	"subscriptions-service/pkg/domain"
)

func TranslateError(err error) error {
	if _, ok := status.FromError(err); ok {
		return err
	}
	if err == domain.ErrScoreIsInvalid ||
		err == domain.ErrBalanceIsNotFound ||
		err == domain.ErrThereAreNotEnoughSymbolsOnTheBalance ||
		err == ErrInvalidUUID ||
		err == app.ErrBalanceIsNotFound {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}
	if err == domain.ErrBalanceIsAlreadyExists {
		return status.Errorf(codes.AlreadyExists, err.Error())
	}
	return status.Errorf(codes.Internal, err.Error())
}
