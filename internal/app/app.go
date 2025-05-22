package app

import (
	grpcapp "ecomUser/internal/app/grpc"
	"ecomUser/internal/services/user"
	"ecomUser/internal/storage/postgres"
	"log/slog"

	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {

	storage, err := postgres.New(storagePath)

	if err != nil {
		panic(err)
	}

	defer storage.Close()

	authService := user.New(log, storage, storage, tokenTTL) //*user.Auth

	grpcApp := grpcapp.New(log, authService, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}
