package sprig

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"testing"
	"text/template"
	"time"

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

func TestSnakeCase(t *testing.T) {
	assert.NoError(t, runt(`{{ snakecase "FirstName" }}`, "first_name"))
	assert.NoError(t, runt(`{{ snakecase "HTTPServer" }}`, "http_server"))
	assert.NoError(t, runt(`{{ snakecase "NoHTTPS" }}`, "no_https"))
	assert.NoError(t, runt(`{{ snakecase "GO_PATH" }}`, "go_path"))
	assert.NoError(t, runt(`{{ snakecase "GO PATH" }}`, "go_path"))
	assert.NoError(t, runt(`{{ snakecase "GO-PATH" }}`, "go_path"))
}

func TestCamelCase(t *testing.T) {
	assert.NoError(t, runt(`{{ camelcase "http_server" }}`, "HttpServer"))
	assert.NoError(t, runt(`{{ camelcase "_camel_case" }}`, "_CamelCase"))
	assert.NoError(t, runt(`{{ camelcase "no_https" }}`, "NoHttps"))
	assert.NoError(t, runt(`{{ camelcase "_complex__case_" }}`, "_Complex_Case_"))
	assert.NoError(t, runt(`{{ camelcase "all" }}`, "All"))
}

func TestKebabCase(t *testing.T) {
	assert.NoError(t, runt(`{{ kebabcase "FirstName" }}`, "first-name"))
	assert.NoError(t, runt(`{{ kebabcase "HTTPServer" }}`, "http-server"))
	assert.NoError(t, runt(`{{ kebabcase "NoHTTPS" }}`, "no-https"))
	assert.NoError(t, runt(`{{ kebabcase "GO_PATH" }}`, "go-path"))
	assert.NoError(t, runt(`{{ kebabcase "GO PATH" }}`, "go-path"))
	assert.NoError(t, runt(`{{ kebabcase "GO-PATH" }}`, "go-path"))
}

func TestShuffle(t *testing.T) {
	defer rand.Seed(time.Now().UnixNano())
	rand.Seed(1)
	// Because we're using a random number generator, we need these to go in
	// a predictable sequence:
	assert.NoError(t, runt(`{{ shuffle "Hello World" }}`, "rldo HWlloe"))
}

func TestRegex(t *testing.T) {
	assert.NoError(t, runt(`{{ regexQuoteMeta "1.2.3" }}`, "1\\.2\\.3"))
	assert.NoError(t, runt(`{{ regexQuoteMeta "pretzel" }}`, "pretzel"))
}

func TestRandStringFromRegex(t *testing.T) {
	tmplStr := "{{randFromRegex \"https://(www[.])example[.]com/[a-zA-Z0-9]{9,15}\"}}"
	tmpl, err := template.New("randFromRegex").Funcs(FuncMap()).Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		panic(err)
	}

	// Print the generated random string
	fmt.Println(buf.String())

	tmplStr = "{{randFromRegex \"/v4/providers/[^/]+/roles/[^/]+/groups$\"}}"
	tmpl, err = template.New("randFromRegex").Funcs(FuncMap()).Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	var buf2 bytes.Buffer
	err = tmpl.Execute(&buf2, nil)
	if err != nil {
		panic(err)
	}

	// Print the generated random string
	fmt.Println(buf2.String())

	/*
		further make sure it accept perl chars
	*/
	tmplStr = `{{randFromRegex "/v\\d+/providers/[^/]+/roles\\b"}}`
	tmpl, err = template.New("randFromRegex").Funcs(FuncMap()).Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	var buf3 bytes.Buffer
	err = tmpl.Execute(&buf3, nil)
	if err != nil {
		panic(err)
	}

	// Print the generated random string
	fmt.Println(buf3.String())
}

func TestRandStringFromUrlRegex(t *testing.T) {
	regexStr := `https://(www[.])example[.]com/[a-zA-Z0-9]{9,15}/test@url/[a-e]+/[^/]+$`
	tmplStr := fmt.Sprintf("{{randFromUrlRegex \"%s\"}}", regexStr)
	tmpl, err := template.New("randFromUrlRegex").Funcs(FuncMap()).Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err1 := tmpl.Execute(&buf, nil)
	if err1 != nil {
		panic(err1)
	}

	randomUrl := buf.String()
	regexMatcher, err2 := regexp.Compile(regexStr)
	if err2 != nil {
		panic(err2)
	}

	matched := regexMatcher.MatchString(randomUrl)
	if !matched {
		panic(errors.New("the generated random url does not match the original regular expression"))
	}

	reservedChars := getIllegalUrlCharMap()
	runeSlice := []rune(randomUrl)
	for _, curChar := range runeSlice {
		if curChar == ':' || curChar == '/' || curChar == '@' {
			continue
		}

		if _, exists := reservedChars[curChar]; exists {
			panic(fmt.Sprintf("char [%c] exists in generated random url string", curChar))
		}
	}
	// Print the generated random string
	fmt.Println(randomUrl)
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
