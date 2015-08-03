package baidu

// #include <nshead.h>
// #include <string.h>
// void fill_head(void *vhead, int len, int logid) {
//     nshead_t *head = (nshead_t *)vhead;
//     head->body_len = len;
//     head->magic_num = NSHEAD_MAGICNUM;
//     strcpy(head->provider, "Go-AHead");
//     head->version = 2;
// }
//
// int valid_pack(void *buf) {
//     nshead_t *head = (nshead_t *)buf;
//     if (head->magic_num != NSHEAD_MAGICNUM) {return -1;}
//     return head->body_len;
// }
// const int NSHEAD_SIZE = sizeof(nshead_t);
// #cgo pkg-config: baidu
import "C"

import (
	"bufio"
	"errors"
	"io"
	"os"
	"unsafe"
)

var (
	NSHEAD_SIZE  = int(C.NSHEAD_SIZE)
	ErrIO        = errors.New("io error")
	ErrMalFormed = errors.New("malformed head")
)

func NsheadPack(body []byte, logid int) []byte {
	bytes := make([]byte, NSHEAD_SIZE+len(body))
	C.fill_head(unsafe.Pointer(&bytes[0]), C.int(len(body)), C.int(logid))
	copy(bytes[NSHEAD_SIZE:], body)
	return bytes
}

func NsheadRead(r io.Reader) ([]byte, []byte, error) {
	head_reader := bufio.NewReaderSize(r, NSHEAD_SIZE)
	head, herr := head_reader.Peek(NSHEAD_SIZE)
	if herr != nil {
		return nil, nil, ErrIO
	}
	ret := int(C.valid_pack(unsafe.Pointer(&head[0])))
	if ret == -1 {
		return nil, nil, ErrMalFormed
	}

	body_reader := bufio.NewReaderSize(r, ret)
	body, berr := body_reader.Peek(ret)
	if berr != nil {
		return nil, nil, ErrIO
	}

	f, _ := os.Create("cspub_response")
	defer f.Close()
	f.Write(head)
	f.Write(body)

	return head, body, nil
}
