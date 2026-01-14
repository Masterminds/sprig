package sprig

import (
	"math"
	"strings"
	"testing"
	"text/template"
	"time"

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

func TestDurationHelpers(t *testing.T) {
	tests := []struct {
		name   string
		tpl    string
		vars   interface{}
		expect string
	}{{
		name:   "durationSeconds parses duration string",
		tpl:    `{{ durationSeconds "1m30s" }}`,
		expect: `90`,
	}, {
		tpl:    `{{ mustToDuration 30 }}`,
		expect: `30s`,
	}, {
		tpl:    `{{ mustToDuration "1m30s" }}`,
		expect: `1m30s`,
	}, {
		name:   "durationSeconds parses numeric string as seconds",
		tpl:    `{{ durationSeconds "2.5" }}`,
		expect: `2.5`,
	}, {
		name:   "durationSeconds trims whitespace around numeric string",
		tpl:    `{{ durationSeconds "  2.5  " }}`,
		expect: `2.5`,
	}, {
		name:   "durationSeconds int treated as seconds",
		tpl:    `{{ durationSeconds 2 }}`,
		expect: `2`,
	}, {
		name:   "durationSeconds float treated as seconds",
		tpl:    `{{ durationSeconds 2.5 }}`,
		expect: `2.5`,
	}, {
		name:   "durationSeconds uint treated as seconds",
		tpl:    `{{ durationSeconds . }}`,
		vars:   uint(2),
		expect: `2`,
	}, {
		name:   "durationSeconds time.Duration passthrough",
		tpl:    `{{ durationSeconds . }}`,
		vars:   1500 * time.Millisecond,
		expect: `1.5`,
	}, {
		name:   "invalid duration string returns 0",
		tpl:    `{{ durationSeconds "nope" }}`,
		expect: `0`,
	}, {
		name:   "empty duration string returns 0",
		tpl:    `{{ durationSeconds "" }}`,
		expect: `0`,
	}, {
		name:   "whitespace-only duration string returns 0",
		tpl:    `{{ durationSeconds "   " }}`,
		expect: `0`,
	}, {
		name:   "nil returns 0",
		tpl:    `{{ durationSeconds . }}`,
		vars:   nil,
		expect: `0`,
	}, {
		name:   "durationSeconds uint overflow returns 0",
		tpl:    `{{ durationSeconds . }}`,
		vars:   uint64(math.MaxInt64) + 1,
		expect: `0`,
	}, {
		name:   "durationMilliseconds int seconds",
		tpl:    `{{ durationMilliseconds 2 }}`,
		expect: `2000`,
	}, {
		name:   "durationMilliseconds float seconds",
		tpl:    `{{ durationMilliseconds 1.5 }}`,
		expect: `1500`,
	}, {
		name:   "durationMicroseconds int seconds",
		tpl:    `{{ durationMicroseconds 2 }}`,
		expect: `2000000`,
	}, {
		name:   "durationNanoseconds int seconds",
		tpl:    `{{ durationNanoseconds 2 }}`,
		expect: `2000000000`,
	}, {
		name:   "durationMinutes parses duration string",
		tpl:    `{{ durationMinutes "90s" }}`,
		expect: `1.5`,
	}, {
		name:   "durationHours parses duration string",
		tpl:    `{{ durationHours "90m" }}`,
		expect: `1.5`,
	}, {
		name:   "durationDays parses duration string",
		tpl:    `{{ durationDays "36h" }}`,
		expect: `1.5`,
	}, {
		name:   "durationDays numeric seconds",
		tpl:    `{{ durationDays 86400 }}`,
		expect: `1`,
	}, {
		name:   "durationRoundTo numeric seconds",
		tpl:    `{{ durationRoundTo 93 60 }}`, // 93s rounded to 60s = 120s
		expect: `2m0s`,
	}, {
		name:   "durationTruncateTo numeric seconds",
		tpl:    `{{ durationTruncateTo 93 60 }}`, // 93s truncated to 60s = 60s
		expect: `1m0s`,
	}, {
		name:   "durationRoundTo accepts duration-string multiplier",
		tpl:    `{{ durationRoundTo "93s" "1m" }}`,
		expect: `2m0s`,
	}, {
		name:   "durationTruncateTo accepts duration-string multiplier",
		tpl:    `{{ durationTruncateTo "93s" "1m" }}`,
		expect: `1m0s`,
	}, {
		name:   "durationRoundTo invalid m returns v unchanged",
		tpl:    `{{ durationRoundTo "93s" "nope" }}`,
		expect: `1m33s`,
	}, {
		name:   "durationTruncateTo invalid m returns v unchanged",
		tpl:    `{{ durationTruncateTo "93s" "nope" }}`,
		expect: `1m33s`,
	}, {
		name:   "durationRoundTo zero m returns v unchanged",
		tpl:    `{{ durationRoundTo "93s" 0 }}`,
		expect: `1m33s`,
	}, {
		name:   "durationTruncateTo negative m returns v unchanged",
		tpl:    `{{ durationTruncateTo "93s" -1 }}`,
		expect: `1m33s`,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b strings.Builder
			err := template.Must(template.New("test").Funcs(TxtFuncMap()).Parse(tt.tpl)).Execute(&b, tt.vars)
			require.NoError(t, err, tt.tpl)
			require.Equal(t, tt.expect, b.String(), tt.tpl)
		})
	}

	mustErrTests := []struct {
		name string
		tpl  string
		vars interface{}
	}{{
		name: "mustToDuration invalid string",
		tpl:  `{{ mustToDuration "nope" }}`,
	}, {
		name: "mustToDuration empty string",
		tpl:  `{{ mustToDuration "" }}`,
	}, {
		name: "mustToDuration whitespace string",
		tpl:  `{{ mustToDuration "   " }}`,
	}, {
		name: "mustToDuration unsupported type",
		tpl:  `{{ mustToDuration . }}`,
		vars: []int{1, 2, 3},
	}, {
		name: "mustToDuration uint overflow",
		tpl:  `{{ mustToDuration . }}`,
		vars: uint64(math.MaxInt64) + 1,
	},
	}

	for _, tt := range mustErrTests {
		t.Run(tt.name, func(t *testing.T) {
			var b strings.Builder
			err := template.Must(template.New("test").Funcs(TxtFuncMap()).Parse(tt.tpl)).Execute(&b, tt.vars)
			require.Error(t, err, tt.tpl)
		})
	}
}
