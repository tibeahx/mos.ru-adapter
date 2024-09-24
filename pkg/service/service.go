package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"test-task/internal/config"
	storage "test-task/internal/store"
	"test-task/internal/types"
	"time"

	"go.uber.org/zap"
)

var (
	defaultApiKey = os.Getenv("defaultApiKey")
	defaultTTL    = 15 * time.Hour
	client        = http.DefaultClient
	errEmptyRes   = "got empty sting from redis"
)

const (
	parkingCategoryId    = 102
	parkingSelfId        = 621
	maxConnectionTimeout = 60 * time.Second
	allRows              = "allRows"
)

type MosService struct {
	logger  *zap.SugaredLogger
	apiKey  string
	rows    []byte
	rc      *storage.RedisClient
	client  *http.Client
	selfUrl string
}

func NewMosService(cfg *config.Config, rc *storage.RedisClient, logger *zap.SugaredLogger) *MosService {
	var apiKey string
	if cfg.ApiKey == "" {
		apiKey = defaultApiKey
	} else {
		apiKey = cfg.ApiKey
	}

	rows := make([]byte, 0)

	return &MosService{
		apiKey:  apiKey,
		rc:      rc,
		rows:    rows,
		client:  client,
		selfUrl: cfg.MosServiceUrl,
		logger:  logger,
	}
}

// надо разнести, чтобы гетАлПаркингс срабатывал только при старте системы
// далее после старта выполняется этот запрос и весь датасет сохраняется в сторадже
// последующие запросы клиента идут в наш внутренний сторадж, так как на нашем сторадже будет быстрее отдавать
// будет проще сделать логику сортировки по полям из задания, методы получения которых не предусмотрены
// в изначальном апи

// the function is runned once the app starts, no more calls in future available

var defaultHeaders = map[string]string{
	"Accept":       "*/*",
	"Connection":   "keep-alive",
	"Content-type": "application/json",
}

func (s *MosService) getAllParkingsFromUpstream() error {
	ctx, cancel := context.WithTimeout(context.Background(), maxConnectionTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.selfUrl+strconv.Itoa(parkingSelfId)+"/rows", nil)
	if err != nil {
		return fmt.Errorf("service.GetAllparkings failed %w", err)
	}

	q := req.URL.Query()
	q.Add("api_key", s.apiKey)
	req.URL.RawQuery = q.Encode()

	for key, value := range defaultHeaders {
		if len(key) > 1 {
			req.Header.Add(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("status != 200 %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var rows []types.Row
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return err
	}

	json, err := json.Marshal(rows)
	if err != nil {
		return err
	}

	s.logger.Infof("marshalled %d bytes to json", len(json))

	s.rows = json

	return nil
}

// редис возвращает строку, необходимо ее кастануть к типу
func (s *MosService) GetRowsFromStorage(ctx context.Context) ([]types.Row, error) {
	s.getAllParkingsFromUpstream()
	s.saveRowsToCache(ctx)
	s.logger.Infof("total saved rows %d", len(s.rows))

	res, err := s.rc.Redis.Get(ctx, allRows).Result()
	if err != nil {
		return nil, fmt.Errorf("mosservice.GetRowsFromStorage failed due to %w", err)
	}

	if res == "" {
		s.logger.Error(errEmptyRes)
		return nil, errors.New(errEmptyRes)
	}

	var rows []types.Row
	if err := json.Unmarshal([]byte(res), &rows); err != nil {
		return nil, fmt.Errorf("mosservice.GetRowsFromStorage.Unmarshal failed due to %w", err)
	}

	return rows, nil
}

func (s *MosService) GetParkingByGlobalId(ctx context.Context, id string) (types.Parking, error) {
	return types.Parking{}, nil
}

func (s *MosService) GetParkingById(ctx context.Context, id string) (types.Parking, error) {
	return types.Parking{}, nil
}

func (s *MosService) GetByMode(ctx context.Context, mode string) (types.Parking, error) {

	return types.Parking{}, nil
}

func (s *MosService) saveRowsToCache(ctx context.Context) error {
	err := s.rc.Redis.Set(ctx, allRows, s.rows, defaultTTL).Err()
	return err
}
