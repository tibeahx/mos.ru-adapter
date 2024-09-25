package main

import (
	"context"
	"log"
	"sync"
	"test-task/internal/config"
	"test-task/internal/handler"
	"test-task/internal/server"
	storage "test-task/internal/store"
	logger "test-task/pkg/log"
	"test-task/pkg/service"
	"time"
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
