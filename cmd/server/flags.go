package main

import (
	"flag"
	"fmt"
	"os"
)

var addr string

func parseFlags() {
	flag.StringVar(&addr, "a", ":8080", "[addr:port] set server address")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Println("Error: unknown flags detected")
		flag.PrintDefaults()
		os.Exit(1)
	}
}
