package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	defaultServer         = ":8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
)

var (
	address string
)

func parseFlag() {
	flag.Int("p", defaultPollInterval, "metrics update interval in seconds")
	flag.Int("r", defaultReportInterval, "metrics sending interval in seconds")
	flag.StringVar(&address, "a", defaultServer, "server address [host:port]")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Println("Error: unknown flags detected")
		flag.PrintDefaults()
		os.Exit(1)
	}

}
