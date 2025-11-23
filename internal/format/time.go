package format

import "time"

func TimeToISO8601(timestamp time.Time) string {
	return timestamp.UTC().Format("2006-01-02T15:04:05.000Z")
}

func ISO8601ToTime(timestamp string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.000Z", timestamp)
}