// Package hani provides Hanyu Pinyin transliteration of Chinese
// (Han / CJK U+4E00..U+9FFF) characters when read as Mandarin.
//
// The per-character reading table is derived from the kMandarin field of
// the Unicode Unihan Database (Unicode Terms of Use). Transliterate
// emits tone-marked pinyin (e.g. "nǐ hǎo") by default; TransliterateAtonal
// strips the diacritics back to ASCII letters ("ni hao").
//
// Known simplifications:
//   - One reading per character. Unihan supplies a single primary
//     Mandarin reading per codepoint; rare characters with multiple
//     pronunciations get only the most common one.
//   - No word segmentation: syllables are emitted character by character
//     with spaces between them.
//   - Only the basic CJK block U+4E00..U+9FFF. Extensions (A, B, C, …)
//     are not covered.
package hani

import "strings"

const (
	BlockStart rune = 0x4E00
	BlockEnd   rune = 0x9FFF
)

// Transliterate returns tone-marked pinyin for s. Han characters emit a
// trailing space between syllables; non-Han runes pass through unchanged.
func Transliterate(s string) string { return transliterate(s, false) }

// TransliterateAtonal returns ASCII pinyin (tone marks stripped, ü → u).
// Useful when downstream consumers can't handle combining diacritics.
func TransliterateAtonal(s string) string { return transliterate(s, true) }

func transliterate(s string, atonal bool) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) * 3)
	prev := false
	for _, r := range s {
		if v, ok := generatedTonalMap[r]; ok {
			if prev {
				b.WriteByte(' ')
			}
			if atonal {
				b.WriteString(stripTones(v))
			} else {
				b.WriteString(v)
			}
			prev = true
			continue
		}
		if r >= BlockStart && r <= BlockEnd {
			prev = false
			continue
		}
		b.WriteRune(r)
		prev = false
	}
	return b.String()
}

// Contains reports whether s has any CJK-block rune.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}

// toneStrip maps every tone-marked pinyin vowel back to its plain ASCII
// base. ü (and its toned forms) all collapse to "u" — most consumers
// either don't render ü or treat it interchangeably with u.
var toneStrip = map[rune]rune{
	'ā': 'a', 'á': 'a', 'ǎ': 'a', 'à': 'a',
	'ē': 'e', 'é': 'e', 'ě': 'e', 'è': 'e',
	'ī': 'i', 'í': 'i', 'ǐ': 'i', 'ì': 'i',
	'ō': 'o', 'ó': 'o', 'ǒ': 'o', 'ò': 'o',
	'ū': 'u', 'ú': 'u', 'ǔ': 'u', 'ù': 'u',
	'ǖ': 'u', 'ǘ': 'u', 'ǚ': 'u', 'ǜ': 'u',
	'ü': 'u',
}

func stripTones(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if base, ok := toneStrip[r]; ok {
			b.WriteRune(base)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
