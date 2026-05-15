// Package guru provides ISO 15919 transliteration of the Gurmukhi
// script (U+0A00..U+0A7F), used for Punjabi.
package guru

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x0A00
	BlockEnd   rune = 0x0A7F
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Gurmukhi)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
