package main

import (
	"context"
	"log"

	"github.com/tibeahx/mos.ru-adapter/internal/config"
	"github.com/tibeahx/mos.ru-adapter/internal/handler"
	"github.com/tibeahx/mos.ru-adapter/internal/server"
	storage "github.com/tibeahx/mos.ru-adapter/internal/store"
	logger "github.com/tibeahx/mos.ru-adapter/pkg/log"
	"github.com/tibeahx/mos.ru-adapter/pkg/service"
)

func main() {
	cfg := config.GetConfig()

	logger := logger.Zap()

	redis := storage.NewRedisClient(cfg)
	mos := service.NewMosService(cfg, redis, logger)

	handler := handler.NewHandler(mos)
	// если не успели за 5 сек достать из апстрима простыню, идем наху
	// todo: если пошшли назуй, долбимся пока не отдаст
	ctx := context.Background()
	// вынес в мейн при старте сервака теперь лезем в апстрим и сохраняем ответ в редис по ключу allParkings
	mos.GetAllParkingsFromUpstream()
	bytes, err := mos.SaveRowsToCache(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("saved %d bytes in redis", bytes)

	srv := server.NewServer(cfg, handler, logger)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

	// хз хуйня какая то но вроде должно работать
	// defer func() { srv.Stop(ctx) }()
}
