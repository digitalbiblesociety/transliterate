// Package knda provides ISO 15919 transliteration of the Kannada
// script (U+0C80..U+0CFF), used for Kannada, Tulu, Konkani.
package knda

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x0C80
	BlockEnd   rune = 0x0CFF
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Kannada)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
