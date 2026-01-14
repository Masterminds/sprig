package sprig

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
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

func durationRound(duration interface{}) string {
	var d time.Duration
	switch duration := duration.(type) {
	default:
		d = 0
	case string:
		d, _ = time.ParseDuration(duration)
	case int64:
		d = time.Duration(duration)
	case time.Time:
		d = time.Since(duration)
	}

	u := uint64(d)
	neg := d < 0
	if neg {
		u = -u
	}

	var (
		year   = uint64(time.Hour) * 24 * 365
		month  = uint64(time.Hour) * 24 * 30
		day    = uint64(time.Hour) * 24
		hour   = uint64(time.Hour)
		minute = uint64(time.Minute)
		second = uint64(time.Second)
	)
	switch {
	case u > year:
		return strconv.FormatUint(u/year, 10) + "y"
	case u > month:
		return strconv.FormatUint(u/month, 10) + "mo"
	case u > day:
		return strconv.FormatUint(u/day, 10) + "d"
	case u > hour:
		return strconv.FormatUint(u/hour, 10) + "h"
	case u > minute:
		return strconv.FormatUint(u/minute, 10) + "m"
	case u > second:
		return strconv.FormatUint(u/second, 10) + "s"
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

// -----------------------------------------------------------------------------
// Duration helpers (numeric-only returns)
// -----------------------------------------------------------------------------

// asDuration converts common template values into a time.Duration.
//
// Supported inputs:
//   - time.Duration
//   - string duration values parsed by time.ParseDuration (e.g. "1h2m3s")
//   - numeric strings treated as seconds (e.g. "2.5")
//   - ints and uints treated as seconds
//   - floats treated as seconds
func asDuration(v interface{}) (time.Duration, error) {
	switch x := v.(type) {
	case time.Duration:
		return x, nil

	case string:
		s := strings.TrimSpace(x)
		if s == "" {
			return 0, fmt.Errorf("empty duration")
		}
		if d, err := time.ParseDuration(s); err == nil {
			return d, nil
		}
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return time.Duration(f * float64(time.Second)), nil
		}
		return 0, fmt.Errorf("could not parse duration %q", x)

	case nil:
		return 0, fmt.Errorf("invalid duration")
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return time.Duration(rv.Int()) * time.Second, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		u := rv.Uint()
		if u > uint64(math.MaxInt64) {
			return 0, fmt.Errorf("duration seconds overflow: %d", u)
		}
		return time.Duration(int64(u)) * time.Second, nil
	case reflect.Float32, reflect.Float64:
		return time.Duration(rv.Float() * float64(time.Second)), nil
	default:
		return 0, fmt.Errorf("unsupported duration type %T", v)
	}
}

// mustToDuration takes an interface, parses a duration, and returns a time.Duration.
// It will panic if there is an error.
//
// This is designed to be called from a template when need to ensure that a
// duration is valid.
func mustToDuration(v interface{}) time.Duration {
	d, err := asDuration(v)
	if err != nil {
		panic(err)
	}
	return d
}

// durationSeconds converts a duration to seconds (float64).
// On error it returns 0.
func durationSeconds(v interface{}) float64 {
	d, err := asDuration(v)
	if err != nil {
		return 0
	}
	return d.Seconds()
}

// durationMilliseconds converts a duration to milliseconds (int64).
// On error it returns 0.
func durationMilliseconds(v interface{}) int64 {
	d, err := asDuration(v)
	if err != nil {
		return 0
	}
	return d.Milliseconds()
}

// durationMicroseconds converts a duration to microseconds (int64).
func durationMicroseconds(v interface{}) int64 {
	d, err := asDuration(v)
	if err != nil {
		return 0
	}
	return d.Microseconds()
}

// durationNanoseconds converts a duration to nanoseconds (int64).
// On error it returns 0.
func durationNanoseconds(v interface{}) int64 {
	d, err := asDuration(v)
	if err != nil {
		return 0
	}
	return d.Nanoseconds()
}

// durationMinutes converts a duration to minutes (float64).
func durationMinutes(v interface{}) float64 {
	d, err := asDuration(v)
	if err != nil {
		return 0
	}
	return d.Minutes()
}

// durationHours converts a duration to hours (float64).
// On error it returns 0.
func durationHours(v interface{}) float64 {
	d, err := asDuration(v)
	if err != nil {
		return 0
	}
	return d.Hours()
}

// durationDays converts a duration to days (float64). (Not in Go's stdlib; handy in templates.)
// On error it returns 0.
func durationDays(v interface{}) float64 {
	d, err := asDuration(v)
	if err != nil {
		return 0
	}
	return d.Hours() / 24.0
}

// durationWeeks converts a duration to weeks (float64). (Not in Go's stdlib; handy in templates.)
// On error it returns 0.
func durationWeeks(v interface{}) float64 {
	d, err := asDuration(v)
	if err != nil {
		return 0
	}
	return d.Hours() / 24.0 / 7.0
}

// durationRoundTo rounds v to the nearest multiple of m.
// Returns a time.Duration.
//
// v and m accept the same forms as asDuration (e.g. "2h13m", "30s").
// On error, it returns time.Duration(0). If m is invalid, it returns v.
func durationRoundTo(v interface{}, m interface{}) time.Duration {
	d, err := asDuration(v)
	if err != nil {
		return 0
	}
	mul, err := asDuration(m)
	if err != nil {
		return d
	}
	return d.Round(mul)
}

// durationTruncateTo truncates v toward zero to a multiple of m.
// Returns a time.Duration.
//
// On error, it returns time.Duration(0). If m is invalid, it returns v.
func durationTruncateTo(v interface{}, m interface{}) time.Duration {
	d, err := asDuration(v)
	if err != nil {
		return 0
	}
	mul, err := asDuration(m)
	if err != nil {
		return d
	}
	return d.Truncate(mul)
}
