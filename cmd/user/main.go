package main

import (
	"ecomUser/internal/app"
	"ecomUser/internal/config"
	"ecomUser/internal/storage/postgres"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener")
	log.Debug("debug messages are enabled")

	storagePath := postgres.SplitStoragePath(cfg.LoginDB, cfg.PasswordDB, cfg.HostDB, cfg.PortDB, cfg.NameDB)

	storage, err := postgres.New(storagePath)

	if err != nil {
		panic(err)
	}

	defer storage.Close()

	application := app.New(log, cfg.GRPCPort, storagePath, cfg.GRPCTimeout, cfg.JWTSecret, storage)

	go application.GRPCSrv.MustRun()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signal := <-stop

	log.Info("stopping app", slog.String("signal", signal.String()))

	application.GRPCSrv.Stop()
	log.Info("application stop")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
