package main

import (
	"fmt"
	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/metrics"
	"github.com/KillReall666/yaproject/internal/service"
	"github.com/KillReall666/yaproject/internal/storage"
	"time"
)

func main() {
	cfg := config.LoadAgentConfig()
	gms := metrics.NewGaugeMetricsStorage()
	store := storage.NewMemStorage()
	client := service.NewService(store, gms)

	tickUpdater := time.NewTicker(2 * time.Second)
	tickSender := time.NewTicker(10 * time.Second)

	defer tickUpdater.Stop()
	defer tickSender.Stop()

	for {
		select {
		case <-tickUpdater.C:
			gms.UpdateMetrics()
			gms.Gauge["PollCount"]++
			for key := range gms.Gauge {
				gms.GaugeStorage[key] = fmt.Sprintf("%f", gms.Gauge[key])
			}
			fmt.Println("Metrics update...")
			tickUpdater.Reset(2 * time.Second)

		case <-tickSender.C:
			go client.MetricsSender(&cfg)
			fmt.Println(gms.GaugeStorage)
			tickSender.Reset(10 * time.Second)
		}
	}

}
