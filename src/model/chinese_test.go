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

func TestWordsum(t *testing.T) {
	cases := []struct {
		word  string
		count int
	}{
		{"你好", 2},
		{"你,好,", 2},
		{"", 0},
		{"，，你好啊", 3},
	}

	for _, c := range cases {
		cnt := Chinese_wc(c.word)
		if cnt != c.count {
			t.Errorf("except wc %d, actual %d", c.count, cnt)
		}
	}
}
