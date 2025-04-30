package grpcapp

import (
	"fmt"
	"log/slog"

	authgrpc "my-grpc/internal/server"

	"net"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	Port       int
	gRPCServer *grpc.Server
}

func New(log *slog.Logger, port int, authService authgrpc.Auth) *App {
	grpcServer := grpc.NewServer()
	authgrpc.Register(grpcServer, authService)

	return &App{
		log:        log,
		Port:       port,
		gRPCServer: grpcServer,
	}

}

func (a *App) Run() error {
	const op = "app.run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.Port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.Port))

	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("%s:%w", op, err)
	}

	if err := a.gRPCServer.Serve(l); err != nil {
		log.Error(err.Error())
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil

}

func (a *App) Stop() error {
	const op = "app.stop"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("app is stopped")

	a.gRPCServer.GracefulStop()

	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}
