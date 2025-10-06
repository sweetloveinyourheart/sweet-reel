package timeutil

import "time"

// NowRoundedForGranularity returns a time.Time rounded to microsecond resolution
// postgres only supports microsecond resolution, so round
// when inserting (otherwise tests that read/write to db may fail
// due to rounding)
func NowRoundedForGranularity() time.Time {
	return RoundForGranularity(Clock.Now().UTC())
}

// RoundForGranularity rounds any time to the granularity specified in postgres.
func RoundForGranularity(t time.Time) time.Time {
	return t.UTC().Round(time.Microsecond)
}
