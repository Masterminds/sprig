package sprig

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	os.Setenv("FOO", "bar")
	tpl := `{{env "FOO"}}`
	if err := runt(tpl, "bar"); err != nil {
		t.Error(err)
	}
}

func TestExpandEnv(t *testing.T) {
	os.Setenv("FOO", "bar")
	tpl := `{{expandenv "Hello $FOO"}}`
	if err := runt(tpl, "Hello bar"); err != nil {
		t.Error(err)
	}
}

func TestBase(t *testing.T) {
	assert.NoError(t, runt(`{{ base "foo/bar" }}`, "bar"))
}

func TestDir(t *testing.T) {
	assert.NoError(t, runt(`{{ dir "foo/bar/baz" }}`, "foo/bar"))
}

func TestIsAbs(t *testing.T) {
	assert.NoError(t, runt(`{{ isAbs "/foo" }}`, "true"))
	assert.NoError(t, runt(`{{ isAbs "foo" }}`, "false"))
}

func TestClean(t *testing.T) {
	assert.NoError(t, runt(`{{ clean "/foo/../foo/../bar" }}`, "/bar"))
}

func TestExt(t *testing.T) {
	assert.NoError(t, runt(`{{ ext "/foo/bar/baz.txt" }}`, ".txt"))
}

// runt runs a template and checks that the output exactly matches the expected string.
func runt(tpl, expect string) error {
	return runtv(tpl, expect, map[string]string{})
}

// runtv takes a template, and expected return, and values for substitution.
//
// It runs the template and verifies that the output is an exact match.
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

// runRaw runs a template with the given variables and returns the result.
func runRaw(tpl string, vars interface{}) (string, error) {
	fmap := TxtFuncMap()
	t := template.Must(template.New("test").Funcs(fmap).Parse(tpl))
	var b bytes.Buffer
	err := t.Execute(&b, vars)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
