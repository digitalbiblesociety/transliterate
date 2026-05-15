// Package newa provides ISO 15919-style romanization of the Newa
// script (U+11400..U+1147F) used for Nepal Bhasa / Newari (Kathmandu Valley, ~1M speakers). Codepoints derived from Aksharamukha (MIT).
package newa

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x11400
	BlockEnd   rune = 0x1147F
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Newa)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
