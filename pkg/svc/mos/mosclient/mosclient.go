package mosclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tibeahx/mos.ru-adapter/internal/config"
	"github.com/tibeahx/mos.ru-adapter/internal/types"
	"go.uber.org/zap"
)

var httpClient = http.DefaultClient

const (
	parkingSelfId        = 621
	maxConnectionTimeout = 5 * time.Second
)

type MosClient struct {
	c      *http.Client
	apiKey string
	url    *url.URL
	logger *zap.SugaredLogger
}

func NewMosClient(cfg *config.Config, logger *zap.SugaredLogger) *MosClient {
	var u *url.URL
	u, err := url.Parse(cfg.MossvcUrl)
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

	resp, err := mc.retry(ctx, http.MethodGet, mc.url.String()+strconv.Itoa(parkingSelfId)+"/rows", nil)
	if err != nil {
		return nil, err
	}

	var parkings []types.Parking
	if err := json.Unmarshal(resp, &parkings); err != nil {
		return nil, err
	}

	rows, err := json.Marshal(parkings)
	if err != nil {
		return nil, err
	}
	mc.logger.Info(string(rows))

	return rows, nil
}

func (mc *MosClient) retry(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	attempt := 1

	// если не успели за 5 сек достать из апстрима простыню, идем нахуй
	// если пошли, долбимся пока не отдаст
	for {
		mc.logger.Infof("%d req attempt to mos.ru...", attempt)

		req, err := mc.makeRequest(ctx, method, url, body)
		if err != nil {
			return nil, err
		}

		status, resp, err := mc.doRequest(req)
		if err != nil {
			mc.logger.Errorf("req failed due to: %w", err)
		}

		mc.checkResponseStatus(status, resp)

		if ctx.Err() == context.Canceled {
			mc.logger.Warnf("context canceled, mos.ru req timed out... retrying")
			attempt++
			time.Sleep(time.Second)
			continue
		}

		if status == http.StatusOK {
			mc.logger.Infof("successfully got response from %d attempt", attempt)
			return resp, nil
		}
	}
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

func (mc *MosClient) makeRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("mc.makeRequest.NewRequestWithContext failed due to :%w", err)
	}

	headers := make(map[string]string, 3)
	headers["Accept"] = "*/*"
	headers["Connection"] = "keep-alive"
	headers["Content-type"] = "application/json"

	mc.applyHeaders(headers, req)

	q := req.URL.Query()
	q.Add("api_key", mc.apiKey)
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func (mc *MosClient) checkResponseStatus(status int, resp []byte) error {
	switch status {
	case http.StatusOK:
		fallthrough
	case http.StatusBadRequest:
		return fmt.Errorf("bad request from mos.ru: %v", resp)
	case http.StatusInternalServerError:
		return fmt.Errorf("internal error from mos.ru: %v", resp)
	case http.StatusNoContent:
		fallthrough
	case http.StatusServiceUnavailable:
		return fmt.Errorf("svc unavailable error from mos.ru: %v", resp)
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
