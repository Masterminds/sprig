package sprig

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHtmlDate(t *testing.T) {
	tm, err := time.Parse("02 Jan 06 15:04:05 MST", "13 Jun 19 20:39:39 GMT")
	if err != nil {
		t.Error(err)
	}
	tpl := `{{ .Time | htmlDate }}`
	if err := runtv(tpl, "2019-06-13", map[string]interface{}{"Time": tm}); err != nil {
		t.Error(err)
	}
}

func TestAgo(t *testing.T) {
	tpl := "{{ ago .Time }}"
	if err := runtv(tpl, "2m5s", map[string]interface{}{"Time": time.Now().Add(-125 * time.Second)}); err != nil {
		t.Error(err)
	}

	if err := runtv(tpl, "2h34m17s", map[string]interface{}{"Time": time.Now().Add(-(2*3600 + 34*60 + 17) * time.Second)}); err != nil {
		t.Error(err)
	}

	if err := runtv(tpl, "-5s", map[string]interface{}{"Time": time.Now().Add(5 * time.Second)}); err != nil {
		t.Error(err)
	}
}

func TestToDate(t *testing.T) {
	tpl := `{{toDate "2006-01-02" "2017-12-31" | date "02/01/2006"}}`
	if err := runt(tpl, "31/12/2017"); err != nil {
		t.Error(err)
	}
}

func TestUnixEpoch(t *testing.T) {
	tm, err := time.Parse("02 Jan 06 15:04:05 MST", "13 Jun 19 20:39:39 GMT")
	if err != nil {
		t.Error(err)
	}
	tpl := `{{unixEpoch .Time}}`

	if err = runtv(tpl, "1560458379", map[string]interface{}{"Time": tm}); err != nil {
		t.Error(err)
	}
}

func TestDateInZone(t *testing.T) {
	tm, err := time.Parse("02 Jan 06 15:04:05 MST", "13 Jun 19 20:39:39 GMT")
	if err != nil {
		t.Error(err)
	}
	tpl := `{{ date_in_zone "02 Jan 06 15:04 -0700" .Time "UTC" }}`

	// Test time.Time input
	if err = runtv(tpl, "13 Jun 19 20:39 +0000", map[string]interface{}{"Time": tm}); err != nil {
		t.Error(err)
	}

	// Test pointer to time.Time input
	if err = runtv(tpl, "13 Jun 19 20:39 +0000", map[string]interface{}{"Time": &tm}); err != nil {
		t.Error(err)
	}

	// Test no time input. This should be close enough to time.Now() we can test
	loc, _ := time.LoadLocation("UTC")
	if err = runtv(tpl, time.Now().In(loc).Format("02 Jan 06 15:04 -0700"), map[string]interface{}{"Time": ""}); err != nil {
		t.Error(err)
	}

	// Test unix timestamp as int64
	if err = runtv(tpl, "13 Jun 19 20:39 +0000", map[string]interface{}{"Time": int64(1560458379)}); err != nil {
		t.Error(err)
	}

	// Test unix timestamp as int32
	if err = runtv(tpl, "13 Jun 19 20:39 +0000", map[string]interface{}{"Time": int32(1560458379)}); err != nil {
		t.Error(err)
	}

	// Test unix timestamp as int
	if err = runtv(tpl, "13 Jun 19 20:39 +0000", map[string]interface{}{"Time": int(1560458379)}); err != nil {
		t.Error(err)
	}

	// Test case of invalid timezone
	tpl = `{{ date_in_zone "02 Jan 06 15:04 -0700" .Time "foobar" }}`
	if err = runtv(tpl, "13 Jun 19 20:39 +0000", map[string]interface{}{"Time": tm}); err != nil {
		t.Error(err)
	}
}

func TestDuration(t *testing.T) {
	tpl := "{{ duration .Secs }}"
	if err := runtv(tpl, "1m1s", map[string]interface{}{"Secs": "61"}); err != nil {
		t.Error(err)
	}
	if err := runtv(tpl, "1h0m0s", map[string]interface{}{"Secs": "3600"}); err != nil {
		t.Error(err)
	}
	// 1d2h3m4s but go is opinionated
	if err := runtv(tpl, "26h3m4s", map[string]interface{}{"Secs": "93784"}); err != nil {
		t.Error(err)
	}
}

func TestDurationRound(t *testing.T) {
	tpl := "{{ durationRound .Time }}"
	if err := runtv(tpl, "2h", map[string]interface{}{"Time": "2h5s"}); err != nil {
		t.Error(err)
	}
	if err := runtv(tpl, "1d", map[string]interface{}{"Time": "24h5s"}); err != nil {
		t.Error(err)
	}
	if err := runtv(tpl, "3mo", map[string]interface{}{"Time": "2400h5s"}); err != nil {
		t.Error(err)
	}
}

