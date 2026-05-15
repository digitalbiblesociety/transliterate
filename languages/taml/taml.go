// Package taml provides ISO 15919 transliteration of the Tamil
// script (U+0B80..U+0BFF), used for Tamil.
package taml

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x0B80
	BlockEnd   rune = 0x0BFF
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Tamil)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
