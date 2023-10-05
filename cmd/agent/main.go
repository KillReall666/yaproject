package main

import (
	"log"

	agent "github.com/KillReall666/yaproject/internal/client"
	metrics2 "github.com/KillReall666/yaproject/internal/client/metrics"
	"github.com/KillReall666/yaproject/internal/config"
)

func main() {
	cfg := config.LoadAgentConfig()
	gms := metrics2.NewGaugeMetricsStorage()
	cli := agent.NewClient(cfg, gms)


	err := cli.Run()
	if err != nil {
		log.Fatalf("client died on error: %v", err)
	}

}
