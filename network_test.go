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

func TestGetHostByAddr(t *testing.T) {
	tpl := `{{"1.1.1.1" | getHostByAddr}}`

	resolvedHost, _ := runRaw(tpl, nil)

	assert.Equal(t, resolvedHost, "one.one.one.one.")
}
