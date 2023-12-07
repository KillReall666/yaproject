package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/handlers"
	"github.com/KillReall666/yaproject/internal/model"
)

func (c *Client) MetricsSender(cfg *config.RunConfig) error {
	for key, value := range c.gms.Gauge {
		metric := model.MetricsJSON{
			ID:    key,
			MType: "gauge",
			Value: handlers.Float64Ptr(value),
		}
		data, err := json.Marshal(metric)
		if err != nil {
			c.logger.LogInfo("ошибка при marshal gauge:", err)
		}

		compressedData := c.Compress(data)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1500)
		defer cancel()

		url := "http://" + cfg.Address + "/update/"
		req, err := http.NewRequestWithContext(ctx, "POST", url, compressedData)
		if err != nil {
			c.logger.LogInfo("ошибка при запросе gauge", err)
		}

		req.Header.Set("Content-Encoding", "gzip")

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.logger.LogInfo("ошибка при получении ответа gauge:", err)
			return err
		}
		defer resp.Body.Close()

		_, err = io.ReadAll(resp.Body)
		if err != nil {
			c.logger.LogInfo(err)
		}
	}

	for key, val := range c.gms.Counter {
		metric := model.MetricsJSON{
			ID:    key,
			MType: "counter",
			Delta: handlers.Int64Ptr(val),
		}

		data, err := json.Marshal(metric)
		if err != nil {
			c.logger.LogInfo("ошибка при marshal counter", err)
		}

		compressedData := c.Compress(data)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1500)
		defer cancel()

		url := "http://" + cfg.Address + "/update/"
		req, err := http.NewRequestWithContext(ctx, "POST", url, compressedData)
		if err != nil {
			c.logger.LogInfo("ошибка при выполнении запроса counter", err)
		}

		req.Header.Set("Content-Encoding", "gzip")

		client := http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			c.logger.LogInfo("ошибка при получении ответа counter:", err)
		}
		defer resp.Body.Close()
		_, err = io.ReadAll(resp.Body)
		if err != nil {
			c.logger.LogInfo(err)
		}
	}

	return nil
}
