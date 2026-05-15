// Package brah provides ISO 15919-style romanization of the Brahmi
// script (U+11000..U+1107F), the ancestor of every Brahmic-family
// writing system. Used for the Aśokan edicts and early Buddhist/Jain
// inscriptions; relevant for biblical scholarship of comparative
// epigraphy in the Indian subcontinent.
//
// Codepoint mappings derived from Aksharamukha (MIT).
package brah

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x11000
	BlockEnd   rune = 0x1107F
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Brahmi)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
