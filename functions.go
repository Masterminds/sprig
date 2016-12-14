/*
Sprig: Template functions for Go.

This package contains a number of utility functions for working with data
inside of Go `html/template` and `text/template` files.

To add these functions, use the `template.Funcs()` method:

	t := templates.New("foo").Funcs(sprig.FuncMap())

Note that you should add the function map before you parse any template files.

	In several cases, Sprig reverses the order of arguments from the way they
	appear in the standard library. This is to make it easier to pipe
	arguments into functions.

Date Functions

	- date FORMAT TIME: Format a date, where a date is an integer type or a time.Time type, and
	  format is a time.Format formatting string.
	- dateModify: Given a date, modify it with a duration: `date_modify "-1.5h" now`. If the duration doesn't
	parse, it returns the time unaltered. See `time.ParseDuration` for info on duration strings.
	- now: Current time.Time, for feeding into date-related functions.
	- htmlDate TIME: Format a date for use in the value field of an HTML "date" form element.
	- dateInZone FORMAT TIME TZ: Like date, but takes three arguments: format, timestamp,
	  timezone.
	- htmlDateInZone TIME TZ: Like htmlDate, but takes two arguments: timestamp,
	  timezone.

String Functions

	- abbrev: Truncate a string with ellipses. `abbrev 5 "hello world"` yields "he..."
	- abbrevboth: Abbreviate from both sides, yielding "...lo wo..."
	- trunc: Truncate a string (no suffix). `trunc 5 "Hello World"` yields "hello".
	- trim: strings.TrimSpace
	- trimAll: strings.Trim, but with the argument order reversed `trimAll "$" "$5.00"` or `"$5.00 | trimAll "$"`
	- trimSuffix: strings.TrimSuffix, but with the argument order reversed: `trimSuffix "-" "ends-with-"`
	- trimPrefix: strings.TrimPrefix, but with the argument order reversed `trimPrefix "$" "$5"`
	- upper: strings.ToUpper
	- lower: strings.ToLower
	- nospace: Remove all space characters from a string. `nospace "h e l l o"` becomes "hello"
	- title: strings.Title
	- untitle: Remove title casing
	- repeat: strings.Repeat, but with the arguments switched: `repeat count str`. (This simplifies common pipelines)
	- substr: Given string, start, and length, return a substr.
	- initials: Given a multi-word string, return the initials. `initials "Matt Butcher"` returns "MB"
	- randAlphaNum: Given a length, generate a random alphanumeric sequence
	- randAlpha: Given a length, generate an alphabetic string
	- randAscii: Given a length, generate a random ASCII string (symbols included)
	- randNumeric: Given a length, generate a string of digits.
	- wrap: Force a line wrap at the given width. `wrap 80 "imagine a longer string"`
	- wrapWith: Wrap a line at the given length, but using 'sep' instead of a newline. `wrapWith 50, "<br>", $html`
	- contains: strings.Contains, but with the arguments switched: `contains substr str`. (This simplifies common pipelines)
	- hasPrefix: strings.hasPrefix, but with the arguments switched
	- hasSuffix: strings.hasSuffix, but with the arguments switched
	- quote: Wrap string(s) in double quotation marks, escape the contents by adding '\' before '"'.
	- squote: Wrap string(s) in double quotation marks, does not escape content.
	- cat: Concatenate strings, separating them by spaces. `cat $a $b $c`.
	- indent: Indent a string using space characters. `indent 4 "foo\nbar"` produces "    foo\n    bar"
	- replace: Replace an old with a new in a string: `$name | replace " " "-"`
	- plural: Choose singular or plural based on length: `len $fish | plural "one anchovy" "many anchovies"`
	- sha256sum: Generate a hex encoded sha256 hash of the input

String Slice Functions:

	- join: strings.Join, but as `join SEP SLICE`
	- split: strings.Split, but as `split SEP STRING`. The results are returned
	  as a map with the indexes set to _N, where N is an integer starting from 0.
	  Use it like this: `{{$v := "foo/bar/baz" | split "/"}}{{$v._0}}` (Prints `foo`)

Integer Slice Functions:

	- until: Given an integer, returns a slice of counting integers from 0 to one
	  less than the given integer: `range $i, $e := until 5`
	- untilStep: Given start, stop, and step, return an integer slice starting at
	  'start', stopping at `stop`, and incrementing by 'step. This is the same
	  as Python's long-form of 'range'.

Conversions:

	- atoi: Convert a string to an integer. 0 if the integer could not be parsed.
	- in64: Convert a string or another numeric type to an int64.
	- int: Convert a string or another numeric type to an int.
	- float64: Convert a string or another numeric type to a float64.

Defaults:

	- default: Give a default value. Used like this: trim "   "| default "empty".
	  Since trim produces an empty string, the default value is returned. For
	  things with a length (strings, slices, maps), len(0) will trigger the default.
	  For numbers, the value 0 will trigger the default. For booleans, false will
	  trigger the default. For structs, the default is never returned (there is
	  no clear empty condition). For everything else, nil value triggers a default.
	- empty: Return true if the given value is the zero value for its type.
	  Caveats: structs are always non-empty. This should match the behavior of
	  {{if pipeline}}, but can be used inside of a pipeline.

OS:
	- env: Resolve an environment variable
	- expandenv: Expand a string through the environment

File Paths:
	- base: Return the last element of a path. https://golang.org/pkg/path#Base
	- dir: Remove the last element of a path. https://golang.org/pkg/path#Dir
	- clean: Clean a path to the shortest equivalent name.  (e.g. remove "foo/.."
	from "foo/../bar.html") https://golang.org/pkg/path#Clean
	- ext: https://golang.org/pkg/path#Ext
	- isAbs: https://golang.org/pkg/path#IsAbs

Encoding:
	- b64enc: Base 64 encode a string.
	- b64dec: Base 64 decode a string.

Reflection:

	- typeOf: Takes an interface and returns a string representation of the type.
	  For pointers, this will return a type prefixed with an asterisk(`*`). So
	  a pointer to type `Foo` will be `*Foo`.
	- typeIs: Compares an interface with a string name, and returns true if they match.
	  Note that a pointer will not match a reference. For example `*Foo` will not
	  match `Foo`.
	- typeIsLike: Compares an interface with a string name and returns true if
	  the interface is that `name` or that `*name`. In other words, if the given
	  value matches the given type or is a pointer to the given type, this returns
	  true.
	- kindOf: Takes an interface and returns a string representation of its kind.
	- kindIs: Returns true if the given string matches the kind of the given interface.

	Note: None of these can test whether or not something implements a given
	interface, since doing so would require compiling the interface in ahead of
	time.

Data Structures:

	- tuple: Takes an arbitrary list of items and returns a slice of items. Its
	  tuple-ish properties are mainly gained through the template idiom, and not
	  through an API provided here.
	- dict: Takes a list of name/values and returns a map[string]interface{}.
	  The first parameter is converted to a string and stored as a key, the
	  second parameter is treated as the value. And so on, with odds as keys and
	  evens as values. If the function call ends with an odd, the last key will
	  be assigned the empty string. Non-string keys are converted to strings as
	  follows: []byte are converted, fmt.Stringers will have String() called.
	  errors will have Error() called. All others will be passed through
	  fmt.Sprtinf("%v").
	- set: Takes a dict, a key, and a value, and sets that key/value pair in
	  the dict. `set $dict $key $value`. For convenience, it returns the dict,
	  even though the dict was modified in place.
	- unset: Takes a dict and a key, and deletes that key/value pair from the
	  dict. `unset $dict $key`. This returns the dict for convenience.
	- hasKey: Takes a dict and a key, and returns boolean true if the key is in
	  the dict.

Math Functions:

Integer functions will convert integers of any width to `int64`. If a
string is passed in, functions will attempt to convert with
`strconv.ParseInt(s, 1064)`. If this fails, the value will be treated as 0.

	- add1: Increment an integer by 1
	- add: Sum an arbitrary number of integers
	- sub: Subtract the second integer from the first
	- div: Divide the first integer by the second
	- mod: Module of first integer divided by second
	- mul: Multiply integers
	- max: Return the biggest of a series of one or more integers
	- min: Return the smallest of a series of one or more integers
	- biggest: DEPRECATED. Return the biggest of a series of one or more integers

Crypto Functions:

	- genPrivateKey: Generate a private key for the given cryptosystem. If no
	  argument is supplied, by default it will generate a private key using
	  the RSA algorithm. Accepted values are `rsa`, `dsa`, and `ecdsa`.

*/
package sprig

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"html/template"
	"math"
	"math/big"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	ttemplate "text/template"
	"time"

	util "github.com/aokoli/goutils"
	uuid "github.com/satori/go.uuid"
)

