package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/KillReall666/yaproject/internal/handlers"
	"github.com/KillReall666/yaproject/internal/model"
	"net/http"
	"time"
	"log"
	"github.com/KillReall666/yaproject/internal/client/metrics"
	"github.com/KillReall666/yaproject/internal/config"
)

type Client struct {
	cfg config.RunConfig
	gms metrics.GaugeMetricsGetter
}

func NewClient(cfg config.RunConfig, gms *metrics.GaugeMetricsGetter) *Client {
	return &Client{
		cfg: cfg,
		gms: *gms,
	}
}

func (c *Client) Run() error {
	tickUpdater := time.NewTicker(2 * time.Second)
	tickSender := time.NewTicker(10 * time.Second)
	defer tickUpdater.Stop()
	defer tickSender.Stop()
	for {
		select {
		case <-tickUpdater.C:
			c.gms.UpdateMetrics()
			c.gms.Counter["PollCount"]++
			tickUpdater.Reset(2 * time.Second)

		case <-tickSender.C:
			c.CounterMetricsSender(&c.cfg)
			c.GaugeMetricsSender(&c.cfg)
			tickSender.Reset(10 * time.Second)
		}
	}
}

func (c *Client) GaugeMetricsPrepare() *bytes.Buffer {
	for key, value := range c.gms.Gauge {
		metric := model.MetricsJSON{
			ID:    key,
			MType: "gauge",
			Value: handlers.Float64Ptr(value),
		}
		data, err := json.Marshal(metric)
		if err != nil {
			log.Println(err)
		}

		compressedData := Compress(data)
		return compressedData
	}
	return nil
}

func (c *Client) CountMetricPrepare() *bytes.Buffer {
	for key, val := range c.gms.Counter {
		metric := model.MetricsJSON{
			ID:    key,
			MType: "counter",
			Delta: handlers.Int64Ptr(val),
		}
		data, err := json.Marshal(metric)
		if err != nil {
			log.Println(err)
		}

		compressedData := Compress(data)
		return compressedData
	}
	return nil
}

func (c *Client) GaugeMetricsSender(cfg *config.RunConfig) {
	url := "http://" + cfg.Address + "/update/"
	resp, err := http.NewRequest("POST", url, c.GaugeMetricsPrepare())
	resp.Header.Set("Content-Encoding", "gzip")
	client := http.Client{}
	client.Do(resp)
	if err != nil {
		fmt.Println("ошибка при выполнении запроса:", err)
	}
	defer resp.Body.Close()
}

func (c *Client) CounterMetricsSender(cfg *config.RunConfig) {
	headers := http.Header{}
	headers.Set("Content-Encoding", "gzip")
	url := "http://" + cfg.Address + "/update/"
	resp, err := http.NewRequest("POST", url, c.CountMetricPrepare())
	resp.Header.Set("Content-Encoding", "gzip")
	client := http.Client{}
	client.Do(resp)
	if err != nil {
		fmt.Println("ошибка при выполнении запроса:", err)
	}
	defer resp.Body.Close()
}



func Compress(data []byte) *bytes.Buffer {
	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)
	gzipWriter.Write(data)
	gzipWriter.Close()
	return &compressedData
}

func (c *Client) MetricsSenderOld(cfg *config.RunConfig) {
	for key, value := range c.gms.GaugeStorage {
		if key == "PollCount" {
			url := "http://" + cfg.Address + "/update/counter/PollCount/" + c.gms.GaugeStorage["PollCount"]
			resp, err := http.Post(url, "text/plain", nil)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()
		} else {
			url := "http://" + cfg.Address + "/update/gauge/" + key + "/" + value
			resp, err := http.Post(url, "text/plain", nil)
			if err != nil {
				fmt.Println("error sending request:", err)
				continue
			}
			defer resp.Body.Close()
		}
	}
}
