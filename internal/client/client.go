package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/KillReall666/yaproject/internal/client/metrics"
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
			c.PackMetricsSender(&c.cfg)
			tickSender.Reset(time.Duration(c.cfg.DefaultReportInterval) * time.Second)

		}
	}
}

func (c *Client) PackMetricsSender(cfg *config.RunConfig) {
	var packDataGauge []model.MetricsJSON

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
	compressedData := c.Compress(data)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	url := "http://" + cfg.Address + "/updates/"
	req, err := http.NewRequestWithContext(ctx, "POST", url, compressedData)
	if err != nil {
		c.logger.LogInfo("ошибка при запросе gauge", err)
	}

	req.Header.Set("Content-Encoding", "gzip")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.logger.LogInfo("ошибка при получении ответа gauge:", err)
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		c.logger.LogInfo(err)
	}
}

func (c *Client) Compress(data []byte) *bytes.Buffer {
	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)
	_, err := gzipWriter.Write(data)
	if err != nil {
		c.logger.LogInfo(err)
	}
	err = gzipWriter.Close()
	if err != nil {
		c.logger.LogInfo(err)

	}
	return &compressedData
}
