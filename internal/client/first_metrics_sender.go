package client

import (
	"net/http"

	"github.com/KillReall666/yaproject/internal/config"
)

func (c *Client) MetricsSenderOld(cfg *config.RunConfig) {
	for key, value := range c.gms.GaugeStorage {
		if key == "PollCount" {
			url := "http://" + cfg.Address + "/update/counter/PollCount/" + c.gms.GaugeStorage["PollCount"]
			resp, err := http.Post(url, "text/plain", nil)
			if err != nil {
				c.logger.LogInfo("ошибка при отправке запроса counter:", err)
			}
			defer resp.Body.Close()
		} else {
			url := "http://" + cfg.Address + "/update/gauge/" + key + "/" + value
			resp, err := http.Post(url, "text/plain", nil)
			if err != nil {
				c.logger.LogInfo("ошибка при отправке запроса gauge:", err)
				continue
			}
			defer resp.Body.Close()
		}
	}
}
