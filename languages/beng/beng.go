// Package beng provides ISO 15919 transliteration of the Bengali-Assamese
// script (U+0980..U+09FF), used to write Bengali, Assamese, Bishnupriya,
// Meitei (historically), and other Indo-Aryan languages.
package beng

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x0980
	BlockEnd   rune = 0x09FF
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.Bengali)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
