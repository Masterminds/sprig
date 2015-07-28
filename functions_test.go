package sprig

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"
)

// This is woefully incomplete. Please help.

func TestSubstr(t *testing.T) {
	tpl := `{{"fooo" | substr 0 3 }}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

func TestTrimall(t *testing.T) {
	tpl := `{{"$foo$" | trimall "$"}}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

func TestSplit(t *testing.T) {
	tpl := `{{$v := "foo$bar$baz" | split "$"}}{{$v._0}}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

func runt(tpl, expect string) error {
	fmap := TxtFuncMap()
	t := template.Must(template.New("test").Funcs(fmap).Parse(tpl))
	var b bytes.Buffer
	err := t.Execute(&b, "")
	if err != nil {
		return err
	}
	if expect != b.String() {
		return fmt.Errorf("Expected '%s', got '%s'", expect, b.String())
	}
	return nil
}
