// Package hani provides Hanyu Pinyin transliteration of Chinese
// (Han / CJK U+4E00..U+9FFF) characters when read as Mandarin.
//
// Readings come from two vendored MIT-licensed sources:
//
//   - mozillazg/go-pinyin (internal/pinyin) — per-character primary
//     reading from the Unicode Unihan Database. Vendored in-tree; not a
//     go.mod dependency.
//   - mozillazg/phrase-pinyin-data (internal/pinyin/data) — ~47k
//     curated multi-character phrases + ~900 manual overrides. Used for
//     polyphone disambiguation: 中 reads zhōng in 中国 but zhòng in 击中.
//
// Lookup is greedy longest-match: at each position we try the longest
// phrase that fits (capped at the dictionary's maximum phrase length),
// then fall back to per-character lookup. Non-Han runes pass through.
//
// Known simplifications:
//   - No word segmentation outside the phrase dictionary. Characters
//     that don't appear in any matched phrase use their primary single-
//     character reading from Unihan even when context might suggest a
//     different one.
//   - Only the basic CJK block U+4E00..U+9FFF. Extensions (A, B, C, …)
//     are not covered (the embedded dict has no entries for them).
package hani

import (
	"strings"

	"github.com/digitalbiblesociety/transliterate/internal/pinyin"
)

const (
	BlockStart rune = 0x4E00
	BlockEnd   rune = 0x9FFF
)

// Transliterate returns tone-marked pinyin for s. Han syllables are
// space-separated; non-Han runes pass through unchanged.
func Transliterate(s string) string { return transliterate(s, false) }

// TransliterateAtonal returns ASCII pinyin (tone marks stripped, ü → u).
// Useful when downstream consumers can't handle combining diacritics.
func TransliterateAtonal(s string) string { return transliterate(s, true) }

// Contains reports whether s has any CJK-block rune.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}

func transliterate(s string, atonal bool) string {
	if s == "" {
		return s
	}
	rs := []rune(s)
	maxLen := pinyin.MaxPhraseRunes()
	var b strings.Builder
	b.Grow(len(s) * 3)
	prev := false

	for i := 0; i < len(rs); {
		r := rs[i]

		if r < BlockStart || r > BlockEnd {
			b.WriteRune(r)
			prev = false
			i++
			continue
		}

		// Longest-match phrase lookup. Cap the scan at the dictionary's
		// known maximum length so we don't pay for impossible windows.
		n := min(len(rs)-i, maxLen)
		matched := 0
		var syls []string
		for L := n; L >= 2; L-- {
			window := string(rs[i : i+L])
			if py, ok := pinyin.LookupPhrase(window); ok && len(py) == L {
				syls = py
				matched = L
				break
			}
		}

		if matched > 0 {
			for _, syl := range syls {
				if prev {
					b.WriteByte(' ')
				}
				if atonal {
					b.WriteString(stripTones(syl))
				} else {
					b.WriteString(syl)
				}
				prev = true
			}
			i += matched
			continue
		}

		// Single-character fallback. SinglePinyin returns the primary
		// reading (Heteronym off) as a one-element slice.
		single := pinyin.SinglePinyin(r, pinyin.Args{Style: pinyin.Tone})
		if len(single) > 0 && single[0] != "" {
			if prev {
				b.WriteByte(' ')
			}
			if atonal {
				b.WriteString(stripTones(single[0]))
			} else {
				b.WriteString(single[0])
			}
			prev = true
		} else {
			// In-block but no reading — drop (matches previous behavior).
			prev = false
		}
		i++
	}
	return b.String()
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
