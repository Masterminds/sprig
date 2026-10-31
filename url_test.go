package sprig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var urlTests = map[string]map[string]interface{}{
	"proto://auth@host:80/path?query#fragment": {
		"fragment": "fragment",
		"host":     "host:80",
		"hostname": "host",
		"opaque":   "",
		"path":     "/path",
		"query":    "query",
		"scheme":   "proto",
		"userinfo": "auth",
	},
	"proto://host:80/path": {
		"fragment": "",
		"host":     "host:80",
		"hostname": "host",
		"opaque":   "",
		"path":     "/path",
		"query":    "",
		"scheme":   "proto",
		"userinfo": "",
	},
	"something": {
		"fragment": "",
		"host":     "",
		"hostname": "",
		"opaque":   "",
		"path":     "something",
		"query":    "",
		"scheme":   "",
		"userinfo": "",
	},
	"proto://user:passwor%20d@host:80/path": {
		"fragment": "",
		"host":     "host:80",
		"hostname": "host",
		"opaque":   "",
		"path":     "/path",
		"query":    "",
		"scheme":   "proto",
		"userinfo": "user:passwor%20d",
	},
	"proto://host:80/pa%20th?key=val%20ue": {
		"fragment": "",
		"host":     "host:80",
		"hostname": "host",
		"opaque":   "",
		"path":     "/pa th",
		"query":    "key=val%20ue",
		"scheme":   "proto",
		"userinfo": "",
	},
}

func TestUrlParse(t *testing.T) {
	// testing that function is exported and working properly
	assert.NoError(t, runt(
		`{{ index ( urlParse "proto://auth@host:80/path?query#fragment" ) "host" }}`,
		"host:80"))

	// testing scenarios
	for url, expected := range urlTests {
		assert.EqualValues(t, expected, urlParse(url))
	}
}

func TestUrlJoin(t *testing.T) {
	tests := map[string]string{
		`{{ urlJoin (dict "fragment" "fragment" "host" "host:80" "path" "/path" "query" "query" "scheme" "proto") }}`:       "proto://host:80/path?query#fragment",
		`{{ urlJoin (dict "fragment" "fragment" "host" "host:80" "path" "/path" "scheme" "proto" "userinfo" "ASDJKJSD") }}`: "proto://ASDJKJSD@host:80/path#fragment",
	}
	for tpl, expected := range tests {
		assert.NoError(t, runt(tpl, expected))
	}

	for expected, urlMap := range urlTests {
		assert.EqualValues(t, expected, urlJoin(urlMap))
	}

}

func TestUrlParseJoinEscapedPath(t *testing.T) {
	for _, input := range []string{
		"https://example.com/a%2Fb",
		"https://example.com/a%2fb",
		"https://example.com/%252F/a%2Fb?key=a%2Fb#fragment",
		"/a%2Fb",
	} {
		t.Run(input, func(t *testing.T) {
			parts := urlParse(input)
			assert.Equal(t, input, urlJoin(parts))
			assert.NoError(t, runtv(`{{ urlParse . | urlJoin }}`, input, input))
		})
	}
}

func TestUrlJoinUpdatedEscapedPath(t *testing.T) {
	parts := urlParse("https://example.com/a%2Fb")
	parts["path"] = "/new path"
	assert.Equal(t, "https://example.com/new%20path", urlJoin(parts))
}