func TestDateModify(t *testing.T) {
	tm, err := time.Parse("02 Jan 06 15:04:05 MST", "13 Jun 19 20:39:39 GMT")
	if err != nil {
		t.Error(err)
	}
	tpl := `{{date_modify "24h" .Time}}`

	if err = runtv(tpl, "2019-06-14 20:39:39 +0000 GMT", map[string]interface{}{"Time": tm}); err != nil {
		t.Error(err)
	}
}

func TestDateMustModifyReturnsErr(t *testing.T) {
	tm, err := time.Parse("02 Jan 06 15:04:05 MST", "13 Jun 19 20:39:39 GMT")
	if err != nil {
		t.Error(err)
	}
	tpl := `{{must_date_modify "1f" .Time}}`

	if err = runtv(tpl, "2019-06-13 21:39:39 +0000 GMT", map[string]interface{}{"Time": tm}); err == nil {
		t.Error("expected err, got nil")
	}
}

func TestTimeRound(t *testing.T) {
	tm1, err := time.Parse("2006-01-02 15:04:05", "2023-05-01 10:07:00")
	require.NoError(t, err)

	tm2, err := time.Parse("2006-01-02 15:04:05", "2023-05-01 10:08:00")
	require.NoError(t, err)

	tm3, err := time.Parse("2006-01-02 15:04:05", "2023-05-01 10:35:00")
	require.NoError(t, err)

	tm4, err := time.Parse("2006-01-02 15:04:05", "2023-05-01 10:20:00")
	require.NoError(t, err)

	// Direct call tests with string duration and time.Time
	res, err := timeRound("15m", tm1)
	require.NoError(t, err)
	assert.Equal(t, "2023-05-01 10:00:00", res.Format("2006-01-02 15:04:05"))

	res, err = timeRound("15m", tm2)
	require.NoError(t, err)
	assert.Equal(t, "2023-05-01 10:15:00", res.Format("2006-01-02 15:04:05"))

	res, err = timeRound("1h", tm3)
	require.NoError(t, err)
	assert.Equal(t, "2023-05-01 11:00:00", res.Format("2006-01-02 15:04:05"))

	res, err = timeRound("1h", tm4)
	require.NoError(t, err)
	assert.Equal(t, "2023-05-01 10:00:00", res.Format("2006-01-02 15:04:05"))

	// Direct call with time.Duration
	res, err = timeRound(15*time.Minute, tm1)
	require.NoError(t, err)
	assert.Equal(t, "2023-05-01 10:00:00", res.Format("2006-01-02 15:04:05"))

	dur := 15 * time.Minute
	res, err = timeRound(&dur, tm1)
	require.NoError(t, err)
	assert.Equal(t, "2023-05-01 10:00:00", res.Format("2006-01-02 15:04:05"))

	// Direct call with pointer to time.Time
	res, err = timeRound("15m", &tm1)
	require.NoError(t, err)
	assert.Equal(t, "2023-05-01 10:00:00", res.Format("2006-01-02 15:04:05"))

	// Direct call with string date (RFC3339)
	res, err = timeRound("15m", "2023-05-01T10:07:00Z")
	require.NoError(t, err)
	assert.Equal(t, "2023-05-01 10:00:00", res.UTC().Format("2006-01-02 15:04:05"))

	// Direct call with string date (standard layout)
	res, err = timeRound("15m", "2023-05-01 10:08:00")
	require.NoError(t, err)
	assert.Equal(t, "2023-05-01 10:15:00", res.Format("2006-01-02 15:04:05"))

	// Template pipeline execution
	assert.NoError(t, runtv(`{{ .Time | timeRound "15m" | date "15:04:05" }}`, "10:00:00", map[string]interface{}{"Time": tm1}))
	assert.NoError(t, runtv(`{{ .Time | timeRound "15m" | date "15:04:05" }}`, "10:15:00", map[string]interface{}{"Time": tm2}))
	assert.NoError(t, runtv(`{{ timeRound "1h" .Time | date "15:04:05" }}`, "11:00:00", map[string]interface{}{"Time": tm3}))
	assert.NoError(t, runtv(`{{ time_round "15m" .Time | date "15:04:05" }}`, "10:00:00", map[string]interface{}{"Time": tm1}))

	// Template with string date input
	assert.NoError(t, runt(`{{ "2023-05-01T10:07:00Z" | timeRound "15m" | date "15:04:05" }}`, "10:00:00"))
	assert.NoError(t, runt(`{{ timeRound "1h" "2023-05-01T10:35:00Z" | date "15:04:05" }}`, "11:00:00"))

	// Error handling tests
	_, err = timeRound("invalid_duration", tm1)
	assert.Error(t, err)

	_, err = timeRound("15m", "invalid_date")
	assert.Error(t, err)

	_, err = timeRound((*time.Duration)(nil), tm1)
	assert.Error(t, err)

	_, err = timeRound("15m", (*time.Time)(nil))
	assert.Error(t, err)

	// Template execution error on invalid inputs
	assert.Error(t, runtv(`{{ .Time | timeRound "invalid" }}`, "", map[string]interface{}{"Time": tm1}))
	assert.Error(t, runt(`{{ "invalid_date" | timeRound "15m" }}`, ""))
}

