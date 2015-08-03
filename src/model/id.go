// id generators, such as gid/book_id/chapter_id/page_id etc..
package model

import (
	"hash/crc"
	"strings"
)

func Sign(s string) int32 {
	return crc.Checksum([]byte(string), nil)
}

func GetBookid(short_id uint64) uint64 {
	return (1 << 32) + short_id
}

func GetGid(novel_name, author_name string) uint32 {
	return Sign(Encoding_G2U8(novel_name) + "@" + Encoding_U82G(author_name))
}
