// Package modi provides ISO 15919-style romanization of the Modi
// script (U+11600..U+1165F) used for Marathi (historical 17th-20th c. correspondence and administration). Codepoints derived from Aksharamukha (MIT).
package modi

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x11600
	BlockEnd   rune = 0x1165F
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Modi)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
