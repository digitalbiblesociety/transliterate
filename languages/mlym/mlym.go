// Package mlym provides ISO 15919 transliteration of the Malayalam
// script (U+0D00..U+0D7F), used for Malayalam.
package mlym

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x0D00
	BlockEnd   rune = 0x0D7F
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Malayalam)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
