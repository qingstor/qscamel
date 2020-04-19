package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type joinCase struct {
	in  []string
	out string
}

func TestJoin(t *testing.T) {
	joinTestCases := []joinCase{
		{[]string{"/"}, ""},

		{[]string{"a"}, "a"},
		{[]string{"/", "/a"}, "/a"},

		{[]string{"a", "b"}, "a/b"},
		{[]string{"a/", "b"}, "a/b"},
		{[]string{"a/", "b/"}, "a/b"},
		{[]string{"/a", "/b"}, "a//b"},

		{[]string{"a", "", "b"}, "a/b"},
		{[]string{"//a/b", "c/"}, "/a/b/c"},
	}

	for _, v := range joinTestCases {
		assert.Equal(t, v.out, Join(v.in...))
	}
}

type relativeCase struct {
	full   string
	prefix string
	out    string
}

func TestRelative(t *testing.T) {
	relativeTestCases := []relativeCase{
		{"/a/b/c", "/a/b", "c"},
		{"/a/b/c", "/a", "b/c"},
		{"/a/b/c", "a/", "b/c"},
		{"/a/b/c", "a/b", "c"},
		{"/a/b/c", "a/b/", "c"},
		{"/a/b/c", "b/", "a/b/c"},
	}

	for _, v := range relativeTestCases {
		assert.Equal(t, v.out, Relative(v.full, v.prefix))
	}
}
