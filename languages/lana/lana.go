// Package lana provides ISO 15919-style romanization of the TaiTham
// script (U+1A20..U+1AAF) used for Northern Thai (Kam Mueang, ~6M), Tai Lue, Khün, Lao Tham. Codepoints derived from Aksharamukha (MIT).
package lana

import "github.com/digitalbiblesociety/transliterate/internal/brahmic"

const (
	BlockStart rune = 0x1A20
	BlockEnd   rune = 0x1AAF
)

func Transliterate(s string) string {
	return brahmic.Transliterate(s, brahmic.TaiTham)
}

func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
