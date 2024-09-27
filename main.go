package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := mos.SaveRowsToCache(ctx); err != nil {
			log.Print("failed to save rows to redis")
			return
		}

		parkings, err := mos.GetParkingsFromStorage(ctx)
		if err != nil {
			return
		}

		j, _ := json.Marshal(parkings)
		fmt.Println(string(j))
	}()
	wg.Wait()

	srv := server.NewServer(cfg, handler, logger)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		if err := srv.Stop(ctx); err != nil {
			log.Fatal(err)
		}
	}()
}
