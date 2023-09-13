package metrics

import (
	"fmt"
	"github.com/KillReall666/yaproject/internal/model"
	"net/http"
	"time"
)

func MetricSender(gs *GaugeMetricsGetter, cfg *model.RunConfig) {
	tickUpdater := time.NewTicker(time.Duration(cfg.DefaultPollInterval))
	tickSender := time.NewTicker(time.Duration(cfg.DefaultReportInterval))
	for {
		select {
		case <-tickUpdater.C:
			gs.UpdateMetrics()
			gs.Gauge["PollCount"]++
			for key := range gs.Gauge {
				gs.GaugeStorage[key] = fmt.Sprintf("%f", gs.Gauge[key])
			}
			fmt.Println("Metrics update...")
			tickUpdater.Reset(2 * time.Second)

		case <-tickSender.C:
			for key, value := range gs.GaugeStorage {
				if gs.GaugeStorage[key] != "PollCount" {
					url := "http://" + cfg.Address + "/update/gauge/" + key + "/" + value
					resp, err := http.Post(url, "text/plain", nil)
					if err != nil {
						fmt.Println(err)
						continue
					}
					defer resp.Body.Close()
				}
			}

			url := "http://" + cfg.Address + "/update/counter/PollCount/" + gs.GaugeStorage["PollCount"]
			resp, err := http.Post(url, "text/plain", nil)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()
			fmt.Println("Metrics sent...")
			fmt.Println(gs.GaugeStorage)
			tickSender.Reset(10 * time.Second)
		}
	}
}
