package clock

import "time"

// передать время в gui

func Get() (msec, sec, minute, hour int) {
	now := time.Now()
	msec = now.Nanosecond() / 1000000
	sec = now.Second()
	minute = now.Minute()
	hour = now.Hour()
	return
}
