package main

import (
	agent "github.com/KillReall666/yaproject/internal/client"
	"github.com/KillReall666/yaproject/internal/client/metrics"
	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/logger"
)

func main() {
	log, err := logger.InitLogger()
	if err != nil {
		panic("cannot initialize zap")
	}

	cfg := config.LoadForAgent()
	gms := metrics.NewGaugeMetricsStorage()

	cli := agent.NewClient(cfg, gms, log)

	err = cli.Run()
	if err != nil {
		log.LogInfo("client died on error: %v", err)
	}

}
