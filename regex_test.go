package sprig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexFindAll(t *testing.T) {
	type args struct {
		regex, s string
		n        int
	}
	cases := []struct {
		expected int
		args     args
	}{
		{1, args{"a{2}", "aa", -1}},
		{1, args{"a{2}", "aaaaaaaa", 1}},
		{2, args{"a{2}", "aaaa", -1}},
		{0, args{"a{2}", "none", -1}},
	}

	for _, c := range cases {
		res, err := regexFindAll(c.args.regex, c.args.s, c.args.n)
		if err != nil {
			t.Errorf("regexFindAll test case %v failed with err %s", c, err)
		}
		assert.Equal(t, c.expected, len(res), "case %#v", c.args)
	}
}

func TestRegexFindl(t *testing.T) {
	type args struct{ regex, s string }
	cases := []struct {
		expected string
		args     args
	}{
		{"foo", args{"fo.?", "foorbar"}},
		{"foo", args{"fo.?", "foo foe fome"}},
		{"", args{"fo.?", "none"}},
	}

	for _, c := range cases {
		res, err := regexFind(c.args.regex, c.args.s)
		if err != nil {
			t.Errorf("regexFind test case %v failed with err %s", c, err)
		}
		assert.Equal(t, c.expected, res, "case %#v", c.args)
	}
}

func TestRegexReplaceAll(t *testing.T) {
	type args struct{ regex, s, repl string }
	cases := []struct {
		expected string
		args     args
	}{
		{"-T-T-", args{"a(x*)b", "-ab-axxb-", "T"}},
		{"--xx-", args{"a(x*)b", "-ab-axxb-", "$1"}},
		{"---", args{"a(x*)b", "-ab-axxb-", "$1W"}},
		{"-W-xxW-", args{"a(x*)b", "-ab-axxb-", "${1}W"}},
	}

	for _, c := range cases {
		res, err := regexReplaceAll(c.args.regex, c.args.s, c.args.repl)
		if err != nil {
			t.Errorf("regexReplaceAll test case %v failed with err %s", c, err)
		}
		assert.Equal(t, c.expected, res, "case %#v", c.args)
	}
}

func TestRegexReplaceAllLiteral(t *testing.T) {
	type args struct{ regex, s, repl string }
	cases := []struct {
		expected string
		args     args
	}{
		{"-T-T-", args{"a(x*)b", "-ab-axxb-", "T"}},
		{"-$1-$1-", args{"a(x*)b", "-ab-axxb-", "$1"}},
		{"-${1}-${1}-", args{"a(x*)b", "-ab-axxb-", "${1}"}},
	}

	for _, c := range cases {
		res, err := regexReplaceAllLiteral(c.args.regex, c.args.s, c.args.repl)
		if err != nil {
			t.Errorf("regexReplaceAllLiteral test case %v failed with err %s", c, err)
		}
		assert.Equal(t, c.expected, res, "case %#v", c.args)
	}
}

func TestRegexSplit(t *testing.T) {
	type args struct {
		regex, s string
		n        int
	}
	cases := []struct {
		expected int
		args     args
	}{
		{4, args{"a", "banana", -1}},
		{0, args{"a", "banana", 0}},
		{1, args{"a", "banana", 1}},
		{2, args{"a", "banana", 2}},
		{2, args{"z+", "pizza", -1}},
		{0, args{"z+", "pizza", 0}},
		{1, args{"z+", "pizza", 1}},
		{2, args{"z+", "pizza", 2}},
	}

	for _, c := range cases {
		res, err := regexSplit(c.args.regex, c.args.s, c.args.n)
		if err != nil {
			t.Errorf("regexSplit test case %v failed with err %s", c, err)
		}
		assert.Equal(t, c.expected, len(res), "case %#v", c.args)
	}
}
