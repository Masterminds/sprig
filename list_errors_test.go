package sprig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustListFunctionsNil(t *testing.T) {
	tests := []struct {
		name string
		call func() error
	}{
		{"push", func() error { _, err := mustPush(nil, 1); return err }},
		{"prepend", func() error { _, err := mustPrepend(nil, 1); return err }},
		{"chunk", func() error { _, err := mustChunk(1, nil); return err }},
		{"last", func() error { _, err := mustLast(nil); return err }},
		{"first", func() error { _, err := mustFirst(nil); return err }},
		{"rest", func() error { _, err := mustRest(nil); return err }},
		{"initial", func() error { _, err := mustInitial(nil); return err }},
		{"reverse", func() error { _, err := mustReverse(nil); return err }},
		{"compact", func() error { _, err := mustCompact(nil); return err }},
		{"uniq", func() error { _, err := mustUniq(nil); return err }},
		{"without", func() error { _, err := mustWithout(nil, 1); return err }},
		{"slice", func() error { _, err := mustSlice(nil); return err }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.NotPanics(t, func() { assert.Error(t, test.call()) })
		})
	}
}

func TestMustChunkInvalidSize(t *testing.T) {
	for _, size := range []int{0, -1, -10} {
		for _, list := range [][]int{nil, {}, {1}, {1, 2, 3}} {
			assert.NotPanics(t, func() {
				_, err := mustChunk(size, list)
				assert.Error(t, err, "size %d, list %v", size, list)
			})
		}
	}
}

func TestMustSliceInvalidBounds(t *testing.T) {
	tests := []struct {
		name    string
		list    interface{}
		indices []interface{}
	}{
		{"negative start", []int{1, 2, 3}, []interface{}{-1}},
		{"negative end", []int{1, 2, 3}, []interface{}{0, -1}},
		{"reversed bounds", []int{1, 2, 3}, []interface{}{2, 1}},
		{"start beyond capacity", []int{1, 2, 3}, []interface{}{4}},
		{"end beyond capacity", []int{1, 2, 3}, []interface{}{0, 4}},
		{"empty slice", []int{}, []interface{}{1}},
		{"nil slice", []int(nil), []interface{}{0, 1}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				_, err := mustSlice(test.list, test.indices...)
				assert.Error(t, err)
			})
		})
	}
}

func TestMustSliceArray(t *testing.T) {
	assert.NotPanics(t, func() {
		got, err := mustSlice([3]int{1, 2, 3}, 1, 3)
		assert.NoError(t, err)
		assert.Equal(t, []int{2, 3}, got)
	})
}

func TestMustSliceCapacity(t *testing.T) {
	list := []int{1, 2, 3}
	got, err := mustSlice(list[:1], 1, 3)
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 3}, got)
}
