package readers

import "time"

// CurrentTime returns current UTC time rounded to seconds.
// This method's purpose is to standarize time formats in this project.
// Use it.
func CurrentTime() time.Time {
	return time.Now().Round(time.Second).UTC()
}
