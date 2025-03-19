package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {

	(&Warp{}).Run()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-sigCh:
			return
		}
	}
}
