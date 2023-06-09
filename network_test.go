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

func TestCidrNetmask(t *testing.T) {
	tpl := `{{"1.2.3.4/32" | cidrNetmask}}`
	netmask, _ := runRaw(tpl, nil)
	assert.NotNil(t, netmask)
	assert.NotEmpty(t, netmask)
	assert.Equal(t, "255.255.255.255", netmask)
}