// Produce the function map.
//
// Use this to pass the functions into the template engine:
//
// 	tpl := template.New("foo").Funcs(sprig.FuncMap))
//
func FuncMap() template.FuncMap {
	return HtmlFuncMap()
}

// HermeticTextFuncMap returns a 'text/template'.FuncMap with only repeatable functions.
func HermeticTxtFuncMap() ttemplate.FuncMap {
	r := TxtFuncMap()
	for _, name := range nonhermeticFunctions {
		delete(r, name)
	}
	return r
}

// HermeticHtmlFuncMap returns an 'html/template'.Funcmap with only repeatable functions.
func HermeticHtmlFuncMap() template.FuncMap {
	r := HtmlFuncMap()
	for _, name := range nonhermeticFunctions {
		delete(r, name)
	}
	return r
}

// TextFuncMap returns a 'text/template'.FuncMap
func TxtFuncMap() ttemplate.FuncMap {
	return ttemplate.FuncMap(genericMap)
}

// HtmlFuncMap returns an 'html/template'.Funcmap
func HtmlFuncMap() template.FuncMap {
	return template.FuncMap(genericMap)
}

// These functions are not guaranteed to evaluate to the same result for given input, because they
// refer to the environemnt or global state.
var nonhermeticFunctions = []string{
	// Date functions
	"date",
	"date_in_zone",
	"date_modify",
	"now",
	"htmlDate",
	"htmlDateInZone",
	"dateInZone",
	"dateModify",

	// Strings
	"randAlphaNum",
	"randAlpha",
	"randAscii",
	"randNumeric",
	"uuidv4",

	// OS
	"env",
	"expandenv",
}

