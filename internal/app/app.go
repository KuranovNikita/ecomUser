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

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration, secret string, storage *postgres.Storage) *App {

	authService := user.New(log, storage, storage, tokenTTL, secret) //*user.Auth

	grpcApp := grpcapp.New(log, authService, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}
