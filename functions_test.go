package sprig

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/aokoli/goutils"
)

// This is woefully incomplete. Please help.

func TestSubstr(t *testing.T) {
	tpl := `{{"fooo" | substr 0 3 }}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

func TestQuote(t *testing.T) {
	tpl := `{{quote "a" "b" "c"}}`
	if err := runt(tpl, `"a" "b" "c"`); err != nil {
		t.Error(err)
	}
}
func TestSquote(t *testing.T) {
	tpl := `{{squote "a" "b" "c"}}`
	if err := runt(tpl, `'a' 'b' 'c'`); err != nil {
		t.Error(err)
	}
}

func TestAdd(t *testing.T) {
	tpl := `{{ 3 | add 1 2}}`
	if err := runt(tpl, `6`); err != nil {
		t.Error(err)
	}
}

func TestMul(t *testing.T) {
	tpl := `{{ 1 | mul "2" 3 "4"}}`
	if err := runt(tpl, `24`); err != nil {
		t.Error(err)
	}
}

func TestHtmlDate(t *testing.T) {
	t.Skip()
	tpl := `{{ htmlDate 0}}`
	if err := runt(tpl, "1970-01-01"); err != nil {
		t.Error(err)
	}
}

func TestBiggest(t *testing.T) {
	tpl := `{{ biggest 1 2 3 345 5 6 7}}`
	if err := runt(tpl, `345`); err != nil {
		t.Error(err)
	}

	tpl = `{{ max 345}}`
	if err := runt(tpl, `345`); err != nil {
		t.Error(err)
	}
}
func TestMin(t *testing.T) {
	tpl := `{{ min 1 2 3 345 5 6 7}}`
	if err := runt(tpl, `1`); err != nil {
		t.Error(err)
	}

	tpl = `{{ min 345}}`
	if err := runt(tpl, `345`); err != nil {
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

	tpl = `{{ .Nothing | default "123" }}`
	if err := runt(tpl, "123"); err != nil {
		t.Error(err)
	}
	tpl = `{{ default "123" }}`
	if err := runt(tpl, "123"); err != nil {
		t.Error(err)
	}
}

func TestToInt64(t *testing.T) {
	target := int64(102)
	if target != toInt64(int8(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt64(int(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt64(int32(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt64(int16(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt64(int64(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt64("102") {
		t.Errorf("Expected 102")
	}
	if 0 != toInt64("frankie") {
		t.Errorf("Expected 0")
	}
	if target != toInt64(uint16(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt64(uint64(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt64(float64(102.1234)) {
		t.Errorf("Expected 102")
	}
	if 1 != toInt64(true) {
		t.Errorf("Expected 102")
	}
}

func TestEmpty(t *testing.T) {
	tpl := `{{if empty 1}}1{{else}}0{{end}}`
	if err := runt(tpl, "0"); err != nil {
		t.Error(err)
	}

	tpl = `{{if empty 0}}1{{else}}0{{end}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}
	tpl = `{{if empty ""}}1{{else}}0{{end}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}
	tpl = `{{if empty 0.0}}1{{else}}0{{end}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}
	tpl = `{{if empty false}}1{{else}}0{{end}}`
	if err := runt(tpl, "1"); err != nil {
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

func TestKindOf(t *testing.T) {
	tpl := `{{kindOf .}}`

	f := fixtureTO{"hello", "world"}
	if err := runtv(tpl, "struct", f); err != nil {
		t.Error(err)
	}

	f2 := []string{"hello"}
	if err := runtv(tpl, "slice", f2); err != nil {
		t.Error(err)
	}

	var f3 *fixtureTO = nil
	if err := runtv(tpl, "ptr", f3); err != nil {
		t.Error(err)
	}
}

func TestTypeIs(t *testing.T) {
	f := &fixtureTO{"hello", "world"}
	tpl := `{{if typeIs "*sprig.fixtureTO" .}}t{{else}}f{{end}}`
	if err := runtv(tpl, "t", f); err != nil {
		t.Error(err)
	}

	f2 := "hello"
	if err := runtv(tpl, "f", f2); err != nil {
		t.Error(err)
	}
}
func TestTypeIsLike(t *testing.T) {
	f := "foo"
	tpl := `{{if typeIsLike "string" .}}t{{else}}f{{end}}`
	if err := runtv(tpl, "t", f); err != nil {
		t.Error(err)
	}

	// Now make a pointer. Should still match.
	f2 := &f
	if err := runtv(tpl, "t", f2); err != nil {
		t.Error(err)
	}
}
func TestKindIs(t *testing.T) {
	f := &fixtureTO{"hello", "world"}
	tpl := `{{if kindIs "ptr" .}}t{{else}}f{{end}}`
	if err := runtv(tpl, "t", f); err != nil {
		t.Error(err)
	}
	f2 := "hello"
	if err := runtv(tpl, "f", f2); err != nil {
		t.Error(err)
	}
}

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

func TestBase64EncodeDecode(t *testing.T) {
	magicWord := "coffee"
	expect := base64.StdEncoding.EncodeToString([]byte(magicWord))

	if expect == magicWord {
		t.Fatal("Encoder doesn't work.")
	}

	tpl := `{{b64enc "coffee"}}`
	if err := runt(tpl, expect); err != nil {
		t.Error(err)
	}
	tpl = fmt.Sprintf("{{b64dec %q}}", expect)
	if err := runt(tpl, magicWord); err != nil {
		t.Error(err)
	}
}
func TestBase32EncodeDecode(t *testing.T) {
	magicWord := "coffee"
	expect := base32.StdEncoding.EncodeToString([]byte(magicWord))

	if expect == magicWord {
		t.Fatal("Encoder doesn't work.")
	}

	tpl := `{{b32enc "coffee"}}`
	if err := runt(tpl, expect); err != nil {
		t.Error(err)
	}
	tpl = fmt.Sprintf("{{b32dec %q}}", expect)
	if err := runt(tpl, magicWord); err != nil {
		t.Error(err)
	}
}

func TestGoutils(t *testing.T) {
	tests := map[string]string{
		`{{abbrev 5 "hello world"}}`:           "he...",
		`{{abbrevboth 5 10 "1234 5678 9123"}}`: "...5678...",
		`{{nospace "h e l l o "}}`:             "hello",
		`{{untitle "First Try"}}`:              "first try", //https://youtu.be/44-RsrF_V_w
		`{{initials "First Try"}}`:             "FT",
		`{{wrap 5 "Hello World"}}`:             "Hello\nWorld",
		`{{wrapWith 5 "\t" "Hello World"}}`:    "Hello\tWorld",
	}
	for k, v := range tests {
		t.Log(k)
		if err := runt(k, v); err != nil {
			t.Errorf("Error on tpl %s: %s", err)
		}
	}
}

func TestRandom(t *testing.T) {
	// One of the things I love about Go:
	goutils.RANDOM = rand.New(rand.NewSource(1))

	// Because we're using a random number generator, we need these to go in
	// a predictable sequence:
	if err := runt(`{{randAlphaNum 5}}`, "9bzRv"); err != nil {
		t.Errorf("Error on tpl %s: %s", err)
	}
	if err := runt(`{{randAlpha 5}}`, "VjwGe"); err != nil {
		t.Errorf("Error on tpl %s: %s", err)
	}
	if err := runt(`{{randAscii 5}}`, "1KA5p"); err != nil {
		t.Errorf("Error on tpl %s: %s", err)
	}
	if err := runt(`{{randNumeric 5}}`, "26018"); err != nil {
		t.Errorf("Error on tpl %s: %s", err)
	}

}

func TestContains(t *testing.T) {
	tpl := `{{"foobar" | contains "foo"}}`
	if err := runt(tpl, "true"); err != nil {
		t.Error(err)
	}
}

func TestTuple(t *testing.T) {
	tpl := `{{$t := tuple 1 "a" "foo"}}{{index $t 2}}{{index $t 0 }}{{index $t 1}}`
	if err := runt(tpl, "foo1a"); err != nil {
		t.Error(err)
	}
}

func TestPluralize(t *testing.T) {
	tpl := `{{ pluralize 1 "abc" "def"}}`
	if err := runt(tpl, "abc"); err != nil {
		t.Error(err)
	}
	tpl = `{{ pluralize 2 "abc" "def" }}`
	if err := runt(tpl, "def"); err != nil {
		t.Error(err)
	}
}

func TestDict(t *testing.T) {
	tpl := `{{$d := dict 1 2 "three" "four" 5}}{{range $k, $v := $d}}{{$k}}{{$v}}{{end}}`
	out, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if len(out) != 12 {
		t.Errorf("Expected length 12, got %d", len(out))
	}
	// dict does not guarantee ordering because it is backed by a map.
	if !strings.Contains(out, "12") {
		t.Error("Expected grouping 12")
	}
	if !strings.Contains(out, "threefour") {
		t.Error("Expected grouping threefour")
	}
	if !strings.Contains(out, "5") {
		t.Error("Expected 5")
	}
	tpl = `{{$t := dict "I" "shot" "the" "albatross"}}{{$t.the}} {{$t.I}}`
	if err := runt(tpl, "albatross shot"); err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	fmap := TxtFuncMap()
	delete(fmap, "split")
	if _, ok := fmap["split"]; ok {
		t.Error("Failed to delete split from map")
	}
}

// NOTE(bacongobbler): this test is really _slow_ because of how long it takes to compute
// and generate a new crypto key.
func TestGenPrivateKey(t *testing.T) {
	// test that calling by default generates an RSA private key
	tpl := `{{genPrivateKey ""}}`
	out, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(out, "RSA PRIVATE KEY") {
		t.Error("Expected RSA PRIVATE KEY")
	}
	// test all acceptable arguments
	tpl = `{{genPrivateKey "rsa"}}`
	out, err = runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(out, "RSA PRIVATE KEY") {
		t.Error("Expected RSA PRIVATE KEY")
	}
	tpl = `{{genPrivateKey "dsa"}}`
	out, err = runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(out, "DSA PRIVATE KEY") {
		t.Error("Expected DSA PRIVATE KEY")
	}
	tpl = `{{genPrivateKey "ecdsa"}}`
	out, err = runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(out, "EC PRIVATE KEY") {
		t.Error("Expected EC PRIVATE KEY")
	}
	// test bad
	tpl = `{{genPrivateKey "bad"}}`
	out, err = runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if out != "Unknown type bad" {
		t.Error("Expected type 'bad' to be an unknown crypto algorithm")
	}
	// ensure that we can base64 encode the string
	tpl = `{{genPrivateKey "rsa" | b64enc}}`
	out, err = runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
}

func runt(tpl, expect string) error {
	return runtv(tpl, expect, map[string]string{})
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
