package sprig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFpBase(t *testing.T) {
	assert.NoError(t, runt(`{{ fpBase "C:\\foo\\bar" }}`, "bar"))
}

func TestFpDir(t *testing.T) {
	assert.NoError(t, runt(`{{ fpDir "C:\\foo\\bar\\baz" }}`, "C:\\foo\\bar"))
}

func TestFpIsAbs(t *testing.T) {
	assert.NoError(t, runt(`{{ fpIsAbs "C:\\foo" }}`, "true"))
	assert.NoError(t, runt(`{{ fpIsAbs "foo" }}`, "false"))
}

func TestFpClean(t *testing.T) {
	assert.NoError(t, runt(`{{ fpClean "C:\\foo\\..\\foo\\..\\bar" }}`, "C:\\bar"))
}

func TestFpExt(t *testing.T) {
	assert.NoError(t, runt(`{{ fpExt "C:\\foo\\bar\\baz.txt" }}`, ".txt"))
}
