package charset

import (
	"fmt"
	"testing"
)

var testdata_iso8859_1 = []byte{0x41, 0x75, 0x66, 0x74, 0x72, 0x61, 0x67, 0x73, 0x62, 0x65, 0x73, 0x74, 0xe4, 0x74, 0x69, 0x67, 0x75, 0x6e, 0x67, 0x2e, 0x70, 0x64, 0x66}

var testdata_utf8 = []byte{0x41, 0x75, 0x66, 0x74, 0x72, 0x61, 0x67, 0x73, 0x62, 0x65, 0x73, 0x74, 0xc3, 0xa4, 0x74, 0x69, 0x67, 0x75, 0x6e, 0x67, 0x2e, 0x70, 0x64, 0x66}

func TestDetect(t *testing.T) {
	testData := []struct {
		data     []byte
		encoding string
	}{
		{
			testdata_iso8859_1,
			"ISO-8859-1",
		},
		{
			testdata_utf8,
			"UTF-8",
		},
		{
			nil, // nil input is an error
			"",
		},
	}

	for _, test := range testData {
		enc, err := Detect(test.data)

		if test.encoding != "" && err != nil {
			t.Errorf("Error detecting the encoding of %s text", test.encoding)
			continue
		}

		if test.encoding == "" && err == nil {
			t.Errorf("Test did not fail as expected for input %v", test.data)
			continue
		}

		if test.encoding != enc {
			t.Errorf("Wrong encoding detection, got %s, expected %s", enc, test.encoding)
			continue
		}
	}
}

func compareData(expected []byte, val []byte) error {
	if len(expected) != len(val) {
		return fmt.Errorf("Unexpected length: expected %d, got %d", len(expected), len(val))
	}

	for i, x := range expected {
		if x != val[i] {
			return fmt.Errorf("Unexpected byte at index %d: expected %x, got %x", i, x, val[i])
		}
	}

	return nil
}

func TestConvert(t *testing.T) {
	testData := []struct {
		srcdata []byte
		srcenc  string
		dstdata []byte
		dstenc  string
		success bool
	}{
		{
			testdata_utf8,
			"UTF-8",
			testdata_iso8859_1,
			"ISO-8859-1",
			true,
		},
	}

	for _, test := range testData {
		data, err := Convert(test.srcdata, test.srcenc, test.dstenc)

		if test.success && err != nil {
			t.Errorf("Conversion from %s to %s unexpectedly failed: %s", test.srcenc, test.dstenc, err)
			continue
		}

		if !test.success && err == nil {
			t.Errorf("Expected error for conversion from %s to %s", test.srcenc, test.dstenc)
			continue
		}

		if err := compareData(test.dstdata, data); err != nil {
			t.Errorf("Unexpected result for conversion from %s to %s: %s", test.srcenc, test.dstenc, err)
			continue
		}
	}
}
