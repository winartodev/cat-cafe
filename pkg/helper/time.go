package helper

import "time"

func TimeUTC() time.Time {
	return time.Now().UTC()
}
