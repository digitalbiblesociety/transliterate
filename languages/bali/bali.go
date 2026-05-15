// Package bali provides ISO 15919 transliteration of the Balinese
// script (U+1B00..U+1B7F), used for Balinese, Kawi.
package bali

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x1B00
	BlockEnd   rune = 0x1B7F
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Balinese)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
