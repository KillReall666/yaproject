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
	"io"
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
			c.MetricsSender(&c.cfg)
			tickSender.Reset(10 * time.Second)
		}
	}
}

func (c *Client) MetricsSender(cfg *config.RunConfig) error {
	for key, value := range c.gms.Gauge {
		metric := model.MetricsJSON{
			ID:    key,
			MType: "gauge",
			Value: handlers.Float64Ptr(value),
		}

		data, err1 := json.Marshal(metric)
		if err1 != nil {
			log.Println("ошибка при marshal gauge:", err1)
		}

		compressedData := Compress(data)

		url := "http://" + cfg.Address + "/update/"
		req, err2 := http.NewRequest("POST", url, compressedData)
		if err2 != nil {
			log.Println("ошибка при запросе gauge", err2)
		}

		req.Header.Set("Content-Encoding", "gzip")

		client := http.Client{}
		resp, err3 := client.Do(req)
		if err3 != nil {
			log.Println("ошибка при получении ответа gauge:", err3)
			return err3
		}
		defer resp.Body.Close()

		_, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		//log.Println("Gauge: ", string(res))
	}

	for key, val := range c.gms.Counter {
		metric := model.MetricsJSON{
			ID:    key,
			MType: "counter",
			Delta: handlers.Int64Ptr(val),
		}

		data, err4 := json.Marshal(metric)
		if err4 != nil {
			log.Println("ошибка при marshal counter", err4)
		}

		compressedData := Compress(data)

		url := "http://" + cfg.Address + "/update/"
		req, err5 := http.NewRequest("POST", url, compressedData)
		if err5 != nil {
			log.Println("ошибка при выполнении запроса counter", err5)
		}

		req.Header.Set("Content-Encoding", "gzip")

		client := http.Client{}
		resp, err6 := client.Do(req)

		if err6 != nil {
			log.Println("ошибка при получении ответа counter:", err6)
		}
		defer resp.Body.Close()

		_, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		//log.Println("Counter: ", string(res))
	}

	return nil
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
