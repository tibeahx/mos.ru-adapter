package mosclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tibeahx/mos.ru-adapter/internal/config"
	"go.uber.org/zap"
)

var httpClient = http.DefaultClient

const (
	parkingSelfId        = 621
	maxConnectionTimeout = 15 * time.Second
)

type MosClient struct {
	c      *http.Client
	apiKey string
	url    *url.URL
	logger *zap.SugaredLogger
}

func NewMosClient(cfg *config.Config, logger *zap.SugaredLogger) *MosClient {
	var u *url.URL
	u, err := url.Parse(cfg.MosUrl)
	if err != nil {
		panic(err)
	}

	return &MosClient{
		c:      httpClient,
		apiKey: cfg.ApiKey,
		url:    u,
	}
}

func (mc *MosClient) GetAllParkingsFromUpstream() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), maxConnectionTimeout)
	defer cancel()

	req, err := mc.makeRequest(ctx, http.MethodGet, mc.url.String()+strconv.Itoa(parkingSelfId)+"/rows", nil)
	if err != nil {
		return nil, err
	}

	status, resp, err := mc.doRequest(req)
	if err != nil {
		return nil, err
	}

	if err := mc.checkResponseStatus(status); err != nil {
		return nil, err
	}

	return resp, nil
}

func (mc *MosClient) doRequest(req *http.Request) (int, []byte, error) {
	res, err := mc.c.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("mc.doRequest failed due to %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("mc.doRequest failed due to %w", err)
	}

	return res.StatusCode, body, nil
}

func (mc *MosClient) makeRequest(
	ctx context.Context,
	method string,
	url string,
	body io.Reader,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("mc.makeRequest.NewRequestWithContext failed due to :%w", err)
	}

	mc.applyHeaders(mc.defaultHeaders(), req)
	mc.setApiKey(req)

	return req, nil
}

func (mc *MosClient) defaultHeaders() map[string]string {
	return map[string]string{
		"Accept":       "*/*",
		"Connection":   "keep-alive",
		"Content-type": "application/json",
	}
}

func (mc *MosClient) setApiKey(req *http.Request) {
	q := req.URL.Query()
	q.Add("api_key", mc.apiKey)
	req.URL.RawQuery = q.Encode()
}

func (mc *MosClient) checkResponseStatus(status int) error {
	switch status {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return fmt.Errorf("bad request from mos.ru")
	case http.StatusInternalServerError:
		return fmt.Errorf("internal error from mos.ru")
	case http.StatusNoContent:
		return nil
	case http.StatusServiceUnavailable:
		return fmt.Errorf("svc unavailable error from mos.ru")
	}
	return nil
}

func (mc *MosClient) applyHeaders(headers map[string]string, r *http.Request) {
	for key, value := range headers {
		if len(key) == 0 || len(value) == 0 {
			mc.logger.Panic("broken headers...")
			return
		}
		r.Header.Add(key, value)
	}
}
