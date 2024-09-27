package mos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/tibeahx/mos.ru-adapter/internal/config"
	"github.com/tibeahx/mos.ru-adapter/pkg/svc/mos/mosclient"
	"github.com/tibeahx/mos.ru-adapter/pkg/svc/redis"

	"github.com/tibeahx/mos.ru-adapter/internal/types"

	"go.uber.org/zap"
)

var (
	defaultTTL  = 15 * time.Hour
	errEmptyRes = "got empty sting from redis"
)

const (
	parkingCategoryId = 102
	allRows           = "allRows"
)

type Mossvc struct {
	mu     sync.Mutex
	client *mosclient.MosClient
	logger *zap.SugaredLogger
	rows   []byte
	rc     *redis.RC
}

func NewMossvc(cfg *config.Config, rc *redis.RC, logger *zap.SugaredLogger, client *mosclient.MosClient) *Mossvc {
	rows := make([]byte, 0)

	return &Mossvc{
		client: client,
		rc:     rc,
		rows:   rows,
		logger: logger,
	}
}

func (s *Mossvc) GetParkingsFromStorage(ctx context.Context) ([]types.Parking, error) {
	res, err := s.rc.Redis.Get(ctx, allRows).Result()
	if err != nil {
		return nil, fmt.Errorf("mossvc.GetRowsFromStorage failed due to %w", err)
	}

	if res == "" {
		return nil, errors.New(errEmptyRes)
	}

	var parkings []types.Parking
	if err := json.Unmarshal([]byte(res), &parkings); err != nil {
		return nil, fmt.Errorf("mossvc.GetRowsFromStorage.Unmarshal failed due to %w", err)
	}

	return parkings, nil
}

func (s *Mossvc) GetParkingByGlobalId(ctx context.Context, id string) (types.Parking, error) {
	return types.Parking{}, nil
}

func (s *Mossvc) GetParkingById(ctx context.Context, id string) (types.Parking, error) {
	return types.Parking{}, nil
}

func (s *Mossvc) GetByMode(ctx context.Context, mode string) (types.Parking, error) {
	return types.Parking{}, nil
}

func (s *Mossvc) SaveRowsToCache(ctx context.Context) error {
	s.mu.Lock()
	err := s.rc.Redis.Set(ctx, allRows, s.rows, defaultTTL).Err()
	s.mu.Unlock()
	return err
}

func (s *Mossvc) AddRows(parkings []types.Parking) error {
	rows, err := s.client.GetAllParkingsFromUpstream()
	if err != nil {
		return err
	}
	s.rows = rows
	return nil
}
