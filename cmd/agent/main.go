package main

import (
	"context"
	"fmt"
	"github.com/KillReall666/yaproject/internal/handlers/metrics"
	"time"
)

func main() {

	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)
	gaugeStorage := metrics.NewGaugeMetricsStorage()
	gaugeStorage.Gauge["PollCount"] = 0
	err := gaugeStorage.ProcessUpdating(ctx)
	if err != nil {
		fmt.Println(err)
	}

}
