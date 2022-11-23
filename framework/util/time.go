package util

import "time"

const (
	TimeFormatterDefault      = ""
	TimeFormatterSeconds      = "2006-01-02 15:04:05"
	TimeFormatterMicroSeconds = "2006-01-02 15:04:05.000000"
	TimeFormatterNanoSeconds  = "2006-01-02 15:04:05.000000000"
)

func CurrentTimeFormatStr(formatter string) string {
	return DefaultTimeFormat(time.Now(), formatter)
}

func DefaultTimeFormatNormal(times time.Time) string {
	return DefaultTimeFormat(times, TimeFormatterDefault)
}

func DefaultTimeFormatSeconds(times time.Time) string {
	return DefaultTimeFormat(times, TimeFormatterSeconds)
}

func DefaultTimeFormatMicroSeconds(times time.Time) string {
	return DefaultTimeFormat(times, TimeFormatterMicroSeconds)
}

func DefaultTimeFormatMicroNanoSeconds(times time.Time) string {
	return DefaultTimeFormat(times, TimeFormatterNanoSeconds)
}

func DefaultTimeFormat(times time.Time, formatter string) string {
	if nowStr := times.Format(formatter); nowStr != "" {
		return nowStr
	}
	return times.Format(time.RFC3339)
}
