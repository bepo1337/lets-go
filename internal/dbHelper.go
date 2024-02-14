package internal

import "fmt"

const (
	DAY   string = "DAY"
	MONTH        = "MONTH"
)

// DateAddString Test
func DateAddString(amount int, size string) string {
	return fmt.Sprintf("DATE_ADD(UTC_TIMESTAMP(), INTERVAL %d %s)", amount, size)
}
