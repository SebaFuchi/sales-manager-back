package timeHelper

import "time"

func Now() *string {
	nowStr := time.Now().Format(time.RFC3339)
	return &nowStr
}

func Now3Months() *string {
	nowStr := time.Now().Add(2160 * time.Hour).Format(time.RFC3339)
	return &nowStr
}
