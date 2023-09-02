package main

import "github.com/KillReall666/yaproject/internal/service"

func main() {
	err := service.Run()
	if err != nil {
		panic(err)
	}
}