var genericMap = map[string]interface{}{
	"hello": func() string { return "Hello!" },

	// Date functions
	"date":           date,
	"date_in_zone":   dateInZone,
	"date_modify":    dateModify,
	"now":            func() time.Time { return time.Now() },
	"htmlDate":       htmlDate,
	"htmlDateInZone": htmlDateInZone,
	"dateInZone":     dateInZone,
	"dateModify":     dateModify,

	// Strings
	"abbrev":     abbrev,
	"abbrevboth": abbrevboth,
	"trunc":      trunc,
	"trim":       strings.TrimSpace,
	"upper":      strings.ToUpper,
	"lower":      strings.ToLower,
	"title":      strings.Title,
	"untitle":    untitle,
	"substr":     substring,
	// Switch order so that "foo" | repeat 5
	"repeat": func(count int, str string) string { return strings.Repeat(str, count) },
	// Deprecated: Use trimAll.
	"trimall": func(a, b string) string { return strings.Trim(b, a) },
	// Switch order so that "$foo" | trimall "$"
	"trimAll":      func(a, b string) string { return strings.Trim(b, a) },
	"trimSuffix":   func(a, b string) string { return strings.TrimSuffix(b, a) },
	"trimPrefix":   func(a, b string) string { return strings.TrimPrefix(b, a) },
	"nospace":      util.DeleteWhiteSpace,
	"initials":     initials,
	"randAlphaNum": randAlphaNumeric,
	"randAlpha":    randAlpha,
	"randAscii":    randAscii,
	"randNumeric":  randNumeric,
	"swapcase":     util.SwapCase,
	"wrap":         func(l int, s string) string { return util.Wrap(s, l) },
	"wrapWith":     func(l int, sep, str string) string { return util.WrapCustom(str, l, sep, true) },
	// Switch order so that "foobar" | contains "foo"
	"contains":  func(substr string, str string) bool { return strings.Contains(str, substr) },
	"hasPrefix": func(substr string, str string) bool { return strings.HasPrefix(str, substr) },
	"hasSuffix": func(substr string, str string) bool { return strings.HasSuffix(str, substr) },
	"quote":     quote,
	"squote":    squote,
	"cat":       cat,
	"indent":    indent,
	"replace":   replace,
	"plural":    plural,
	"sha256sum": sha256sum,

	// Wrap Atoi to stop errors.
	"atoi":    func(a string) int { i, _ := strconv.Atoi(a); return i },
	"int64":   toInt64,
	"int":     toInt,
	"float64": toFloat64,

	//"gt": func(a, b int) bool {return a > b},
	//"gte": func(a, b int) bool {return a >= b},
	//"lt": func(a, b int) bool {return a < b},
	//"lte": func(a, b int) bool {return a <= b},

	// split "/" foo/bar returns map[int]string{0: foo, 1: bar}
	"split": split,

	"until":     until,
	"untilStep": untilStep,

	// VERY basic arithmetic.
	"add1": func(i interface{}) int64 { return toInt64(i) + 1 },
	"add": func(i ...interface{}) int64 {
		var a int64 = 0
		for _, b := range i {
			a += toInt64(b)
		}
		return a
	},
	"sub": func(a, b interface{}) int64 { return toInt64(a) - toInt64(b) },
	"div": func(a, b interface{}) int64 { return toInt64(a) / toInt64(b) },
	"mod": func(a, b interface{}) int64 { return toInt64(a) % toInt64(b) },
	"mul": func(a interface{}, v ...interface{}) int64 {
		val := toInt64(a)
		for _, b := range v {
			val = val * toInt64(b)
		}
		return val
	},
	"biggest": max,
	"max":     max,
	"min":     min,

	// string slices. Note that we reverse the order b/c that's better
	// for template processing.
	"join": func(sep string, ss []string) string { return strings.Join(ss, sep) },

	// Defaults
	"default": dfault,
	"empty":   empty,

	// Reflection
	"typeOf":     typeOf,
	"typeIs":     typeIs,
	"typeIsLike": typeIsLike,
	"kindOf":     kindOf,
	"kindIs":     kindIs,

	// OS:
	"env":       func(s string) string { return os.Getenv(s) },
	"expandenv": func(s string) string { return os.ExpandEnv(s) },

	// File Paths:
	"base":  path.Base,
	"dir":   path.Dir,
	"clean": path.Clean,
	"ext":   path.Ext,
	"isAbs": path.IsAbs,

	// Encoding:
	"b64enc": base64encode,
	"b64dec": base64decode,
	"b32enc": base32encode,
	"b32dec": base32decode,

	// Data Structures:
	"tuple":  tuple,
	"dict":   dict,
	"set":    set,
	"unset":  unset,
	"hasKey": hasKey,

	// Crypto:
	"genPrivateKey": generatePrivateKey,

	// UUIDs:
	"uuidv4": uuidv4,
}

