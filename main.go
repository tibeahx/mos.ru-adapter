package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/tibeahx/mos.ru-adapter/internal/config"
	"github.com/tibeahx/mos.ru-adapter/internal/handler"
	"github.com/tibeahx/mos.ru-adapter/internal/server"
	storage "github.com/tibeahx/mos.ru-adapter/internal/store"
	logger "github.com/tibeahx/mos.ru-adapter/pkg/log"
	"github.com/tibeahx/mos.ru-adapter/pkg/service"
)

var once sync.Once

func main() {
	cfg := config.GetConfig()

	logger := logger.Zap()

	redis := storage.NewRedisClient(cfg)
	mos := service.NewMosService(cfg, redis, logger)

	handler := handler.NewHandler(mos)

	// если не успели за 5 сек достать из апстрима простыню, идем наху
	// todo: если пошшли назуй, долбимся пока не отдаст
	ctx, err := context.WithTimeout(context.Background(), time.Second*5)
	if err != nil {
		log.Fatal(err)
	}
	// вынес в мейн при старте сервака теперь лезем в апстрим и сохраняем ответ в редис по ключу allParkings
	once.Do(func() { getAndSaveParkingsToRedis(mos, ctx) })

	srv := server.NewServer(cfg, handler, logger)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

	// хз хуйня какая то но вроде должно работать
	defer func() { srv.Stop(ctx) }()
}

func getAndSaveParkingsToRedis(mos *service.MosService, ctx context.Context) {
	mos.GetAllParkingsFromUpstream()
	mos.SaveRowsToCache(ctx)
}
