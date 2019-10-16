package cloudfunc

import (
	"time"
)

// GetMinsOfWeek return mins passed from week start
func GetMinsOfWeek(t time.Time) int {
	mins := t.Minute() + t.Hour()*60 + int(t.Weekday())*24*60
	return mins
}
