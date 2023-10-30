package client

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/KillReall666/yaproject/internal/client/metrics"
	"github.com/avast/retry-go"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/handlers"
	"github.com/KillReall666/yaproject/internal/logger"
	"github.com/KillReall666/yaproject/internal/model"
)

type Client struct {
	cfg    config.RunConfig
	gms    metrics.GaugeMetricsGetter
	logger *logger.Logger
}

func NewClient(cfg config.RunConfig, gms *metrics.GaugeMetricsGetter, log *logger.Logger) *Client {
	return &Client{
		cfg:    cfg,
		gms:    *gms,
		logger: log,
	}
}

func (c *Client) Run() error {
	tickUpdater := time.NewTicker(time.Duration(c.cfg.DefaultPollInterval) * time.Second)
	tickSender := time.NewTicker(time.Duration(c.cfg.DefaultReportInterval) * time.Second)
	defer tickUpdater.Stop()
	defer tickSender.Stop()
	for {
		select {
		case <-tickUpdater.C:
			c.gms.UpdateMetrics()
			c.gms.Counter["PollCount"]++
			tickUpdater.Reset(time.Duration(c.cfg.DefaultPollInterval) * time.Second)

		case <-tickSender.C:
			err := c.PackMetricsSender(&c.cfg)
			if err != nil {
				c.logger.LogInfo("ошибка при попытке отправки метрик на сервер:", err)
			}
			tickSender.Reset(time.Duration(c.cfg.DefaultReportInterval) * time.Second)

		}
	}
}

func (c *Client) PackMetricsSender(cfg *config.RunConfig) error {
	var packDataGauge []model.MetricsJSON
	var compressedData *bytes.Buffer
	var hash string

	for key, value := range c.gms.Gauge {
		metric := model.MetricsJSON{
			ID:    key,
			MType: "gauge",
			Value: handlers.Float64Ptr(value),
		}
		packDataGauge = append(packDataGauge, metric)
	}

	for key, val := range c.gms.Counter {
		metric := model.MetricsJSON{
			ID:    key,
			MType: "counter",
			Delta: handlers.Int64Ptr(val),
		}
		packDataGauge = append(packDataGauge, metric)
	}

	data, err1 := json.Marshal(packDataGauge)
	if err1 != nil {
		c.logger.LogInfo("ошибка при marshal gauge:", err1)
	}

	if c.cfg.HashKey != "" {
		hash = c.computeSHA256Hash(data)
		compressedData = c.Compress(data)
	} else {
		compressedData = c.Compress(data)
	}

	url := "http://" + cfg.Address + "/updates/"
	err := retry.Do(
		func() error {
			req, err := http.NewRequest("POST", url, compressedData)
			if err != nil {
				c.logger.LogInfo("ошибка при запросе gauge", err)
			}

			if c.cfg.HashKey != "" {
				req.Header.Set("HashSHA256", hash)
			}

			req.Header.Set("Content-Encoding", "gzip")

			client := http.Client{}

			resp, err := client.Do(req)
			if err != nil {
				var netErr net.Error
				if (errors.As(err, &netErr) && netErr.Timeout()) ||
					strings.Contains(err.Error(), "EOF") ||
					strings.Contains(err.Error(), "connection reset by peer") {
					return err
				}
				return retry.Unrecoverable(err)
			}

			defer resp.Body.Close()

			_, err = io.ReadAll(resp.Body)
			if err != nil {
				c.logger.LogInfo("error when reading response body: ", err)
			}
			return err
		},
		retry.Attempts(3),
		retry.Delay(time.Second),
		retry.DelayType(retry.BackOffDelay),
	)
	return err
}

func (c *Client) Compress(data []byte) *bytes.Buffer {
	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)
	_, err := gzipWriter.Write(data)
	if err != nil {
		c.logger.LogInfo("compression error: ", err)
	}
	err = gzipWriter.Close()
	if err != nil {
		c.logger.LogInfo(err)
	}
	return &compressedData
}

func (c *Client) computeSHA256Hash(data []byte) string {
	if c.cfg.HashKey == "" {
		return ""
	}
	hash := hmac.New(sha256.New, []byte(c.cfg.HashKey))
	hash.Write(data)
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
