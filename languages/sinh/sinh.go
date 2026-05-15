// Package sinh provides ISO 15919 transliteration of the Sinhala
// script (U+0D80..U+0DFF), used for Sinhala.
package sinh

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x0D80
	BlockEnd   rune = 0x0DFF
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Sinhala)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
