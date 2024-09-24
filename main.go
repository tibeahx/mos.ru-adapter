package main

import (
	"context"
	"log"
	"test-task/internal/config"
	"test-task/internal/handler"

	"test-task/internal/server"
	storage "test-task/internal/store"
	logger "test-task/pkg/log"
	"test-task/pkg/service"
	"time"
)

func main() {
	cfg := config.GetConfig()

	logger := logger.Zap()

	redis := storage.NewRedisClient(cfg)
	mos := service.NewMosService(cfg, redis, logger)

	handler := handler.NewHandler(mos)

	srv := server.NewServer(cfg, handler, logger)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

	ctx, err := context.WithTimeout(context.Background(), time.Second*5)
	if err != nil {
		log.Fatal(err)
	}

	srv.Stop(ctx)
}
