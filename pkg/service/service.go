package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tibeahx/mos.ru-adapter/internal/config"
	storage "github.com/tibeahx/mos.ru-adapter/internal/store"
	"github.com/tibeahx/mos.ru-adapter/internal/types"

	"go.uber.org/zap"
)

var (
	defaultTTL  = 15 * time.Hour
	client      = http.DefaultClient
	errEmptyRes = "got empty sting from redis"
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
	rows := make([]byte, 0)

	return &MosService{
		apiKey:  cfg.ApiKey,
		rc:      rc,
		rows:    rows,
		client:  client,
		selfUrl: cfg.MosServiceUrl,
		logger:  logger,
	}
}

var defaultHeaders = map[string]string{
	"Accept":       "*/*",
	"Connection":   "keep-alive",
	"Content-type": "application/json",
}

func (s *MosService) GetAllParkingsFromUpstream() error {
	ctx, cancel := context.WithTimeout(context.Background(), maxConnectionTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.selfUrl+strconv.Itoa(parkingSelfId)+"/rows", nil)
	if err != nil {
		return fmt.Errorf("service.GetAllparkings failed %w", err)
	}

	q := req.URL.Query()
	q.Add("api_key", s.apiKey)
	req.URL.RawQuery = q.Encode()

	s.applyHeaders(defaultHeaders, req)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("status != 200 \ncurrent status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var parkings []types.Parking
	if err := json.NewDecoder(resp.Body).Decode(&parkings); err != nil {
		return err
	}

	json, err := json.Marshal(parkings)
	if err != nil {
		return err
	}

	s.logger.Infof("marshalled %d bytes to json", len(json))

	s.rows = json

	return nil
}

// редис возвращает строку, необходимо ее кастануть к типу
func (s *MosService) GetParkingsFromStorage(ctx context.Context) ([]types.Parking, error) {
	res, err := s.rc.Redis.Get(ctx, allRows).Result()
	if err != nil {
		return nil, fmt.Errorf("mosservice.GetRowsFromStorage failed due to %w", err)
	}

	if res == "" {
		s.logger.Error(errEmptyRes)
		return nil, errors.New(errEmptyRes)
	}

	var parkings []types.Parking
	if err := json.Unmarshal([]byte(res), &parkings); err != nil {
		return nil, fmt.Errorf("mosservice.GetRowsFromStorage.Unmarshal failed due to %w", err)
	}

	return parkings, nil
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

func (s *MosService) SaveRowsToCache(ctx context.Context) error {
	err := s.rc.Redis.Set(ctx, allRows, s.rows, defaultTTL).Err()
	return err
}

func (s *MosService) applyHeaders(headers map[string]string, r *http.Request) {
	for key, value := range headers {
		if len(key) == 0 {
			return
		}
		r.Header.Add(key, value)
	}
}
