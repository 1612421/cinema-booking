package time

import (
	"time"
)

func MonthDuration(months int) time.Duration {
	return time.Duration(months) * 30 * DayDuration(1)
}

func DayDuration(days int) time.Duration {
	return time.Duration(days) * HourDuration(24)
}

func HourDuration(hours int) time.Duration {
	return time.Duration(hours) * time.Hour
}

func MinuteDuration(minutes int) time.Duration {
	return time.Duration(minutes) * time.Minute
}

func SecondDuration(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}
