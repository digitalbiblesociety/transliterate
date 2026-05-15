// Package deva provides ISO 15919 transliteration of Devanagari script
// (U+0900..U+097F), used to write Hindi, Marathi, Sanskrit, Nepali,
// Bhojpuri, and other Indo-Aryan languages.
package deva

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x0900
	BlockEnd   rune = 0x097F
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Devanagari)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
