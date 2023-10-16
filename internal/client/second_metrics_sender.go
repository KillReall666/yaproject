package client

import (
	"context"
	"encoding/json"
	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/handlers"
	"github.com/KillReall666/yaproject/internal/model"
	"io"
	"net/http"
	"time"
)

func (c *Client) MetricsSender(cfg *config.RunConfig) error {
	for key, value := range c.gms.Gauge {
		metric := model.MetricsJSON{
			ID:    key,
			MType: "gauge",
			Value: handlers.Float64Ptr(value),
		}
		data, err1 := json.Marshal(metric)
		if err1 != nil {
			c.logger.LogInfo("ошибка при marshal gauge:", err1)
		}

		compressedData := c.Compress(data)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1500)
		defer cancel()

		url := "http://" + cfg.Address + "/update/"
		req, err2 := http.NewRequestWithContext(ctx, "POST", url, compressedData)
		if err2 != nil {
			c.logger.LogInfo("ошибка при запросе gauge", err2)
		}

		req.Header.Set("Content-Encoding", "gzip")

		client := http.Client{}
		resp, err3 := client.Do(req)
		if err3 != nil {
			c.logger.LogInfo("ошибка при получении ответа gauge:", err3)
			return err3
		}
		defer resp.Body.Close()

		_, err := io.ReadAll(resp.Body)
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

		data, err4 := json.Marshal(metric)
		if err4 != nil {
			c.logger.LogInfo("ошибка при marshal counter", err4)
		}

		compressedData := c.Compress(data)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1500)
		defer cancel()

		url := "http://" + cfg.Address + "/update/"
		req, err5 := http.NewRequestWithContext(ctx, "POST", url, compressedData)
		if err5 != nil {
			c.logger.LogInfo("ошибка при выполнении запроса counter", err5)
		}

		req.Header.Set("Content-Encoding", "gzip")

		client := http.Client{}
		resp, err6 := client.Do(req)

		if err6 != nil {
			c.logger.LogInfo("ошибка при получении ответа counter:", err6)
		}
		defer resp.Body.Close()
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			c.logger.LogInfo(err)
		}
	}

	return nil
}
