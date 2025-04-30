package app

import (
	"log/slog"
	grpcapp "my-grpc/internal/app/grpc"
	"my-grpc/internal/config"
	"my-grpc/internal/repository/postgres"
	"my-grpc/internal/service/auth"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	db config.DBConfig,
	tokenTTL time.Duration,
) *App {
	storage, err := postgres.NewPostgresDB(db)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		GRPCServer: grpcApp,
	}
}
