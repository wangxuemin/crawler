package model

import (
	"testing"
)

func TestChineseFilter(t *testing.T) {
	cases := []struct {
		out, in string
	}{
		{"你好", "你好"},
		{"你好abc123", "你好abc123"},
		{"你好", "你，好"},
		{"你好", "你   好"},
		{"你好abc123", "你好a b c 1 2 3"},
	}

	for _, c := range cases {
		tmp := ChineseU8_punctuation_filter(c.in)
		if tmp != c.out {
			t.Errorf("testing error %s ==> %s != %s", c.in, tmp, c.out)
		}
	}
}