func split(sep, orig string) map[string]string {
	parts := strings.Split(orig, sep)
	res := make(map[string]string, len(parts))
	for i, v := range parts {
		res["_"+strconv.Itoa(i)] = v
	}
	return res
}

// substring creates a substring of the given string.
//
// If start is < 0, this calls string[:length].
//
// If start is >= 0 and length < 0, this calls string[start:]
//
// Otherwise, this calls string[start, length].
func substring(start, length int, s string) string {
	if start < 0 {
		return s[:length]
	}
	if length < 0 {
		return s[start:]
	}
	return s[start:length]
}

// Given a format and a date, format the date string.
//
// Date can be a `time.Time` or an `int, int32, int64`.
// In the later case, it is treated as seconds since UNIX
// epoch.
func date(fmt string, date interface{}) string {
	return dateInZone(fmt, date, "Local")
}

func htmlDate(date interface{}) string {
	return dateInZone("2006-01-02", date, "Local")
}

func htmlDateInZone(date interface{}, zone string) string {
	return dateInZone("2006-01-02", date, zone)
}

func dateInZone(fmt string, date interface{}, zone string) string {
	var t time.Time
	switch date := date.(type) {
	default:
		t = time.Now()
	case time.Time:
		t = date
	case int64:
		t = time.Unix(date, 0)
	case int:
		t = time.Unix(int64(date), 0)
	case int32:
		t = time.Unix(int64(date), 0)
	}

	loc, err := time.LoadLocation(zone)
	if err != nil {
		loc, _ = time.LoadLocation("UTC")
	}

	return t.In(loc).Format(fmt)
}

