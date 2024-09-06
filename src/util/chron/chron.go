package chron

import "time"

type Stopwatch struct {
	startTime time.Time
	endTime   time.Time
}

func Start() Stopwatch {
	return Stopwatch{
		startTime: time.Now(),
	}
}

func (sw *Stopwatch) Stop() time.Duration {
	sw.endTime = time.Now()
	return sw.endTime.Sub(sw.startTime)
}

func (sw *Stopwatch) Reset() {
	sw.startTime = time.Now()
	sw.endTime = time.Time{}
}
