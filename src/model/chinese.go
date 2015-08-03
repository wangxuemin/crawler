package model

import (
	"fmt"
	"github.com/qiniu/iconv"
	"strings"
)

func Encoding_U82G(src string) string {
	cd, err := iconv.Open("GBK", "UTF-8//IGNORE")
	if err != nil {
		fmt.Println("open iconv error")
		return ""
	}
	defer cd.Close()

	return cd.ConvString(src)
}

func Encoding_G2U8(src string) string {
	cd, err := iconv.Open("UTF-8", "GBK")
	if err != nil {
		fmt.Println("open iconv error")
		return ""
	}
	defer cd.Close()

	return cd.ConvString(src)
}

func ChineseU8_punctuation_filter(s string) string {
	return Encoding_G2U8(Chinese_punctuation_filter(Encoding_U82G(s)))
}

func Encoding_U82U(src string) string {
	cd, err := iconv.Open("UTF8", "UNICODE")
	if err != nil {
		fmt.Println("open iconv error")
		return ""
	}
	defer cd.Close()

	return cd.ConvString(src)
}

//compatiable to pirate novel
//can accept only gbk encoding
func Chinese_punctuation_filter(s string) string {
	bytes := []byte(s)
	output := make([]byte, len(bytes))

	oidx := 0
	for idx := 0; idx < len(bytes); {
		ch := bytes[idx]
		if (ch >= '0' && ch <= '9') ||
			(ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') {
			output[oidx] = ch
			oidx++
			idx++
			continue
		}

		//check special chars
		if ch <= 0xff && ch >= 0x80 && idx+1 < len(bytes) {
			ch2 := bytes[idx+1]
			if (ch == 0xA2 && (ch2 >= 0xA1 && ch2 <= 0xFE)) ||
				((ch >= 0xB0 && ch <= 0xF7) && (ch2 >= 0xA1 && ch2 <= 0xFE)) ||
				((ch >= 0x81 && ch <= 0xA0) && (ch2 >= 0x40 && ch2 <= 0xFE)) ||
				((ch >= 0xAA && ch <= 0xFE) && (ch2 >= 0x40 && ch2 <= 0xA0)) ||
				((ch == 0xA3) && ((ch2 >= 0xB0 && ch2 <= 0xB9) ||
					(ch2 >= 0xC1 && ch2 <= 0xDA) ||
					(ch2 >= 0xE1 && ch2 <= 0xFA))) {
				//set whole gbk char here
				output[oidx] = ch
				oidx++
				output[oidx] = ch2
				oidx++
				idx += 2
				continue
			}
		}

		idx += 1 + int((ch&0x80)>>7) //TODO what's this
	}

	return strings.ToLower(string(output[:oidx]))
}

func Chinese_wc(s string) int {
	cnt := 0
	for _, c := range s {
		fmt.Printf("%x\n", c)
		if (c >= 0x4E00 && c <= 0xA000) ||
			(c >= 0x3400 && c <= 0x4DC0) ||
			(c >= 0x200000 && c <= 0x2A6E0) ||
			(c >= 0x2A700 && c <= 0x2B740) ||
			(c >= 0x2B740 && c <= 0x2B820) ||
			(c >= 0xF900 && c <= 0xFB00) ||
			(c >= 0x9FA6 && c <= 0x9FCC) ||
			(c >= 0x2F800 && c <= 0x2FA20) {
			cnt += 1
		}
	}
	return cnt
}
