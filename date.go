package sprig

import (
	"strconv"
	"time"
)

// Given a format and a date, format the date string.
//
// Date can be a `time.Time` or an `int, int32, int64`.
// In the later case, it is treated as seconds since UNIX
// epoch.
func date(fmt string, date interface{}) string {
	return dateInZone(fmt, date, "Local")
}

func htmlDate(date interface{}) string {
	return dateInZone("2006-01-02", date, "Local")
}

func htmlDateInZone(date interface{}, zone string) string {
	return dateInZone("2006-01-02", date, zone)
}

func dateInZone(fmt string, date interface{}, zone string) string {
	var t time.Time
	switch date := date.(type) {
	default:
		t = time.Now()
	case time.Time:
		t = date
	case *time.Time:
		t = *date
	case int64:
		t = time.Unix(date, 0)
	case int:
		t = time.Unix(int64(date), 0)
	case int32:
		t = time.Unix(int64(date), 0)
	}

	loc, err := time.LoadLocation(zone)
	if err != nil {
		loc, _ = time.LoadLocation("UTC")
	}

	return t.In(loc).Format(fmt)
}

func dateModify(fmt string, date time.Time) time.Time {
	d, err := time.ParseDuration(fmt)
	if err != nil {
		return date
	}
	return date.Add(d)
}

func mustDateModify(fmt string, date time.Time) (time.Time, error) {
	d, err := time.ParseDuration(fmt)
	if err != nil {
		return time.Time{}, err
	}
	return date.Add(d), nil
}

func dateAgo(date interface{}) string {
	var t time.Time

	switch date := date.(type) {
	default:
		t = time.Now()
	case time.Time:
		t = date
	case int64:
		t = time.Unix(date, 0)
	case int:
		t = time.Unix(int64(date), 0)
	}
	// Drop resolution to seconds
	duration := time.Since(t).Round(time.Second)
	return duration.String()
}

func duration(sec interface{}) string {
	var n int64
	switch value := sec.(type) {
	default:
		n = 0
	case string:
		n, _ = strconv.ParseInt(value, 10, 64)
	case int64:
		n = value
	}
	return (time.Duration(n) * time.Second).String()
}

const (
	year  = time.Hour * 24 * 365
	month = time.Hour * 24 * 30
	day   = time.Hour * 24
)

func durationRound(duration interface{}) string {
	var d time.Duration

	switch duration := duration.(type) {
	default:
		return "0s"
	case string:
		d, _ = time.ParseDuration(duration)
	case int:
		// We handle these cases similar to how `duration` does.
		d = time.Duration(duration) * time.Second
	case int64:
		// Considering the given value as seconds instead of nanoseconds might be a breaking
		// change, but it is more consistent with the other cases and most likely closer to what
		// the user expects.
		d = time.Duration(duration) * time.Second
	case float64:
		d = time.Duration(duration) * time.Second
	case time.Time:
		d = time.Since(duration)
	case time.Duration:
		d = duration
	}

	// Not sure if this actually makes much sense, but removing it would be a breaking change.
	if d < 0 {
		d = -d
	}

	if d > year {
		return strconv.FormatInt(int64(d/year), 10) + "y"
	}
	if d > month {
		return strconv.FormatInt(int64(d/month), 10) + "mo"
	}
	if d > day {
		return strconv.FormatInt(int64(d/day), 10) + "d"
	}
	if d > time.Hour {
		return strconv.FormatInt(int64(d/time.Hour), 10) + "h"
	}
	if d > time.Minute {
		return strconv.FormatInt(int64(d/time.Minute), 10) + "m"
	}
	if d > time.Second {
		return strconv.FormatInt(int64(d/time.Second), 10) + "s"
	}

	return "0s"
}

func toDate(fmt, str string) time.Time {
	t, _ := time.ParseInLocation(fmt, str, time.Local)
	return t
}

func mustToDate(fmt, str string) (time.Time, error) {
	return time.ParseInLocation(fmt, str, time.Local)
}

func unixEpoch(date time.Time) string {
	return strconv.FormatInt(date.Unix(), 10)
}
