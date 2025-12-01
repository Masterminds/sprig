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
		"port":     "80",
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
		"port":     "80",
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
		"port":     "",
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
		"port":     "80",
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
		"port":     "80",
		"opaque":   "",
		"path":     "/pa th",
		"query":    "key=val%20ue",
		"scheme":   "proto",
		"userinfo": "",
	},
	"proto://[2001:db8::1]:8080/path": {
		"fragment": "",
		"host":     "[2001:db8::1]:8080",
		"hostname": "2001:db8::1",
		"port":     "8080",
		"opaque":   "",
		"path":     "/path",
		"query":    "",
		"scheme":   "proto",
		"userinfo": "",
	},
	"proto://[::1]:8080/path?query#fragment": {
		"fragment": "fragment",
		"host":     "[::1]:8080",
		"hostname": "::1",
		"port":     "8080",
		"opaque":   "",
		"path":     "/path",
		"query":    "query",
		"scheme":   "proto",
		"userinfo": "",
	},
	"proto://user@[2001:db8::1]:8080/path": {
		"fragment": "",
		"host":     "[2001:db8::1]:8080",
		"hostname": "2001:db8::1",
		"port":     "8080",
		"opaque":   "",
		"path":     "/path",
		"query":    "",
		"scheme":   "proto",
		"userinfo": "user",
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
		// Test IPv6 with hostname and port
		`{{ urlJoin (dict "hostname" "2001:db8::1" "port" "8080" "path" "/path" "scheme" "proto") }}`: "proto://[2001:db8::1]:8080/path",
		`{{ urlJoin (dict "hostname" "::1" "port" "8080" "path" "/path" "scheme" "proto") }}`:         "proto://[::1]:8080/path",
	}
	for tpl, expected := range tests {
		assert.NoError(t, runt(tpl, expected))
	}

	for expected, urlMap := range urlTests {
		assert.EqualValues(t, expected, urlJoin(urlMap))
	}

}
