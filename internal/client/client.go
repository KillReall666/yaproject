package client

import (

	"bytes"
	"encoding/json"
	"fmt"
	"github.com/KillReall666/yaproject/internal/handlers"
	"github.com/KillReall666/yaproject/internal/model"

	"fmt"

	"net/http"
	"time"

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

			c.gms.Gauge["PollCount"]++
			for key := range c.gms.Gauge {
				c.gms.GaugeStorage[key] = fmt.Sprintf("%f", c.gms.Gauge[key])
			}
			fmt.Println("Metrics update...")

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
		data, err := json.Marshal(metric)
		if err != nil {
			return err
		}
		resp, err := http.Post("http://"+cfg.Address+"/update/", "application/json", bytes.NewBuffer(data))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
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
			return err
		}
		resp, err := http.Post("http://"+cfg.Address+"/update/", "application/json", bytes.NewBuffer(data))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
		}
	}
	return nil
}

func (c *Client) MetricsSenderOld(cfg *config.RunConfig) {
	for key, value := range c.gms.GaugeStorage {
		if key == "PollCount" {
			url := "http://" + cfg.Address + "/update/counter/PollCount/" + c.gms.GaugeStorage["PollCount"]

			fmt.Println(c.gms.GaugeStorage)
			tickSender.Reset(10 * time.Second)
		}
	}
	//return nil
}

func (c *Client) MetricsSender(cfg *config.RunConfig) {
	for key, value := range c.gms.GaugeStorage {
		switch c.gms.GaugeStorage[key] {
		case "PollCount":
			url := "http://" + cfg.Address + "/html/counter/PollCount/" + c.gms.GaugeStorage["PollCount"]

			resp, err := http.Post(url, "text/plain", nil)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()

		} else {
			url := "http://" + cfg.Address + "/update/gauge/" + key + "/" + value


		default:
			url := "http://" + cfg.Address + "/html/gauge/" + key + "/" + value

			resp, err := http.Post(url, "text/plain", nil)
			if err != nil {
				fmt.Println("error sending request:", err)
				continue
			}
			defer resp.Body.Close()
			//fmt.Println("request sent successfully:", resp.Status)

		}
	}
}
