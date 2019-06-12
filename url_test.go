package sprig

import (
	"testing"
)

func TestPathEscape(t *testing.T) {
	tpl := `{{ "Hello, World + Sprig!" | pathEscape }}`
	if err := runt(tpl, "Hello%2C%20World%20+%20Sprig%21"); err != nil {
		t.Error(err)
	}
}

func TestPathUnescape(t *testing.T) {
	tpl := `{{ "Hello%2C%20World%20+%20Sprig%21" | pathUnescape }}`
	if err := runt(tpl, "Hello, World + Sprig!"); err != nil {
		t.Error(err)
	}
}

func TestQueryEscape(t *testing.T) {
	tpl := `{{ "Hello, World + Sprig!" | queryEscape }}`
	if err := runt(tpl, "Hello%2C+World+%2B+Sprig%21"); err != nil {
		t.Error(err)
	}
}

func TestQueryUnescape(t *testing.T) {
	tpl := `{{ "Hello%2C+World+%2B+Sprig%21" | queryUnescape }}`
	if err := runt(tpl, "Hello, World + Sprig!"); err != nil {
		t.Error(err)
	}
}
