package mos

import (
	"context"
	"encoding/json"
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
	defaultTTL = 15 * time.Hour
)

const (
	parkingCategoryId = 102
	allRows           = "allRows"
)

type Mos struct {
	mu     sync.Mutex
	client *mosclient.MosClient
	logger *zap.SugaredLogger
	rows   []byte
	rc     *redis.RC
}

func NewMos(
	cfg *config.Config,
	rc *redis.RC,
	logger *zap.SugaredLogger,
	client *mosclient.MosClient,
) *Mos {
	rows := make([]byte, 0)

	return &Mos{
		client: client,
		rc:     rc,
		rows:   rows,
		logger: logger,
	}
}

func (s *Mos) GetParkingsFromStorage(ctx context.Context) ([]types.Parking, error) {
	data, err := s.rc.Redis.Get(ctx, allRows).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get parkings from Redis: %w", err)
	}
	s.logger.Infof("got %d bytes from store", len([]byte(data)))
	if data == "" {
		return nil, fmt.Errorf("got empty data from redis")
	}

	var parkings []types.Parking
	if err := json.Unmarshal([]byte(data), &parkings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parkings from Redis: %w", err)
	}

	return parkings, nil
}

func (s *Mos) GetParkingByGlobalId(ctx context.Context, id string) (types.Parking, error) {
	return types.Parking{}, nil
}

func (s *Mos) GetParkingById(ctx context.Context, id string) (types.Parking, error) {
	return types.Parking{}, nil
}

func (s *Mos) GetByMode(ctx context.Context, mode string) (types.Parking, error) {
	return types.Parking{}, nil
}

func (s *Mos) SaveRowsToCache(ctx context.Context) error {
	s.mu.Lock()
	if err := s.addRows(); err != nil {
		return err
	}

	if err := s.rc.Redis.Set(ctx, allRows, s.rows, defaultTTL).Err(); err != nil {
		return err
	}
	s.mu.Unlock()
	return nil
}

func (s *Mos) addRows() error {
	rows, err := s.client.GetAllParkingsFromUpstream()
	if err != nil {
		return err
	}
	s.rows = rows
	return nil
}
