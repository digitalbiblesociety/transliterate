// Package georgian provides BGN/PCGN transliteration of the Georgian
// scripts (Mkhedruli, Asomtavruli) to the Latin alphabet.
//
// Modern Georgian uses Mkhedruli (U+10D0-U+10FA, 33 letters). Older
// liturgical texts may use Asomtavruli (U+10A0-U+10C5); we map those
// to the same Latin forms.
package geor

import "strings"

const (
	BlockStart rune = 0x10A0
	BlockEnd   rune = 0x10FF
)

// Transliterate returns the BGN/PCGN romanization of s. Non-Georgian
// runes pass through unchanged.
func Transliterate(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) * 2)
	for _, r := range s {
		if v, ok := letterMap[r]; ok {
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

// Contains reports whether s has any Georgian-block rune.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}

// letterMap holds Mkhedruli (lowercase) and its Asomtavruli twins.
// Georgian has no case distinction in modern usage, so both scripts map
// to lowercase Latin.
var letterMap = map[rune]string{
	// Mkhedruli (modern, U+10D0..U+10FA).
	0x10D0: "a",  0x10D1: "b",  0x10D2: "g",  0x10D3: "d",  0x10D4: "e",
	0x10D5: "v",  0x10D6: "z",  0x10D7: "t",  0x10D8: "i",  0x10D9: "k'",
	0x10DA: "l",  0x10DB: "m",  0x10DC: "n",  0x10DD: "o",  0x10DE: "p'",
	0x10DF: "zh", 0x10E0: "r",  0x10E1: "s",  0x10E2: "t'", 0x10E3: "u",
	0x10E4: "p",  0x10E5: "k",  0x10E6: "gh", 0x10E7: "q'", 0x10E8: "sh",
	0x10E9: "ch", 0x10EA: "ts", 0x10EB: "dz", 0x10EC: "ts'", 0x10ED: "ch'",
	0x10EE: "kh", 0x10EF: "j",  0x10F0: "h",
	// Archaic Mkhedruli.
	0x10F1: "ē", 0x10F2: "y", 0x10F3: "w", 0x10F4: "ḫ", 0x10F5: "ǰ",
	0x10F6: "f", 0x10F7: "ʿ", 0x10F8: "ʾ", 0x10F9: "ǎ", 0x10FA: "ŭ",
	// Asomtavruli (capital/old, U+10A0..U+10C5) — same Latin forms.
	0x10A0: "a",  0x10A1: "b",  0x10A2: "g",  0x10A3: "d",  0x10A4: "e",
	0x10A5: "v",  0x10A6: "z",  0x10A7: "t",  0x10A8: "i",  0x10A9: "k'",
	0x10AA: "l",  0x10AB: "m",  0x10AC: "n",  0x10AD: "o",  0x10AE: "p'",
	0x10AF: "zh", 0x10B0: "r",  0x10B1: "s",  0x10B2: "t'", 0x10B3: "u",
	0x10B4: "p",  0x10B5: "k",  0x10B6: "gh", 0x10B7: "q'", 0x10B8: "sh",
	0x10B9: "ch", 0x10BA: "ts", 0x10BB: "dz", 0x10BC: "ts'", 0x10BD: "ch'",
	0x10BE: "kh", 0x10BF: "j",  0x10C0: "h",
}
