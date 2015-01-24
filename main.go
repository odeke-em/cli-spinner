package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func spin() chan interface{} {
	outChan := make(chan interface{})
	go func() {
		var spinner = map[string]string{
			" ⏳ ": "\033[00m",
			" ⌛ ": "\033[33m",
		}
		throttle := time.Tick(1e9 / 10)
		running := true
		for running {
			for segment, color := range spinner {
				select {
				case in := <-outChan:
					if in == nil {
						running = false
						break
					}
				default:
				}
				fmt.Printf("%s%s\033[00m\r", color, segment)
				<-throttle
			}
		}
	}()
	return outChan
}

func main() {
	spinChan := spin()

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
