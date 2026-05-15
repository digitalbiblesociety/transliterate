// Package sund provides ISO 15919 transliteration of the Sundanese
// script (U+1B80..U+1BBF), used for Sundanese.
package sund

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x1B80
	BlockEnd   rune = 0x1BBF
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Sundanese)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
