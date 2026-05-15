// Package tirh provides ISO 15919-style romanization of the Tirhuta
// script (U+11480..U+114DF) used for Maithili (Bihar, Nepal — ~50M speakers). Codepoints derived from Aksharamukha (MIT).
package tirh

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x11480
	BlockEnd   rune = 0x114DF
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Tirhuta)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
