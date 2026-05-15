// Package cansyl provides transliteration of Canadian Aboriginal
// Syllabics (U+1400-U+167F) to the Latin alphabet. The block is used
// for Cree, Inuktitut, Ojibwe, Naskapi, Blackfoot, Chipewyan, and
// related First Nations languages.
//
// The mapping table is generated mechanically from Unicode character
// names (see table.go) — each glyph is named after the syllable it
// represents (e.g., "CANADIAN SYLLABICS WE" → "we"). Dialect markers
// like "WEST-CREE" or "NASKAPI" are stripped, so glyphs for cognate
// syllables across dialects collapse to the same Latin form.
package cans

import "strings"

const (
	BlockStart rune = 0x1400
	BlockEnd   rune = 0x167F
)

// Transliterate returns the romanization of s. Non-syllabics runes
// pass through unchanged.
func Transliterate(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) * 2)
	for _, r := range s {
		if v, ok := generatedMap[r]; ok {
			b.WriteString(v)
			continue
		}
		if r >= BlockStart && r <= BlockEnd {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// Contains reports whether s has any Canadian Syllabics rune.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
