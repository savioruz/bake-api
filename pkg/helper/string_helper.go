package helper

import "time"

func FormatTime(date time.Time) string {
	return date.UTC().Add(time.Hour * 7).Format(time.RFC3339)
}

func StrToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
