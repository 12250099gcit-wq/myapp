// myapp/utils/date/date.go
package date

import "time"

// apiDateLayout is the ISO 8601 format used for date strings in API responses.
const apiDateLayout = "2006-01-02T15:04:05Z"

// GetDate returns the current UTC date and time as a formatted string.
// Go uses a reference time (Jan 2, 2006 at 15:04:05 UTC) to define formats.
func GetDate() string {
    d := time.Now().UTC()
    return d.Format(apiDateLayout)
}