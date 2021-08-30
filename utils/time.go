package utils

import "time"

func Time(number float64) time.Time {
	return time.Unix(int64(number), 0)
}
