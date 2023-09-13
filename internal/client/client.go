package client

import (
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
			c.gms.Gauge["PollCount"]++
			for key := range c.gms.Gauge {
				c.gms.GaugeStorage[key] = fmt.Sprintf("%f", c.gms.Gauge[key])
			}
			fmt.Println("Metrics update...")
			tickUpdater.Reset(2 * time.Second)

		case <-tickSender.C:
			c.MetricsSender(&c.cfg)
			fmt.Println(c.gms.GaugeStorage)
			tickSender.Reset(10 * time.Second)
		}
	}
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
