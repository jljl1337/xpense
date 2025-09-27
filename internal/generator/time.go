package generator

import (
	"time"

	"github.com/jljl1337/xpense/internal/format"
)

func NowISO8601() string {
	return format.TimeToISO8601(time.Now())
}