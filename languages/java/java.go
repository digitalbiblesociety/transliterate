// Package java provides ISO 15919 transliteration of the Javanese
// script (U+A980..U+A9DF), used for Javanese, Kawi (Old Javanese).
package java

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0xA980
	BlockEnd   rune = 0xA9DF
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Javanese)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
