// Package yueh provides Jyutping transliteration of Han characters
// (U+4E00..U+9FFF) when read as Cantonese (Yue Chinese).
//
// The per-character reading table is derived from the kCantonese field of
// the Unicode Unihan Database (Unicode Terms of Use). Each emitted
// syllable follows the Jyutping scheme: ASCII letters with a tone digit
// (1-6) appended; the neutral / mid-level tone is "1" through "6" with
// no diacritics.
//
// Known simplifications:
//   - One reading per character. Unihan supplies a single primary
//     Cantonese reading per codepoint; rare characters with multiple
//     pronunciations get only the most common one.
//   - No word segmentation: syllables are emitted character by character
//     with spaces between them.
//   - Only the basic CJK block U+4E00..U+9FFF. Extensions (A, B, C, …)
//     are not covered.
//   - Han characters lacking a Cantonese reading in Unihan are dropped
//     (not all Han codepoints have Cantonese assignments).
package yueh

import "strings"

const (
	BlockStart rune = 0x4E00
	BlockEnd   rune = 0x9FFF
)

// Transliterate returns the Jyutping romanization of s with tone digits.
// Han characters are emitted with a trailing space so syllables remain
// readable. Non-Han runes pass through unchanged.
func Transliterate(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) * 3)
	prev := false
	for _, r := range s {
		if v, ok := generatedJyutpingMap[r]; ok {
			if prev {
				b.WriteByte(' ')
			}
			b.WriteString(v)
			prev = true
			continue
		}
		if r >= BlockStart && r <= BlockEnd {
			// In-block but unmapped (no Cantonese reading in Unihan) — drop.
			prev = false
			continue
		}
		b.WriteRune(r)
		prev = false
	}
	return b.String()
}

// TransliterateAtonal is like Transliterate but strips the trailing tone
// digit from each Jyutping syllable, producing ASCII letters only.
func TransliterateAtonal(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) * 3)
	prev := false
	for _, r := range s {
		if v, ok := generatedJyutpingMap[r]; ok {
			if prev {
				b.WriteByte(' ')
			}
			b.WriteString(stripToneDigit(v))
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

func stripToneDigit(jyutping string) string {
	n := len(jyutping)
	if n == 0 {
		return jyutping
	}
	last := jyutping[n-1]
	if last >= '1' && last <= '6' {
		return jyutping[:n-1]
	}
	return jyutping
}
