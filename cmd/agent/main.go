package main

import (
	"fmt"
	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/metrics"
)

func main() {
	cfg := config.LoadAgentConfig()
	fmt.Println(cfg)
	gs := metrics.NewGaugeMetricsStorage()

	metrics.MetricSender(gs, &cfg)

}
