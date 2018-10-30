package readiness

import (
	"time"
)

type Readiness interface {
	Interval() uint
	IsReady() bool
}

// Keep alive
// omit event if readiness updated
func Ready(readiness Readiness) chan bool {
	interrupt := make(chan bool)

	go func() {
		lastReady := false
		for range time.Tick(time.Second * time.Duration(readiness.Interval())) {
			isReady := readiness.IsReady()
			if lastReady == isReady {
				continue
			}

			lastReady = isReady // update last ready status

			interrupt <- isReady
		}
	}()

	return interrupt
}
