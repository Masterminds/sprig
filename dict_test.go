package sprig

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		`{{ keys (dict "foo" 1) (dict "bar" 2) (dict "bar" 3) | uniq | sortAlpha }}`: "[bar foo]",
	}
	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
	}
}

func TestPick(t *testing.T) {
	tests := map[string]string{
		`{{- $d := dict "one" 1 "two" 222222 }}{{ pick $d "two" | len -}}`:               "1",
		`{{- $d := dict "one" 1 "two" 222222 }}{{ pick $d "two" -}}`:                     "map[two:222222]",
		`{{- $d := dict "one" 1 "two" 222222 }}{{ pick $d "one" "two" | len -}}`:         "2",
		`{{- $d := dict "one" 1 "two" 222222 }}{{ pick $d "one" "two" "three" | len -}}`: "2",
		`{{- $d := dict }}{{ pick $d "two" | len -}}`:                                    "0",
	}
	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
	}
}
func TestOmit(t *testing.T) {
	tests := map[string]string{
		`{{- $d := dict "one" 1 "two" 222222 }}{{ omit $d "one" | len -}}`:         "1",
		`{{- $d := dict "one" 1 "two" 222222 }}{{ omit $d "one" -}}`:               "map[two:222222]",
		`{{- $d := dict "one" 1 "two" 222222 }}{{ omit $d "one" "two" | len -}}`:   "0",
		`{{- $d := dict "one" 1 "two" 222222 }}{{ omit $d "two" "three" | len -}}`: "1",
		`{{- $d := dict }}{{ omit $d "two" | len -}}`:                              "0",
	}
	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
	}
}

func TestGet(t *testing.T) {
	tests := map[string]string{
		`{{- $d := dict "one" 1 }}{{ get $d "one" -}}`:           "1",
		`{{- $d := dict "one" 1 "two" "2" }}{{ get $d "two" -}}`: "2",
		`{{- $d := dict }}{{ get $d "two" -}}`:                   "",
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

func TestMerge(t *testing.T) {
	dict := map[string]interface{}{
		"src2": map[string]interface{}{
			"h": 10,
			"i": "i",
			"j": "j",
		},
		"src1": map[string]interface{}{
			"a": 1,
			"b": 2,
			"d": map[string]interface{}{
				"e": "four",
			},
			"g": []int{6, 7},
			"i": "aye",
			"j": "jay",
			"k": map[string]interface{}{
				"l": false,
				"m": true,
			},
			"z": 10,
		},
		"dst": map[string]interface{}{
			"a": "one",
			"c": 3,
			"d": map[string]interface{}{
				"f": 5,
			},
			"g": []int{8, 9},
			"i": "eye",
			"j": nil,
			"k": map[string]interface{}{
				"l": true,
				"m": false,
			},
			"z": 0,
		},
	}
	tpl := `{{merge .dst .src1 .src2}}`
	_, err := runRaw(tpl, dict)
	if err != nil {
		t.Error(err)
	}
	expected := map[string]interface{}{
		"a": "one", // key overridden
		"b": 2,     // merged from src1
		"c": 3,     // merged from dst
		"d": map[string]interface{}{ // deep merge
			"e": "four",
			"f": 5,
		},
		"g": []int{8, 9}, // overridden - arrays are not merged
		"h": 10,          // merged from src2
		"i": "eye",       // overridden twice
		"j": "jay",       // overridden and merged
		"k": map[string]interface{}{
			"l": true, // overridden
			"m": false,
		},
		"z": 0, // zero value should be preserved
	}
	assert.Equal(t, expected, dict["dst"])
}

func TestMergeOverwrite(t *testing.T) {
	dict := map[string]interface{}{
		"src2": map[string]interface{}{
			"h": 10,
			"i": "i",
			"j": "j",
		},
		"src1": map[string]interface{}{
			"a": 1,
			"b": 2,
			"d": map[string]interface{}{
				"e": "four",
			},
			"g": []int{6, 7},
			"i": "aye",
			"j": "jay",
			"k": map[string]interface{}{
				"l": false,
			},
		},
		"dst": map[string]interface{}{
			"a": "one",
			"c": 3,
			"d": map[string]interface{}{
				"f": 5,
			},
			"g": []int{8, 9},
			"i": "eye",
			"k": map[string]interface{}{
				"l": true,
			},
		},
	}
	tpl := `{{mergeOverwrite .dst .src1 .src2}}`
	_, err := runRaw(tpl, dict)
	if err != nil {
		t.Error(err)
	}
	expected := map[string]interface{}{
		"a": 1, // key overwritten from src1
		"b": 2, // merged from src1
		"c": 3, // merged from dst
		"d": map[string]interface{}{ // deep merge
			"e": "four",
			"f": 5,
		},
		"g": []int{6, 7}, // overwritten src1 wins
		"h": 10,          // merged from src2
		"i": "i",         // overwritten twice src2 wins
		"j": "j",         // overwritten twice src2 wins
		"k": map[string]interface{}{ // deep merge
			"l": false, // overwritten src1 wins
		},
	}
	assert.Equal(t, expected, dict["dst"])
}

func TestValues(t *testing.T) {
	tests := map[string]string{
		`{{- $d := dict "a" 1 "b" 2 }}{{ values $d | sortAlpha | join "," }}`:       "1,2",
		`{{- $d := dict "a" "first" "b" 2 }}{{ values $d | sortAlpha | join "," }}`: "2,first",
	}

	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
	}
}

func TestDeepCopy(t *testing.T) {
	tests := map[string]string{
		`{{- $d := dict "a" 1 "b" 2 | deepCopy }}{{ values $d | sortAlpha | join "," }}`: "1,2",
		`{{- $d := dict "a" 1 "b" 2 | deepCopy }}{{ keys $d | sortAlpha | join "," }}`:   "a,b",
		`{{- $one := dict "foo" (dict "bar" "baz") "qux" true -}}{{ deepCopy $one }}`:    "map[foo:map[bar:baz] qux:true]",
	}

	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
	}
}

func TestMustDeepCopy(t *testing.T) {
	tests := map[string]string{
		`{{- $d := dict "a" 1 "b" 2 | mustDeepCopy }}{{ values $d | sortAlpha | join "," }}`: "1,2",
		`{{- $d := dict "a" 1 "b" 2 | mustDeepCopy }}{{ keys $d | sortAlpha | join "," }}`:   "a,b",
		`{{- $one := dict "foo" (dict "bar" "baz") "qux" true -}}{{ mustDeepCopy $one }}`:    "map[foo:map[bar:baz] qux:true]",
	}

	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
	}
}

func TestDig(t *testing.T) {
	tests := map[string]string{
		`{{- $d := dict "a" (dict "b" (dict "c" 1)) }}{{ dig "a" "b" "c" "" $d }}`:  "1",
		`{{- $d := dict "a" (dict "b" (dict "c" 1)) }}{{ dig "a" "b" "z" "2" $d }}`: "2",
		`{{ dict "a" 1 | dig "a" "" }}`:                                             "1",
		`{{ dict "a" 1 | dig "z" "2" }}`:                                            "2",
	}

	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
	}
}
