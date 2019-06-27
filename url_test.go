package sprig

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlParse(t *testing.T) {
	tests := map[string]string{
		`{{ urlParse "proto://auth@host:80/path?query#fragment" }}`:
			"map[fragment:fragment host:host:80 hostname:host opaque: path:/path query:query scheme:proto userinfo:auth]",
		`{{ urlParse "proto://host:80/path" }}`:
			"map[fragment: host:host:80 hostname:host opaque: path:/path query: scheme:proto userinfo:]",
		`{{ urlParse "something" }}`:
			"map[fragment: host: hostname: opaque: path:something query: scheme: userinfo:]",
		`{{ urlParse "proto://user:passwor%20d@host:80/path" }}`:
			"map[fragment: host:host:80 hostname:host opaque: path:/path query: scheme:proto userinfo:user:passwor%20d]",
		`{{ urlParse "proto://@host:80/pa%20th" }}`:
			"map[fragment: host:host:80 hostname:host opaque: path:/pa th query: scheme:proto userinfo:]",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestUrlJoin(t *testing.T) {
	tests := map[string]string{
		`{{ urlJoin (dict "fragment" "fragment" "host" "host:80" "path" "/path" "query" "query" "scheme" "proto") }}`:
			"proto://host:80/path?query#fragment",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}
