package sprig

import (
	"testing"
	"time"
)

func TestHtmlDate(t *testing.T) {
	t.Skip()
	tpl := `{{ htmlDate 0}}`
	if err := runt(tpl, "1970-01-01"); err != nil {
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
