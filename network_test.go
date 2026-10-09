package sprig

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostByName(t *testing.T) {
	tpl := `{{"www.google.com" | getHostByName}}`

	resolvedIP, _ := runRaw(tpl, nil)

	ip := net.ParseIP(resolvedIP)
	assert.NotNil(t, ip)
	assert.NotEmpty(t, ip)
}

func TestJoinHostPort(t *testing.T) {
	tests := []struct {
		name     string
		host     interface{}
		port     interface{}
		expected string
	}{
		{
			name:     "IPv4 address with integer port",
			host:     "127.0.0.1",
			port:     80,
			expected: "127.0.0.1:80",
		},
		{
			name:     "IPv4 address with string port",
			host:     "192.168.1.1",
			port:     "8080",
			expected: "192.168.1.1:8080",
		},
		{
			name:     "IPv4 address with int64 port",
			host:     "10.0.0.1",
			port:     int64(443),
			expected: "10.0.0.1:443",
		},
		{
			name:     "IPv6 address with integer port",
			host:     "::1",
			port:     80,
			expected: "[::1]:80",
		},
		{
			name:     "IPv6 address with string port",
			host:     "::1",
			port:     "8080",
			expected: "[::1]:8080",
		},
		{
			name:     "IPv6 full address with integer port",
			host:     "2001:db8::1",
			port:     80,
			expected: "[2001:db8::1]:80",
		},
		{
			name:     "IPv6 full address with named service port",
			host:     "2001:db8::1",
			port:     "https",
			expected: "[2001:db8::1]:https",
		},
		{
			name:     "IPv6 scoped link-local address",
			host:     "fe80::1%eth0",
			port:     53,
			expected: "[fe80::1%eth0]:53",
		},
		{
			name:     "localhost with integer port",
			host:     "localhost",
			port:     80,
			expected: "localhost:80",
		},
		{
			name:     "localhost with service name port",
			host:     "localhost",
			port:     "http",
			expected: "localhost:http",
		},
		{
			name:     "domain name with integer port",
			host:     "example.com",
			port:     8080,
			expected: "example.com:8080",
		},
		{
			name:     "empty host with integer port",
			host:     "",
			port:     80,
			expected: ":80",
		},
		{
			name:     "empty host with string port",
			host:     "",
			port:     "8080",
			expected: ":8080",
		},
		{
			name:     "empty host and empty port",
			host:     "",
			port:     "",
			expected: ":",
		},
		{
			name:     "host with empty port",
			host:     "localhost",
			port:     "",
			expected: "localhost:",
		},
		{
			name:     "IPv6 host with empty port",
			host:     "::1",
			port:     "",
			expected: "[::1]:",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := joinHostPort(tc.host, tc.port)
			assert.Equal(t, tc.expected, actual)
		})
	}

	templateTests := map[string]string{
		`{{ joinHostPort "127.0.0.1" 80 }}`:    "127.0.0.1:80",
		`{{ joinHostPort "::1" 8080 }}`:        "[::1]:8080",
		`{{ joinHostPort "localhost" "http" }}`: "localhost:http",
		`{{ joinHostPort "example.com" 443 }}`:  "example.com:443",
		`{{ joinHostPort "" 8080 }}`:           ":8080",
		`{{ joinHostPort "" "" }}`:             ":",
		`{{ 80 | joinHostPort "127.0.0.1" }}`:  "127.0.0.1:80",
	}

	for tpl, expected := range templateTests {
		assert.NoError(t, runt(tpl, expected))
	}

	assert.NoError(t, runtv(
		`{{ joinHostPort .host .port }}`,
		"[::1]:9090",
		map[string]interface{}{"host": "::1", "port": 9090},
	))
}
