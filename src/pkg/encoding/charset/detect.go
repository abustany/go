package charset

// #cgo pkg-config: icu-i18n
// #include <unicode/ucsdet.h>
import "C"

import (
	"errors"
	"unsafe"
)

func makeError(code C.UErrorCode) error {
	return errors.New(C.GoString(C.u_errorName(code)))
}

// Tries to detect the character encoding of text
func Detect(text []byte) (string, error) {
	if text == nil || len(text) == 0 {
		return "", errors.New("empty input")
	}

	var uerr C.UErrorCode = C.U_ZERO_ERROR
	var det *C.UCharsetDetector = C.ucsdet_open(&uerr)

	if uerr != C.U_ZERO_ERROR {
		return "", makeError(uerr)
	}

	defer C.ucsdet_close(det)

	C.ucsdet_setText(det, (*C.char)(unsafe.Pointer(&text[0])), C.int32_t(len(text)), &uerr)

	if uerr != C.U_ZERO_ERROR {
		return "", makeError(uerr)
	}

	var match *C.UCharsetMatch = C.ucsdet_detect(det, &uerr)

	if uerr != C.U_ZERO_ERROR {
		return "", makeError(uerr)
	}

	var c_encname *C.char = C.ucsdet_getName(match, &uerr)

	if uerr != C.U_ZERO_ERROR {
		return "", makeError(uerr)
	}

	return C.GoString(c_encname), nil
}
