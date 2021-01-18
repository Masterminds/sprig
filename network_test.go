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

	tpl = `{{"nope.local" | getHostByName}}`

	resolvedIP, _ = runRaw(tpl, nil)

	assert.Empty(t, resolvedIP)
}
