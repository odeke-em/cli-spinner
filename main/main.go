package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/odeke-em/cli-spinner"
)

func main() {
	spinChan := spinner.Spin(1)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	go func() {
		_ = <-sigChan
		spinChan <- nil
		os.Exit(-1)
	}()

	throttle := time.Tick(1e9 / 1)
	<-throttle
	spinChan <- nil
	<-throttle
}
