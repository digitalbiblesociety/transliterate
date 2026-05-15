// Package orya provides ISO 15919 transliteration of the Oriya
// script (U+0B00..U+0B7F), used for Odia (Oriya), Santali (historically).
package orya

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x0B00
	BlockEnd   rune = 0x0B7F
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Oriya)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
