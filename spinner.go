package spinner

import (
	"fmt"
	"time"
)

func Spin(freq int64) chan interface{} {
	outChan := make(chan interface{})
	go func() {
		var spinner = map[string]string{
			" ⏳ ": "\033[00m",
			" ⌛ ": "\033[33m",
		}
		if freq < 1 {
			freq = 10
		}
		throttle := time.Tick(time.Duration(1e9 / freq))
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
