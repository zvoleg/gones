package internal

import "time"

func ClockWaiter(timeSince time.Time, period time.Duration) {
	for time.Since(timeSince) < period {
	}
}
