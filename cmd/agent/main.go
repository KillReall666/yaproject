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
	gs := metrics.NewGaugeMetricsStorage()
	store := storage.NewMemStorage()
	serv := service.NewService(store, gs)

	tickUpdater := time.NewTicker(2 * time.Second)
	tickSender := time.NewTicker(10 * time.Second)

	defer tickUpdater.Stop()
	defer tickSender.Stop()

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
			go serv.MetricsSender(&cfg)
			fmt.Println(gs.GaugeStorage)
			tickSender.Reset(10 * time.Second)
		}
	}

}
