package sprig

import (
	"testing"
)

func TestIssue188(t *testing.T) {
	tests := map[string]string{

		// This first test shows two merges and the merge is NOT A DEEP COPY MERGE.
		// The first merge puts $one on to $target. When the second merge of $two
		// on to $target the nested dict brought over from $one is changed on
		// $one as well as $target.
		`{{- $target := dict -}}
			{{- $one := dict "foo" (dict "bar" "baz") "qux" true -}}
			{{- $two := dict "foo" (dict "bar" "baz2") "qux" false -}}
			{{- mergeOverwrite $target $one | toString | trunc 0 }}{{ $__ := mergeOverwrite $target $two }}{{ $one }}`: "map[foo:map[bar:baz2] qux:true]",

		// This test uses deepCopy on $one to create a deep copy and then merge
		// that. In this case the merge of $two on to $target does not affect
		// $one because a deep copy was used for that merge.
		`{{- $target := dict -}}
			{{- $one := dict "foo" (dict "bar" "baz") "qux" true -}}
			{{- $two := dict "foo" (dict "bar" "baz2") "qux" false -}}
			{{- deepCopy $one | mergeOverwrite $target | toString | trunc 0 }}{{ $__ := mergeOverwrite $target $two }}{{ $one }}`: "map[foo:map[bar:baz] qux:true]",
	}

	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
	}
}
