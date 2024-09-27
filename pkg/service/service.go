package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	Logger *zap.SugaredLogger
	apiKey string
	rows   []byte
	rc     *storage.RedisClient
	client *http.Client
	mosUrl *url.URL
}

func NewMosService(cfg *config.Config, rc *storage.RedisClient, logger *zap.SugaredLogger) *MosService {
	var u *url.URL
	u, err := url.Parse(cfg.MosServiceUrl)
	if err != nil {
		panic(err)
	}

	rows := make([]byte, 0)

	return &MosService{
		apiKey: cfg.ApiKey,
		rc:     rc,
		rows:   rows,
		client: client,
		mosUrl: u,
		Logger: logger,
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

	req, err := s.makeRequest(ctx, http.MethodGet, s.mosUrl.String()+strconv.Itoa(parkingSelfId)+"/rows", nil)
	if err != nil {
		return err
	}
	s.Logger.Infof("request sucessfully composed %v", req)

	status, resp, err := s.doRequest(req)
	if err != nil {
		return err
	}
	s.Logger.Infof("done request: status:%d\n, resp:%s", status, string(resp))

	if status != 200 {
		return fmt.Errorf("status != 200\nactual:%d", status)
	}

	var parkings []types.Parking
	if err := json.Unmarshal(resp, &parkings); err != nil {
		return err
	}

	s.rows, err = json.Marshal(parkings)
	if err != nil {
		return err
	}
	s.Logger.Infof("got:%d bytes from response\n", len(s.rows))

	return nil
}

func (s *MosService) GetParkingsFromStorage(ctx context.Context) ([]types.Parking, error) {
	res, err := s.rc.Redis.Get(ctx, allRows).Result()
	if err != nil {
		return nil, fmt.Errorf("mosservice.GetRowsFromStorage failed due to %w", err)
	}

	if res == "" {
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

func (s *MosService) SaveRowsToCache(ctx context.Context) (int, error) {
	res, err := s.rc.Redis.Set(ctx, allRows, s.rows, defaultTTL).Result()
	if err != nil {
		return 0, err
	}
	return len([]byte(res)), nil
}

func (s *MosService) doRequest(req *http.Request) (int, []byte, error) {
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("mosservice.doRequest failed due to %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("mosservice.doRequest failed due to %w", err)
	}

	return res.StatusCode, body, nil
}

func (s *MosService) makeRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("mosservice.makeRequest.NewRequestWithContext failed due to :%w", err)
	}

	s.applyHeaders(defaultHeaders, req)

	q := req.URL.Query()
	q.Add("api_key", s.apiKey)
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func (s *MosService) applyHeaders(headers map[string]string, r *http.Request) {
	for key, value := range headers {
		if len(key) == 0 {
			return
		}
		r.Header.Add(key, value)
	}
}
