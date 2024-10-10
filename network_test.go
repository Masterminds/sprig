package sprig

import (
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostByName(t *testing.T) {
	// GIVEN a valid hostname
	tpl := `{{"google.com" | getHostByName}}`

	// WHEN getHostByName is executed
	resolvedIP, err := runRaw(tpl, nil)

	// THEN the resolved IP should not be empty and no error should be returned
	assert.NotEmpty(t, resolvedIP)
	assert.NoError(t, err)

	// result has type string, but it should be a slice of strings
	// convert it to a slice of strings
	resolvedIPs := strings.Split(resolvedIP[1:len(resolvedIP)-1], " ")

	// Check if the resolved IP is a valid IP address
	parsedIP := net.ParseIP(resolvedIPs[0])
	assert.NotNil(t, parsedIP)
}

func TestGetHostByNameNXDomain(t *testing.T) {
	// GIVEN an invalid hostname
	tpl := `{{"invalid.invalid" | getHostByName}}`

	// WHEN getHostByName is executed
	resolvedIP, err := runRaw(tpl, nil)

	// THEN the resolved IP should be empty and an error should be returned
	assert.Empty(t, resolvedIP)
	assert.Error(t, err)
}
