package strutil

import (
	"path"
	"unicode"

	"k8s.io/client-go/util/homedir"
)

func ContainsBinary(str string) bool {
	//buf := []byte(str)
	//size := 0
	//for start := 0; start < len(buf); start += size {
	//	var r rune
	//	if r, size = utf8.DecodeRune(buf[start:]); r == utf8.RuneError {
	//		return true
	//	}
	//}
	rs := []rune(str)
	for _, r := range rs {
		if r == unicode.ReplacementChar {
			return true
		}
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			return true
		}
	}
	return false
}

func IsSameChar(str string) bool {
	if len(str) <= 0 {
		return false
	}

	rs := []rune(str)
	first := rs[0]
	for _, r := range rs {
		if r != first {
			return false
		}
	}

	return true
}

func RootPath() string {
	return path.Join(homedir.HomeDir(), ".KTT")
}
