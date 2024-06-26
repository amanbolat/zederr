// Code generated by zederr generator. DO NOT EDIT.
package zederr

import (
	_ "embed"
)

//go:embed locale.en.toml
var enErrorCodeMessages []byte

//go:embed locale.zh.toml
var zhErrorCodeMessages []byte

var errorMessagesMap = map[string][]byte{
	"en": enErrorCodeMessages,
	"zh": zhErrorCodeMessages,
}

// ErrorMessages returns a map of error messages for.
// It is safe to modify the returned map and its values.
func ErrorMessages() map[string][]byte {
	res := map[string][]byte{}
	for lang, msgs := range errorMessagesMap {
		b := make([]byte, len(msgs))
		copy(b, msgs)
		res[lang] = b
	}

	return res
}
