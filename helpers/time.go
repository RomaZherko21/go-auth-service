package helpers

import "time"

func GetCurrentTime() time.Time {
	currentTime := time.Now().UTC()

	// return currentTime.Format(time.RFC3339)
	return currentTime
}
