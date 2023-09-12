package main

import (
	"context"
	"fmt"
	"github.com/KillReall666/yaproject/internal/handlers/metrics"
)

func main() {
	setEnv()

	parseFlag()
	//ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)
	ctx := context.Background()
	gaugeStorage := metrics.NewGaugeMetricsStorage()
	gaugeStorage.Gauge["PollCount"] = 0
	err := gaugeStorage.ProcessUpdating(ctx, defaultPollInterval, defaultReportInterval)
	if err != nil {
		fmt.Println(err)
	}

}