func dateModify(fmt string, date time.Time) time.Time {
	d, err := time.ParseDuration(fmt)
	if err != nil {
		return date
	}
	return date.Add(d)
}

func max(a interface{}, i ...interface{}) int64 {
	aa := toInt64(a)
	for _, b := range i {
		bb := toInt64(b)
		if bb > aa {
			aa = bb
		}
	}
	return aa
}

func min(a interface{}, i ...interface{}) int64 {
	aa := toInt64(a)
	for _, b := range i {
		bb := toInt64(b)
		if bb < aa {
			aa = bb
		}
	}
	return aa
}

// dfault checks whether `given` is set, and returns default if not set.
//
// This returns `d` if `given` appears not to be set, and `given` otherwise.
//
// For numeric types 0 is unset.
// For strings, maps, arrays, and slices, len() = 0 is considered unset.
// For bool, false is unset.
// Structs are never considered unset.
//
// For everything else, including pointers, a nil value is unset.
func dfault(d interface{}, given ...interface{}) interface{} {

	if empty(given) || empty(given[0]) {
		return d
	}
	return given[0]
}

// empty returns true if the given value has the zero value for its type.
func empty(given interface{}) bool {
	g := reflect.ValueOf(given)
	if !g.IsValid() {
		return true
	}

	// Basically adapted from text/template.isTrue
	switch g.Kind() {
	default:
		return g.IsNil()
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return g.Len() == 0
	case reflect.Bool:
		return g.Bool() == false
	case reflect.Complex64, reflect.Complex128:
		return g.Complex() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return g.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return g.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return g.Float() == 0
	case reflect.Struct:
		return false
	}
	return true
}

// typeIs returns true if the src is the type named in target.
func typeIs(target string, src interface{}) bool {
	return target == typeOf(src)
}

func typeIsLike(target string, src interface{}) bool {
	t := typeOf(src)
	return target == t || "*"+target == t
}

func typeOf(src interface{}) string {
	return fmt.Sprintf("%T", src)
}

func kindIs(target string, src interface{}) bool {
	return target == kindOf(src)
}

func kindOf(src interface{}) string {
	return reflect.ValueOf(src).Kind().String()
}

