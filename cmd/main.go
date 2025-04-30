package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"my-grpc/internal/app"
	"my-grpc/internal/config"
	"my-grpc/internal/lib/logger"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.InitConfig()

	if err != nil {
		log.Fatalf("config is not set:%s", err)
	}

	log := logger.SetLogger(cfg.Env)

	if err != nil {
		log.Error(err.Error())
	}

	log.Info("<app are started>")

	application := app.New(log, cfg.GRPC.Port, cfg.DB, cfg.TokenTTL)

	go application.GRPCServer.Run()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GRPCServer.Stop()

}
