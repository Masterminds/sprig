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
	"github.com/stretchr/testify/assert"
)

// This is woefully incomplete. Please help.

func TestSubstr(t *testing.T) {
	tpl := `{{"fooo" | substr 0 3 }}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

func TestTrunc(t *testing.T) {
	tpl := `{{ "foooooo" | trunc 3 }}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

func TestQuote(t *testing.T) {
	tpl := `{{quote "a" "b" "c"}}`
	if err := runt(tpl, `"a" "b" "c"`); err != nil {
		t.Error(err)
	}
	tpl = `{{quote "\"a\"" "b" "c"}}`
	if err := runt(tpl, `"\"a\"" "b" "c"`); err != nil {
		t.Error(err)
	}
	tpl = `{{quote 1 2 3 }}`
	if err := runt(tpl, `"1" "2" "3"`); err != nil {
		t.Error(err)
	}
}
func TestSquote(t *testing.T) {
	tpl := `{{squote "a" "b" "c"}}`
	if err := runt(tpl, `'a' 'b' 'c'`); err != nil {
		t.Error(err)
	}
	tpl = `{{squote 1 2 3 }}`
	if err := runt(tpl, `'1' '2' '3'`); err != nil {
		t.Error(err)
	}
}

func TestContains(t *testing.T) {
	// Mainly, we're just verifying the paramater order swap.
	tests := []string{
		`{{if contains "cat" "fair catch"}}1{{end}}`,
		`{{if hasPrefix "cat" "catch"}}1{{end}}`,
		`{{if hasSuffix "cat" "ducat"}}1{{end}}`,
	}
	for _, tt := range tests {
		if err := runt(tt, "1"); err != nil {
			t.Error(err)
		}
	}
}

func TestTrim(t *testing.T) {
	tests := []string{
		`{{trim "   5.00   "}}`,
		`{{trimAll "$" "$5.00$"}}`,
		`{{trimPrefix "$" "$5.00"}}`,
		`{{trimSuffix "$" "5.00$"}}`,
	}
	for _, tt := range tests {
		if err := runt(tt, "5.00"); err != nil {
			t.Error(err)
		}
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

func TestToFloat64(t *testing.T) {
	target := float64(102)
	if target != toFloat64(int8(102)) {
		t.Errorf("Expected 102")
	}
	if target != toFloat64(int(102)) {
		t.Errorf("Expected 102")
	}
	if target != toFloat64(int32(102)) {
		t.Errorf("Expected 102")
	}
	if target != toFloat64(int16(102)) {
		t.Errorf("Expected 102")
	}
	if target != toFloat64(int64(102)) {
		t.Errorf("Expected 102")
	}
	if target != toFloat64("102") {
		t.Errorf("Expected 102")
	}
	if 0 != toFloat64("frankie") {
		t.Errorf("Expected 0")
	}
	if target != toFloat64(uint16(102)) {
		t.Errorf("Expected 102")
	}
	if target != toFloat64(uint64(102)) {
		t.Errorf("Expected 102")
	}
	if 102.1234 != toFloat64(float64(102.1234)) {
		t.Errorf("Expected 102.1234")
	}
	if 1 != toFloat64(true) {
		t.Errorf("Expected 102")
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

func TestToInt(t *testing.T) {
	target := int(102)
	if target != toInt(int8(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt(int(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt(int32(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt(int16(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt(int64(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt("102") {
		t.Errorf("Expected 102")
	}
	if 0 != toInt("frankie") {
		t.Errorf("Expected 0")
	}
	if target != toInt(uint16(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt(uint64(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt(float64(102.1234)) {
		t.Errorf("Expected 102")
	}
	if 1 != toInt(true) {
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

	dict := map[string]interface{}{"top": map[string]interface{}{}}
	tpl = `{{if empty .top.NoSuchThing}}1{{else}}0{{end}}`
	if err := runtv(tpl, "1", dict); err != nil {
		t.Error(err)
	}
	tpl = `{{if empty .bottom.NoSuchThing}}1{{else}}0{{end}}`
	if err := runtv(tpl, "1", dict); err != nil {
		t.Error(err)
	}
}

func TestCoalesce(t *testing.T) {
	tests := map[string]string{
		`{{ coalesce 1 }}`:                            "1",
		`{{ coalesce "" 0 nil 2 }}`:                   "2",
		`{{ $two := 2 }}{{ coalesce "" 0 nil $two }}`: "2",
		`{{ $two := 2 }}{{ coalesce "" $two 0 0 0 }}`: "2",
		`{{ $two := 2 }}{{ coalesce "" $two 3 4 5 }}`: "2",
		`{{ coalesce }}`:                              "<no value>",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}

	dict := map[string]interface{}{"top": map[string]interface{}{}}
	tpl := `{{ coalesce .top.NoSuchThing .bottom .bottom.dollar "airplane"}}`
	if err := runtv(tpl, "airplane", dict); err != nil {
		t.Error(err)
	}
}

func TestCompact(t *testing.T) {
	tests := map[string]string{
		`{{ list 1 0 "" "hello" | compact }}`: `[1 hello]`,
		`{{ list "" "" | compact }}`:          `[]`,
		`{{ list | compact }}`:                `[]`,
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestSplit(t *testing.T) {
	tpl := `{{$v := "foo$bar$baz" | split "$"}}{{$v._0}}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

func TestToString(t *testing.T) {
	tpl := `{{ toString 1 | kindOf }}`
	assert.NoError(t, runt(tpl, "string"))
}

func TestToStrings(t *testing.T) {
	tpl := `{{ $s := list 1 2 3 | toStrings }}{{ index $s 1 | kindOf }}`
	assert.NoError(t, runt(tpl, "string"))
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

func TestCat(t *testing.T) {
	tpl := `{{$b := "b"}}{{"c" | cat "a" $b}}`
	if err := runt(tpl, "a b c"); err != nil {
		t.Error(err)
	}
}

func TestIndent(t *testing.T) {
	tpl := `{{indent 4 "a\nb\nc"}}`
	if err := runt(tpl, "    a\n    b\n    c"); err != nil {
		t.Error(err)
	}
}

func TestReplace(t *testing.T) {
	tpl := `{{"I Am Henry VIII" | replace " " "-"}}`
	if err := runt(tpl, "I-Am-Henry-VIII"); err != nil {
		t.Error(err)
	}
}

func TestPlural(t *testing.T) {
	tpl := `{{$num := len "two"}}{{$num}} {{$num | plural "1 char" "chars"}}`
	if err := runt(tpl, "3 chars"); err != nil {
		t.Error(err)
	}
	tpl = `{{len "t" | plural "cheese" "%d chars"}}`
	if err := runt(tpl, "cheese"); err != nil {
		t.Error(err)
	}
}

func TestSha256Sum(t *testing.T) {
	tpl := `{{"abc" | sha256sum}}`
	if err := runt(tpl, "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"); err != nil {
		t.Error(err)
	}
}

func TestTuple(t *testing.T) {
	tpl := `{{$t := tuple 1 "a" "foo"}}{{index $t 2}}{{index $t 0 }}{{index $t 1}}`
	if err := runt(tpl, "foo1a"); err != nil {
		t.Error(err)
	}
}

func TestList(t *testing.T) {
	tpl := `{{$t := list 1 "a" "foo"}}{{index $t 2}}{{index $t 0 }}{{index $t 1}}`
	if err := runt(tpl, "foo1a"); err != nil {
		t.Error(err)
	}
}

func TestPush(t *testing.T) {
	// Named `append` in the function map
	tests := map[string]string{
		`{{ $t := tuple 1 2 3  }}{{ append $t 4 | len }}`:        "4",
		`{{ $t := tuple 1 2 3 4  }}{{ append $t 5 | join "-" }}`: "1-2-3-4-5",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}
func TestPrepend(t *testing.T) {
	tests := map[string]string{
		`{{ $t := tuple 1 2 3  }}{{ prepend $t 0 | len }}`:        "4",
		`{{ $t := tuple 1 2 3 4  }}{{ prepend $t 0 | join "-" }}`: "0-1-2-3-4",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestFirst(t *testing.T) {
	tests := map[string]string{
		`{{ list 1 2 3 | first }}`: "1",
		`{{ list | first }}`:       "<no value>",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}
func TestLast(t *testing.T) {
	tests := map[string]string{
		`{{ list 1 2 3 | last }}`: "3",
		`{{ list | last }}`:       "<no value>",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestInitial(t *testing.T) {
	tests := map[string]string{
		`{{ list 1 2 3 | initial | len }}`:   "2",
		`{{ list 1 2 3 | initial | last }}`:  "2",
		`{{ list 1 2 3 | initial | first }}`: "1",
		`{{ list | initial }}`:               "[]",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestRest(t *testing.T) {
	tests := map[string]string{
		`{{ list 1 2 3 | rest | len }}`:   "2",
		`{{ list 1 2 3 | rest | last }}`:  "3",
		`{{ list 1 2 3 | rest | first }}`: "2",
		`{{ list | rest }}`:               "[]",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestReverse(t *testing.T) {
	tests := map[string]string{
		`{{ list 1 2 3 | reverse | first }}`:        "3",
		`{{ list 1 2 3 | reverse | rest | first }}`: "2",
		`{{ list 1 2 3 | reverse | last }}`:         "1",
		`{{ list 1 2 3 4 | reverse }}`:              "[4 3 2 1]",
		`{{ list 1 | reverse }}`:                    "[1]",
		`{{ list | reverse }}`:                      "[]",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
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

func TestUnset(t *testing.T) {
	tpl := `{{- $d := dict "one" 1 "two" 222222 -}}
	{{- $_ := unset $d "two" -}}
	{{- range $k, $v := $d}}{{$k}}{{$v}}{{- end -}}
	`

	expect := "one1"
	if err := runt(tpl, expect); err != nil {
		t.Error(err)
	}
}
func TestHasKey(t *testing.T) {
	tpl := `{{- $d := dict "one" 1 "two" 222222 -}}
	{{- if hasKey $d "one" -}}1{{- end -}}
	`

	expect := "1"
	if err := runt(tpl, expect); err != nil {
		t.Error(err)
	}
}

func TestPluck(t *testing.T) {
	tpl := `
	{{- $d := dict "one" 1 "two" 222222 -}}
	{{- $d2 := dict "one" 1 "two" 33333 -}}
	{{- $d3 := dict "one" 1 -}}
	{{- $d4 := dict "one" 1 "two" 4444 -}}
	{{- pluck "two" $d $d2 $d3 $d4 -}}
	`

	expect := "[222222 33333 4444]"
	if err := runt(tpl, expect); err != nil {
		t.Error(err)
	}
}

func TestKeys(t *testing.T) {
	tests := map[string]string{
		`{{ dict "foo" 1 "bar" 2 | keys | sortAlpha }}`: "[bar foo]",
		`{{ dict | keys }}`:                             "[]",
	}
	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
	}
}

func TestSet(t *testing.T) {
	tpl := `{{- $d := dict "one" 1 "two" 222222 -}}
	{{- $_ := set $d "two" 2 -}}
	{{- $_ := set $d "three" 3 -}}
	{{- if hasKey $d "one" -}}{{$d.one}}{{- end -}}
	{{- if hasKey $d "two" -}}{{$d.two}}{{- end -}}
	{{- if hasKey $d "three" -}}{{$d.three}}{{- end -}}
	`

	expect := "123"
	if err := runt(tpl, expect); err != nil {
		t.Error(err)
	}
}

func TestUntil(t *testing.T) {
	tests := map[string]string{
		`{{range $i, $e := until 5}}{{$i}}{{$e}}{{end}}`:   "0011223344",
		`{{range $i, $e := until -5}}{{$i}}{{$e}} {{end}}`: "00 1-1 2-2 3-3 4-4 ",
	}
	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
	}
}
func TestUntilStep(t *testing.T) {
	tests := map[string]string{
		`{{range $i, $e := untilStep 0 5 1}}{{$i}}{{$e}}{{end}}`:     "0011223344",
		`{{range $i, $e := untilStep 3 6 1}}{{$i}}{{$e}}{{end}}`:     "031425",
		`{{range $i, $e := untilStep 0 -10 -2}}{{$i}}{{$e}} {{end}}`: "00 1-2 2-4 3-6 4-8 ",
		`{{range $i, $e := untilStep 3 0 1}}{{$i}}{{$e}}{{end}}`:     "",
		`{{range $i, $e := untilStep 3 99 0}}{{$i}}{{$e}}{{end}}`:    "",
		`{{range $i, $e := untilStep 3 99 -1}}{{$i}}{{$e}}{{end}}`:   "",
		`{{range $i, $e := untilStep 3 0 0}}{{$i}}{{$e}}{{end}}`:     "",
	}
	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
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

func TestJoin(t *testing.T) {
	assert.NoError(t, runt(`{{ tuple "a" "b" "c" | join "-" }}`, "a-b-c"))
	assert.NoError(t, runt(`{{ tuple 1 2 3 | join "-" }}`, "1-2-3"))
	assert.NoError(t, runtv(`{{ join "-" .V }}`, "a-b-c", map[string]interface{}{"V": []string{"a", "b", "c"}}))
	assert.NoError(t, runtv(`{{ join "-" .V }}`, "abc", map[string]interface{}{"V": "abc"}))
	assert.NoError(t, runtv(`{{ join "-" .V }}`, "1-2-3", map[string]interface{}{"V": []int{1, 2, 3}}))
}

func TestSortAlpha(t *testing.T) {
	// Named `append` in the function map
	tests := map[string]string{
		`{{ list "c" "a" "b" | sortAlpha | join "" }}`: "abc",
		`{{ list 2 1 4 3 | sortAlpha | join "" }}`:     "1234",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestDelete(t *testing.T) {
	fmap := TxtFuncMap()
	delete(fmap, "split")
	if _, ok := fmap["split"]; ok {
		t.Error("Failed to delete split from map")
	}
}

func TestDerivePassword(t *testing.T) {
	expectations := map[string]string{
		`{{derivePassword 1 "long" "password" "user" "example.com"}}`:    "ZedaFaxcZaso9*",
		`{{derivePassword 2 "long" "password" "user" "example.com"}}`:    "Fovi2@JifpTupx",
		`{{derivePassword 1 "maximum" "password" "user" "example.com"}}`: "pf4zS1LjCg&LjhsZ7T2~",
		`{{derivePassword 1 "medium" "password" "user" "example.com"}}`:  "ZedJuz8$",
		`{{derivePassword 1 "basic" "password" "user" "example.com"}}`:   "pIS54PLs",
		`{{derivePassword 1 "short" "password" "user" "example.com"}}`:   "Zed5",
		`{{derivePassword 1 "pin" "password" "user" "example.com"}}`:     "6685",
	}

	for tpl, result := range expectations {
		out, err := runRaw(tpl, nil)
		if err != nil {
			t.Error(err)
		}
		if 0 != strings.Compare(out, result) {
			t.Error("Generated password does not match for", tpl)
		}
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

func TestUUIDGeneration(t *testing.T) {
	tpl := `{{uuidv4}}`
	out, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}

	if len(out) != 36 {
		t.Error("Expected UUID of length 36")
	}

	out2, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}

	if out == out2 {
		t.Error("Expected subsequent UUID generations to be different")
	}
}

func TestXMLEncode(t *testing.T) {
	tpl := `{{"test&" | xml }}`
	if err := runt(tpl, "test&amp;"); err != nil {
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
