// Package batk provides ISO 15919 transliteration of the Batak
// script (U+1BC0..U+1BFF), used for Toba, Karo, Mandailing, Pakpak, Simalungun Batak.
package batk

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x1BC0
	BlockEnd   rune = 0x1BFF
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Batak)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
