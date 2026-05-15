// Package bugi provides ISO 15919 transliteration of the Buginese
// script (U+1A00..U+1A1F), used for Bugis, Makassar.
package bugi

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x1A00
	BlockEnd   rune = 0x1A1F
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Buginese)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
