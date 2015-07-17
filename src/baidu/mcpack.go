package baidu

// #include <stdlib.h>
// #include <stddef.h>
// #include <sys/types.h>
// #include <mc_pack_def.h>
// #include <mc_pack_c.h>
// #include <nshead.h>
// int go_mc_pack_put_int32(mc_pack_t *ppack, char *key, int value) {
//     return mc_pack_put_int32(ppack, key, value);
// }
// int go_mc_pack_get_int32(mc_pack_t *ppack, char *key, int *value){
//     return mc_pack_get_int32(ppack, key, value);
// }
// int go_mc_pack_put_int64(mc_pack_t *ppack, char *key, long long value) {
//     return mc_pack_put_int64(ppack, key, value);
// }
// int go_mc_pack_get_int64(mc_pack_t *ppack, char *key, long long *value){
//     return mc_pack_get_int64(ppack, key, value);
// }
// const char* go_mc_pack_get_raw(mc_pack_t *ppack, char *key, int *len) {
//     return mc_pack_get_raw(ppack, key, len);
// }
// #cgo pkg-config: baidu
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

const (
	PACK_DEFAULT_SIZE = 102400
	PACK_TMP_BUFSIZE  = 102400
)

type McpackError struct {
	s string
}

func (error *McpackError) Error() string {
	return error.s
}

func Marshal(i interface{}) ([]byte, *McpackError) {
	buf := C.malloc(PACK_DEFAULT_SIZE)
	defer C.free(buf)
	tmpbuf := C.malloc(PACK_TMP_BUFSIZE)
	defer C.free(tmpbuf)

	pack := C.mc_pack_open_w(C.int(2), (*C.char)(buf), PACK_DEFAULT_SIZE,
		(*C.char)(tmpbuf), PACK_TMP_BUFSIZE)
	defer C.mc_pack_close(pack)

	if err := marshal(pack, reflect.ValueOf(i)); err != nil {
		return nil, err
	}
	packet := C.GoBytes(buf, C.mc_pack_get_size(pack))

	return packet, nil
}

func Unmarshal(bytes []byte, i interface{}) *McpackError {
	buf := unsafe.Pointer(&bytes[0])
	tmpsize := len(bytes) * 3
	tmpbuf := C.malloc(C.size_t(tmpsize))
	defer C.free(tmpbuf)
	pack := C.mc_pack_open_r((*C.char)(buf), C.int(len(bytes)),
		(*C.char)(tmpbuf), C.int(tmpsize))
	defer C.mc_pack_close(pack)

	if err := unmarshal(pack, i); err != nil {
		return err
	}
	return nil
}

func unpackValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface && !v.IsNil() {
		v = v.Elem()
	}
	return v
}

func pack_kv(pack *C.struct___0, key string, v reflect.Value) *McpackError {
	switch v.Kind() {
	case reflect.Int32, reflect.Int:
		C.go_mc_pack_put_int32(pack, C.CString(key), C.int(v.Int()))
	case reflect.Int64:
		C.go_mc_pack_put_int64(pack, C.CString(key), C.longlong(v.Int()))
	case reflect.String:
		C.mc_pack_put_str(pack, C.CString(key), C.CString(v.String()))
	case reflect.Struct:
		subpack := C.mc_pack_put_object(pack, C.CString(key))
		if err := marshal(subpack, v); err != nil {
			return err
		}
		C.mc_pack_finish(subpack)
	default:
		return &McpackError{"unknown type: " + v.Kind().String()}
	}
	return nil
}

func marshal(pack *C.struct___0, v reflect.Value) *McpackError {
	vt := v.Type()
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			vtf := vt.Field(i)
			realName := vtf.Tag.Get("mcpack")
			if realName == "" {
				realName = vtf.Name
			}

			if err := pack_kv(pack, realName, unpackValue(v.Field(i))); err != nil {
				return err
			}
		}
	default:
		return &McpackError{"unknown type: " + v.Kind().String()}
	}
	return nil
}

func unpack_kv(pack *C.struct___0, v reflect.Value) *McpackError {
	vt := v.Type()

	for i := 0; i < v.NumField(); i++ {
		vtf := vt.Field(i)
		vf := v.Field(i)

		realName := vtf.Tag.Get("mcpack")
		if realName == "" {
			realName = vtf.Name
		}

		switch vf.Kind() {
		case reflect.Struct:
			subpack := C.mc_pack_get_object(pack, C.CString(realName))
			if err := unpack_kv(subpack, vf); err != nil {
				return err
			}
		case reflect.String:
			fmt.Println(realName)
			value_c := C.mc_pack_get_str(pack, C.CString(realName))
			vf.SetString(C.GoString(value_c))
		case reflect.Int, reflect.Int32:
			value_i := C.int(0)
			C.go_mc_pack_get_int32(pack, C.CString(realName), &value_i)
			vf.SetInt(int64(value_i))
		case reflect.Int64:
			value_i := C.longlong(0)
			C.go_mc_pack_get_int64(pack, C.CString(realName), &value_i)
			vf.SetInt(int64(value_i))
		case reflect.Slice, reflect.Array:
			if vf.Type().Elem().Kind() != reflect.Uint8 {
				return &McpackError{"unspport slice type"}
			}
			value_len := C.int(0)
			value_r := C.go_mc_pack_get_raw(pack, C.CString(realName), &value_len)
			vf.SetBytes(C.GoBytes(unsafe.Pointer(value_r), value_len))

		default:
			return &McpackError{"unsupport type: " + v.Kind().String()}
		}
	}

	return nil
}

func unmarshal(pack *C.struct___0, i interface{}) *McpackError {
	rv := reflect.ValueOf(i)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &McpackError{"invalid unmarshal type " + reflect.TypeOf(rv).String()}
	}

	return unpack_kv(pack, rv.Elem())

}
