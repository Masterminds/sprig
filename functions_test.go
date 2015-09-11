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

func TestDefault(t *testing.T) {
	tpl := `{{"" | default "foo"}}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
	tpl = `{{default "foo" 234}}`
	if err := runt(tpl, "234"); err != nil {
		t.Error(err)
	}
	tpl = `{{default "foo" 2.34}}`
	if err := runt(tpl, "2.34"); err != nil {
		t.Error(err)
	}
}

func TestSplit(t *testing.T) {
	tpl := `{{$v := "foo$bar$baz" | split "$"}}{{$v._0}}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

type fixtureTO struct {
	Name, Value string
}

func TestTypeOf(t *testing.T) {
	f := &fixtureTO{"hello", "world"}
	tpl := `{{typeOf .}}`
	if err := runtv(tpl, "*sprig.fixtureTO", f); err != nil {
		t.Error(err)
	}
}

func runt(tpl, expect string) error {
	return runtv(tpl, expect, "")
}
func runtv(tpl, expect string, vars interface{}) error {
	fmap := TxtFuncMap()
	t := template.Must(template.New("test").Funcs(fmap).Parse(tpl))
	var b bytes.Buffer
	err := t.Execute(&b, vars)
	if err != nil {
		return err
	}
	if expect != b.String() {
		return fmt.Errorf("Expected '%s', got '%s'", expect, b.String())
	}
	return nil
}