func base64encode(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func base64decode(v string) string {
	data, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func base32encode(v string) string {
	return base32.StdEncoding.EncodeToString([]byte(v))
}

func base32decode(v string) string {
	data, err := base32.StdEncoding.DecodeString(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func abbrev(width int, s string) string {
	if width < 4 {
		return s
	}
	r, _ := util.Abbreviate(s, width)
	return r
}

func abbrevboth(left, right int, s string) string {
	if right < 4 || left > 0 && right < 7 {
		return s
	}
	r, _ := util.AbbreviateFull(s, left, right)
	return r
}
func initials(s string) string {
	// Wrap this just to eliminate the var args, which templates don't do well.
	return util.Initials(s)
}

func randAlphaNumeric(count int) string {
	// It is not possible, it appears, to actually generate an error here.
	r, _ := util.RandomAlphaNumeric(count)
	return r
}

func randAlpha(count int) string {
	r, _ := util.RandomAlphabetic(count)
	return r
}

func randAscii(count int) string {
	r, _ := util.RandomAscii(count)
	return r
}

func randNumeric(count int) string {
	r, _ := util.RandomNumeric(count)
	return r
}

func untitle(str string) string {
	return util.Uncapitalize(str)
}

func quote(str ...interface{}) string {
	out := make([]string, len(str))
	for i, s := range str {
		out[i] = fmt.Sprintf("%q", strval(s))
	}
	return strings.Join(out, " ")
}

func squote(str ...interface{}) string {
	out := make([]string, len(str))
	for i, s := range str {
		out[i] = fmt.Sprintf("'%v'", s)
	}
	return strings.Join(out, " ")
}

func tuple(v ...interface{}) []interface{} {
	return v
}

func set(d map[string]interface{}, key string, value interface{}) map[string]interface{} {
	d[key] = value
	return d
}

func unset(d map[string]interface{}, key string) map[string]interface{} {
	delete(d, key)
	return d
}

func hasKey(d map[string]interface{}, key string) bool {
	_, ok := d[key]
	return ok
}

func dict(v ...interface{}) map[string]interface{} {
	dict := map[string]interface{}{}
	lenv := len(v)
	for i := 0; i < lenv; i += 2 {
		key := strval(v[i])
		if i+1 >= lenv {
			dict[key] = ""
			continue
		}
		dict[key] = v[i+1]
	}
	return dict
}

func strval(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// toFloat64 converts 64-bit floats
func toFloat64(v interface{}) float64 {
	if str, ok := v.(string); ok {
		iv, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0
		}
		return iv
	}

	val := reflect.Indirect(reflect.ValueOf(v))
	switch val.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return float64(val.Int())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return float64(val.Uint())
	case reflect.Uint, reflect.Uint64:
		return float64(val.Uint())
	case reflect.Float32, reflect.Float64:
		return val.Float()
	case reflect.Bool:
		if val.Bool() == true {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func toInt(v interface{}) int {
	//It's not optimal. Bud I don't want duplicate toInt64 code.
	return int(toInt64(v))
}

// toInt64 converts integer types to 64-bit integers
func toInt64(v interface{}) int64 {
	if str, ok := v.(string); ok {
		iv, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return 0
		}
		return iv
	}

	val := reflect.Indirect(reflect.ValueOf(v))
	switch val.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return val.Int()
	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return int64(val.Uint())
	case reflect.Uint, reflect.Uint64:
		tv := val.Uint()
		if tv <= math.MaxInt64 {
			return int64(tv)
		}
		// TODO: What is the sensible thing to do here?
		return math.MaxInt64
	case reflect.Float32, reflect.Float64:
		return int64(val.Float())
	case reflect.Bool:
		if val.Bool() == true {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func generatePrivateKey(typ string) string {
	var priv interface{}
	var err error
	switch typ {
	case "", "rsa":
		// good enough for government work
		priv, err = rsa.GenerateKey(rand.Reader, 4096)
	case "dsa":
		key := new(dsa.PrivateKey)
		// again, good enough for government work
		if err = dsa.GenerateParameters(&key.Parameters, rand.Reader, dsa.L2048N256); err != nil {
			return fmt.Sprintf("failed to generate dsa params: %s", err)
		}
		err = dsa.GenerateKey(key, rand.Reader)
		priv = key
	case "ecdsa":
		// again, good enough for government work
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	default:
		return "Unknown type " + typ
	}
	if err != nil {
		return fmt.Sprintf("failed to generate private key: %s", err)
	}

	return string(pem.EncodeToMemory(pemBlockForKey(priv)))
}

type DSAKeyFormat struct {
	Version       int
	P, Q, G, Y, X *big.Int
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *dsa.PrivateKey:
		val := DSAKeyFormat{
			P: k.P, Q: k.Q, G: k.G,
			Y: k.Y, X: k.X,
		}
		bytes, _ := asn1.Marshal(val)
		return &pem.Block{Type: "DSA PRIVATE KEY", Bytes: bytes}
	case *ecdsa.PrivateKey:
		b, _ := x509.MarshalECPrivateKey(k)
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func trunc(c int, s string) string {
	if len(s) <= c {
		return s
	}
	return s[0:c]
}

func cat(v ...interface{}) string {
	r := strings.TrimSpace(strings.Repeat("%v ", len(v)))
	return fmt.Sprintf(r, v...)
}

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

func replace(old, new, src string) string {
	return strings.Replace(src, old, new, -1)
}

func plural(one, many string, count int) string {
	if count == 1 {
		return one
	}
	return many
}

func sha256sum(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func until(count int) []int {
	step := 1
	if count < 0 {
		step = -1
	}
	return untilStep(0, count, step)
}

func untilStep(start, stop, step int) []int {
	v := []int{}

	if stop < start {
		if step >= 0 {
			return v
		}
		for i := start; i > stop; i += step {
			v = append(v, i)
		}
		return v
	}

	if step <= 0 {
		return v
	}
	for i := start; i < stop; i += step {
		v = append(v, i)
	}
	return v
}

// uuidv4 provides a safe and secure UUID v4 implementation
func uuidv4() string {
	return fmt.Sprintf("%s", uuid.NewV4())
}
