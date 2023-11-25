package cpu6502

import (
	"time"
)

const freequencee float64 = 1789773.0
const clockTime float64 = 1.0 / freequencee

type clock struct {
	clocks uint
}

func (c *clock) waitExecution() {
	executionTimeSec := float64(c.clocks) * clockTime
	executionTimeNs := time.Duration(executionTimeSec * 1000000000)
	time.Sleep(executionTimeNs)
}
