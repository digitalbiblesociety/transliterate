// Package gujr provides ISO 15919 transliteration of the Gujarati
// script (U+0A80..U+0AFF), used for Gujarati.
package gujr

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x0A80
	BlockEnd   rune = 0x0AFF
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Gujarati)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
