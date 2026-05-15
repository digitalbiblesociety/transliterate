// Package telu provides ISO 15919 transliteration of the Telugu
// script (U+0C00..U+0C7F), used for Telugu, Lambadi.
package telu

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x0C00
	BlockEnd   rune = 0x0C7F
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Telugu)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
