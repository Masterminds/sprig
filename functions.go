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
	- trim: strings.TrimSpace
	- trimall: strings.Trim, but with the argument order reversed `trimall "$" "$5.00"` or `"$5.00 | trimall "$"`
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
	- quote: Wrap string(s) in double quotation marks.
	- squote: Wrap string(s) in double quotation marks.

String Slice Functions:

	- join: strings.Join, but as `join SEP SLICE`
	- split: strings.Split, but as `split SEP STRING`. The results are returned
	  as a map with the indexes set to _N, where N is an integer starting from 0.
	  Use it like this: `{{$v := "foo/bar/baz" | split "/"}}{{$v._0}}` (Prints `foo`)

Conversions:

	- atoi: Convert a string to an integer. 0 if the integer could not be parsed.
	- toInt64: Convert a string or another numeric type to an int64.

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

REMOVED (implemented in Go 1.2)

	- gt: Greater than (integer)
	- lt: Less than (integer)
	- gte: Greater than or equal to (integer)
	- lte: Less than or equal to (integer)

*/
package sprig

import (
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"html/template"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	ttemplate "text/template"
	"time"

	util "github.com/aokoli/goutils"
)

// Produce the function map.
//
// Use this to pass the functions into the template engine:
//
// 	tpl := template.New("foo").Funcs(sprig.FuncMap))
//
func FuncMap() template.FuncMap {
	return template.FuncMap(genericMap)
}

// TextFuncMap returns a 'text/template'.FuncMap
func TxtFuncMap() ttemplate.FuncMap {
	return ttemplate.FuncMap(genericMap)
}

// HtmlFuncMap returns an 'html/template'.Funcmap
func HtmlFuncMap() template.FuncMap {
	return template.FuncMap(genericMap)
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
	"trim":       strings.TrimSpace,
	"upper":      strings.ToUpper,
	"lower":      strings.ToLower,
	"title":      strings.Title,
	"untitle":    untitle,
	"substr":     substring,
	// Switch order so that "foo" | repeat 5
	"repeat": func(count int, str string) string { return strings.Repeat(str, count) },
	// Switch order so that "$foo" | trimall "$"
	"trimall":      func(a, b string) string { return strings.Trim(b, a) },
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
	"contains": func(substr string, str string) bool { return strings.Contains(str, substr) },
	"quote":    quote,
	"squote":   squote,

	// Wrap Atoi to stop errors.
	"atoi":  func(a string) int { i, _ := strconv.Atoi(a); return i },
	"int64": toInt64,

	//"gt": func(a, b int) bool {return a > b},
	//"gte": func(a, b int) bool {return a >= b},
	//"lt": func(a, b int) bool {return a < b},
	//"lte": func(a, b int) bool {return a <= b},

	// split "/" foo/bar returns map[int]string{0: foo, 1: bar}
	"split": split,

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

	// Encoding:
	"b64enc": base64encode,
	"b64dec": base64decode,
	"b32enc": base32encode,
	"b32dec": base32decode,

	// Data Structures:
	"tuple": tuple,
	"dict":  dict,
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

func quote(str ...string) string {
	for i, s := range str {
		str[i] = fmt.Sprintf("%q", s)
	}
	return strings.Join(str, " ")
}

func squote(str ...string) string {
	for i, s := range str {
		str[i] = fmt.Sprintf("'%s'", s)
	}
	return strings.Join(str, " ")
}

func tuple(v ...interface{}) []interface{} {
	return v
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
