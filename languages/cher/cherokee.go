// Package cherokee provides transliteration of the Cherokee syllabary
// (U+13A0-U+13FF for uppercase, U+AB70-U+ABBF for lowercase) to the
// Latin alphabet using Sequoyah's syllable values. The mapping is
// generated mechanically from Unicode character names — each Cherokee
// glyph is named after the syllable it represents (e.g., "CHEROKEE
// LETTER GA" → "ga").
package cher

import "strings"

const (
	BlockStart      rune = 0x13A0
	BlockEnd        rune = 0x13FF
	LowerBlockStart rune = 0xAB70
	LowerBlockEnd   rune = 0xABBF
)

// Transliterate returns the romanization of s. Non-Cherokee runes pass
// through unchanged.
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
		if inBlock(r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// Contains reports whether s has any Cherokee-block rune.
func Contains(s string) bool {
	for _, r := range s {
		if inBlock(r) {
			return true
		}
	}
	return false
}

func inBlock(r rune) bool {
	return (r >= BlockStart && r <= BlockEnd) ||
		(r >= LowerBlockStart && r <= LowerBlockEnd)
}
