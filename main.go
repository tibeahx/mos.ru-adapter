package main

import (
	"context"
	"log"

	"github.com/tibeahx/mos.ru-adapter/internal/config"
	"github.com/tibeahx/mos.ru-adapter/internal/handler"
	"github.com/tibeahx/mos.ru-adapter/internal/server"

	logger "github.com/tibeahx/mos.ru-adapter/pkg/log"
	"github.com/tibeahx/mos.ru-adapter/pkg/svc/mos"
	"github.com/tibeahx/mos.ru-adapter/pkg/svc/mos/mosclient"
	"github.com/tibeahx/mos.ru-adapter/pkg/svc/redis"
)

func main() {
	cfg := config.GetConfig()
	logger := logger.Zap()
	rc := redis.NewRC(cfg)
	mosClient := mosclient.NewMosClient(cfg, logger)
	mos := mos.NewMossvc(cfg, rc, logger, mosClient)
	handler := handler.NewHandler(mos)

	ctx := context.Background()

	go func() {
		if err := mos.SaveRowsToCache(ctx); err != nil {
			log.Print("failed to save rows to redis")
			return
		}
	}()

	srv := server.NewServer(cfg, handler, logger)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

	// хз хуйня какая то но вроде должно работать
	// defer func() { srv.Stop(ctx) }()
}
