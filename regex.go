package sprig

import (
	"regexp"
)

func regexMatch(regex string, s string) bool {
	match, _ := regexp.MatchString(regex, s)
	return match
}

func mustRegexMatch(regex string, s string) (bool, error) {
	return regexp.MatchString(regex, s)
}

func regexFindAll(regex string, s string, n int) []string {
	r := regexp.MustCompile(regex)
	return r.FindAllString(s, n)
}

func mustRegexFindAll(regex string, s string, n int) ([]string, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return []string{}, err
	}
	return r.FindAllString(s, n), nil
}

// regexFindAllGlobal is a pipeline-friendly variant of regexFindAll: it takes
// the subject string as its last argument, so it can be used at the end of a
// template pipeline, e.g. `$s | regexFindAllGlobal "a{2}" -1`.
func regexFindAllGlobal(regex string, n int, s string) []string {
	return regexFindAll(regex, s, n)
}

func mustRegexFindAllGlobal(regex string, n int, s string) ([]string, error) {
	return mustRegexFindAll(regex, s, n)
}

func regexFind(regex string, s string) string {
	r := regexp.MustCompile(regex)
	return r.FindString(s)
}

func mustRegexFind(regex string, s string) (string, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return "", err
	}
	return r.FindString(s), nil
}

func regexReplaceAll(regex string, s string, repl string) string {
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, repl)
}

func mustRegexReplaceAll(regex string, s string, repl string) (string, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return "", err
	}
	return r.ReplaceAllString(s, repl), nil
}

// regexReplaceAllGlobal is a pipeline-friendly variant of regexReplaceAll: it
// takes the subject string as its last argument, so it can be used at the end
// of a template pipeline, e.g. `$s | regexReplaceAllGlobal "a(x*)b" "${1}W"`.
func regexReplaceAllGlobal(regex string, repl string, s string) string {
	return regexReplaceAll(regex, s, repl)
}

func mustRegexReplaceAllGlobal(regex string, repl string, s string) (string, error) {
	return mustRegexReplaceAll(regex, s, repl)
}

func regexReplaceAllLiteral(regex string, s string, repl string) string {
	r := regexp.MustCompile(regex)
	return r.ReplaceAllLiteralString(s, repl)
}

func mustRegexReplaceAllLiteral(regex string, s string, repl string) (string, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return "", err
	}
	return r.ReplaceAllLiteralString(s, repl), nil
}

// regexReplaceAllLiteralGlobal is a pipeline-friendly variant of
// regexReplaceAllLiteral: it takes the subject string as its last argument,
// so it can be used at the end of a template pipeline, e.g.
// `$s | regexReplaceAllLiteralGlobal "a(x*)b" "${1}"`.
func regexReplaceAllLiteralGlobal(regex string, repl string, s string) string {
	return regexReplaceAllLiteral(regex, s, repl)
}

func mustRegexReplaceAllLiteralGlobal(regex string, repl string, s string) (string, error) {
	return mustRegexReplaceAllLiteral(regex, s, repl)
}

func regexSplit(regex string, s string, n int) []string {
	r := regexp.MustCompile(regex)
	return r.Split(s, n)
}

func mustRegexSplit(regex string, s string, n int) ([]string, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return []string{}, err
	}
	return r.Split(s, n), nil
}

// regexSplitGlobal is a pipeline-friendly variant of regexSplit: it takes the
// subject string as its last argument, so it can be used at the end of a
// template pipeline, e.g. `$s | regexSplitGlobal "z+" -1`.
func regexSplitGlobal(regex string, n int, s string) []string {
	return regexSplit(regex, s, n)
}

func mustRegexSplitGlobal(regex string, n int, s string) ([]string, error) {
	return mustRegexSplit(regex, s, n)
}

func regexQuoteMeta(s string) string {
	return regexp.QuoteMeta(s)
}
