package timestamp

import (
	"time"
)

var zeroTime = time.Unix(0, 0)

func DayEnd(t time.Time) time.Time {
	if t.Unix() <= 0 {
		return zeroTime
	}
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, 0, t.Location())
}

func DayStart(t time.Time) time.Time {
	if t.Unix() <= 0 {
		return zeroTime
	}
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
