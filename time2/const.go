package time2

import (
	"time"
)

const (
	// Day is a 24 hour long duration.
	Day = 24 * time.Hour
	// Week is a 7 days long duration.
	Week = 7 * Day
	// Month is a 30 days long duration.
	Month = 30 * Day
	// Year is a 365 days long duration.
	Year = 365 * Day
)
