package model

import (
	"testing"
)

func TestGidGen(t *testing.T) {
	cases := []struct {
		book_name, author_name string
		sign                   uint32
	}{
		//		{"布老虎传记文库·巨人百传丛书——科尔", "李尊玉", 3510365374},
		{"渎神曲", "涤生", 2700582757},
	}

	for _, c := range cases {
		gid := GetGid(c.book_name, c.author_name)
		if gid != c.sign {
			t.Errorf("except gid is %d, actual is %d", c.sign, gid)
		}

	}
}
