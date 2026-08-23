package rule

import "time"

func WindowBounds(at time.Time, seconds int) (time.Time, time.Time) {
	if at.IsZero() {
		at = time.Now().UTC()
	}
	if seconds <= 0 {
		seconds = 60
	}
	return at.Add(-time.Duration(seconds) * time.Second), at
}
func InWindow(at, start, end time.Time) bool { return !at.Before(start) && !at.After(end) }
func IsLate(at, watermark time.Time, grace time.Duration) bool {
	return at.Before(watermark.Add(-grace))
}
