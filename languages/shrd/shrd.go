// Package shrd provides ISO 15919-style romanization of the Sharada
// script (U+11180..U+111DF), used historically for Kashmiri Śaiva and
// Sanskrit texts (8th-19th century).
//
// Codepoint mappings derived from Aksharamukha (MIT).
package shrd

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x11180
	BlockEnd   rune = 0x111DF
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Sharada)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
