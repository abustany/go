package charset

// #cgo pkg-config: icu-i18n
// #include <unicode/ucnv.h>
// #include <stdlib.h>
import "C"

import (
	"errors"
	"unsafe"
)

// makeError is declared in detect.go

// Converts data from one character encoding to another
func Convert(text []byte, srcenc string, dstenc string) ([]byte, error) {
	if text == nil {
		return nil, errors.New("nil input")
	}

	if len(text) == 0 {
		return make([]byte, 0), nil
	}

	var uerr C.UErrorCode = C.U_ZERO_ERROR

	c_srcenc := C.CString(srcenc)
	defer C.free(unsafe.Pointer(c_srcenc))

	c_dstenc := C.CString(dstenc)
	defer C.free(unsafe.Pointer(c_dstenc))

	// Preflight
	var sz C.int32_t = C.ucnv_convert(c_dstenc,
		c_srcenc,
		(*C.char)(unsafe.Pointer(nil)),
		0,
		(*C.char)(unsafe.Pointer(&text[0])),
		C.int32_t(len(text)),
		&uerr)

	if uerr != C.U_BUFFER_OVERFLOW_ERROR {
		return nil, errors.New("Cannot compute length of target buffer")
	}

	target := make([]byte, sz)

	uerr = C.U_ZERO_ERROR

	C.ucnv_convert(c_dstenc,
		c_srcenc,
		(*C.char)(unsafe.Pointer(&target[0])),
		C.int32_t(len(target)),
		(*C.char)(unsafe.Pointer(&text[0])),
		C.int32_t(len(text)),
		&uerr)

	if uerr != C.U_ZERO_ERROR && uerr != C.U_STRING_NOT_TERMINATED_WARNING {
		return nil, makeError(uerr)
	}

	return target, nil
}
