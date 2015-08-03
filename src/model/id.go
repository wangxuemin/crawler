// id generators, such as gid/book_id/chapter_id/page_id etc..
package model

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"
)

func Sign(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}

func GetBookid(short_id uint64) uint64 {
	return (1 << 32) + short_id
}

func GetGid(novel_name, author_name string) uint32 {
	novel_name_gbk := Encoding_U82G(novel_name)
	author_name_gbk := Encoding_U82G(author_name)
	return Sign(Chinese_punctuation_filter(novel_name_gbk) +
		"@" +
		Chinese_punctuation_filter(author_name_gbk))
}
