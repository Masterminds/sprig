package sprig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFpBase(t *testing.T) {
	assert.NoError(t, runt(`{{ fpBase "foo/bar" }}`, "bar"))
}

func TestFpDir(t *testing.T) {
	assert.NoError(t, runt(`{{ fpDir "foo/bar/baz" }}`, "foo/bar"))
}

func TestFpIsAbs(t *testing.T) {
	assert.NoError(t, runt(`{{ fpIsAbs "/foo" }}`, "true"))
	assert.NoError(t, runt(`{{ fpIsAbs "foo" }}`, "false"))
}

func TestFpClean(t *testing.T) {
	assert.NoError(t, runt(`{{ fpClean "/foo/../foo/../bar" }}`, "/bar"))
}

func TestFpExt(t *testing.T) {
	assert.NoError(t, runt(`{{ fpExt "/foo/bar/baz.txt" }}`, ".txt"))
}
