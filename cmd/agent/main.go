package main

import (
	"github.com/KillReall666/yaproject/internal/appclient"
	agent "github.com/KillReall666/yaproject/internal/client"
	metrics2 "github.com/KillReall666/yaproject/internal/client/metrics"
	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/logger"
)

func main() {
	log, err1 := logger.InitLogger()
	if err1 != nil {
		panic("cannot initialize zap")
	}

	cfg := config.LoadAgentConfig()
	gms := metrics2.NewGaugeMetricsStorage()

	cli := agent.NewClient(cfg, gms, log)
	
	app := appclient.NewAgentService(log, cli)

	err := cli.Run()
	if err != nil {
		app.LogInfo("client died on error: %v", err)
	}

}
