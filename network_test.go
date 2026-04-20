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
		tpl      string
		expected string
	}{
		{`{{joinHostPort "192.0.2.1" "80"}}`, "192.0.2.1:80"},
		{`{{joinHostPort "2001:db8::1" "80"}}`, "[2001:db8::1]:80"},
		{`{{joinHostPort "example.com" "443"}}`, "example.com:443"},
	}
	for _, tt := range tests {
		result, err := runRaw(tt.tpl, nil)
		assert.NoError(t, err)
		assert.Equal(t, tt.expected, result)
	}
}

func TestSplitHostPort(t *testing.T) {
	tests := []struct {
		tpl          string
		expectedHost string
		expectedPort string
	}{
		{`{{with splitHostPort "192.0.2.1:80"}}{{.host}}/{{.port}}{{end}}`, "192.0.2.1", "80"},
		{`{{with splitHostPort "[2001:db8::1]:80"}}{{.host}}/{{.port}}{{end}}`, "2001:db8::1", "80"},
		{`{{with splitHostPort "example.com:443"}}{{.host}}/{{.port}}{{end}}`, "example.com", "443"},
	}
	for _, tt := range tests {
		result, err := runRaw(tt.tpl, nil)
		assert.NoError(t, err)
		assert.Equal(t, tt.expectedHost+"/"+tt.expectedPort, result)
	}
}
