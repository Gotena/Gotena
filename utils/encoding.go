package utils

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func ReadUTF16String(data []byte) string {
	win16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	utf16bom := unicode.BOMOverride(win16le.NewDecoder())
	unicodeReader := transform.NewReader(bytes.NewReader(data), utf16bom)
	decoded, err := ioutil.ReadAll(unicodeReader)
	if err != nil {
		panic(err)
	}

	return string(decoded)
}
func WriteUTF16String(data string) []byte {
	win16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	utf16bom := unicode.BOMOverride(win16le.NewEncoder())
	unicodeWriter := transform.NewReader(bytes.NewReader([]byte(data)), utf16bom)
	encoded, err := ioutil.ReadAll(unicodeWriter)
	if err != nil {
		panic(err)
	}

	return encoded
}
