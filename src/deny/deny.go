package deny

import (
	"bytes"
)

var KindWord = [][]byte{[]byte("$"), []byte("{"), []byte("}")}

func InjectionPass(word []byte) bool {
	for _, v := range KindWord {
		if bytes.Contains(word, v) {
			return false
		}
	}
	return true
}
