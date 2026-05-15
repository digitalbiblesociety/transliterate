// Package burmese provides BGN/PCGN transliteration of the Myanmar
// (Burmese) script (U+1000-U+109F) to the Latin alphabet.
//
// Burmese is an Abugida: each consonant carries an inherent /a/ unless
// followed by a vowel sign (matra) or asat (virama). Medial consonants
// (ya, ra, wa, ha) attach to the head consonant before the vowel.
//
// Known simplifications:
//   - Tone marks (U+1037, U+1038) are dropped; tone is not represented.
//   - "Killed" consonants (asat-terminated) are emitted bare without
//     marking syllable closure beyond suppressing the inherent /a/.
package mymr

import "strings"

const (
	BlockStart rune = 0x1000
	BlockEnd   rune = 0x109F
	Asat       rune = 0x103A // ◌် virama
)

// Transliterate returns the BGN/PCGN romanization of s. Non-Burmese
// runes pass through unchanged.
func Transliterate(s string) string {
	if s == "" {
		return s
	}
	rs := []rune(s)
	var b strings.Builder
	b.Grow(len(s) * 2)
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		if base, ok := consonants[r]; ok {
			b.WriteString(base)
			// Consume medial consonants and decorations.
			vowelSeen := false
			killedByAsat := false
			for i+1 < len(rs) {
				next := rs[i+1]
				if m, ok := medials[next]; ok {
					b.WriteString(m)
					i++
					continue
				}
				if v, ok := vowelSigns[next]; ok {
					b.WriteString(v)
					vowelSeen = true
					i++
					continue
				}
				if next == Asat {
					killedByAsat = true
					i++
					continue
				}
				if v, ok := tones[next]; ok {
					b.WriteString(v)
					i++
					continue
				}
				break
			}
			if !vowelSeen && !killedByAsat {
				b.WriteString("a")
			}
			continue
		}
		if v, ok := independentVowels[r]; ok {
			b.WriteString(v)
			continue
		}
		if r >= 0x1040 && r <= 0x1049 { // digits
			b.WriteRune('0' + (r - 0x1040))
			continue
		}
		if v, ok := punctuation[r]; ok {
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

// Contains reports whether s has any Burmese-block rune.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}

var consonants = map[rune]string{
	0x1000: "k",  // က
	0x1001: "hk", // ခ
	0x1002: "g",  // ဂ
	0x1003: "gh", // ဃ
	0x1004: "ng", // င
	0x1005: "s",  // စ
	0x1006: "hs", // ဆ
	0x1007: "z",  // ဇ
	0x1008: "jh", // ဈ
	0x1009: "ny", // ဉ
	0x100A: "ny", // ည
	0x100B: "t",  // ဋ
	0x100C: "ht", // ဌ
	0x100D: "d",  // ဍ
	0x100E: "dh", // ဎ
	0x100F: "n",  // ဏ
	0x1010: "t",  // တ
	0x1011: "ht", // ထ
	0x1012: "d",  // ဒ
	0x1013: "dh", // ဓ
	0x1014: "n",  // န
	0x1015: "p",  // ပ
	0x1016: "hp", // ဖ
	0x1017: "b",  // ဗ
	0x1018: "bh", // ဘ
	0x1019: "m",  // မ
	0x101A: "y",  // ယ
	0x101B: "r",  // ရ
	0x101C: "l",  // လ
	0x101D: "w",  // ဝ
	0x101E: "th", // သ
	0x101F: "h",  // ဟ
	0x1020: "l",  // ဠ
	0x1021: "ʾ",  // အ (glottal stop; often elided)
}

var independentVowels = map[rune]string{
	0x1023: "i",  // ဣ
	0x1024: "ī",  // ဤ
	0x1025: "u",  // ဥ
	0x1026: "ū",  // ဦ
	0x1027: "e",  // ဧ
	0x1029: "o",  // ဩ
	0x102A: "au", // ဪ
}

var vowelSigns = map[rune]string{
	0x102B: "ā",  // ◌ာ
	0x102C: "ā",  // ◌ါ
	0x102D: "i",  // ◌ိ
	0x102E: "ī",  // ◌ီ
	0x102F: "u",  // ◌ု
	0x1030: "ū",  // ◌ူ
	0x1031: "e",  // ◌ေ (prepositioned; we keep linear order)
	0x1032: "ai", // ◌ဲ
}

var medials = map[rune]string{
	0x103B: "y", // ◌ျ
	0x103C: "r", // ◌ြ
	0x103D: "w", // ◌ွ
	0x103E: "h", // ◌ှ
}

var tones = map[rune]string{
	0x1036: "ṁ", // ◌ံ anusvara
	0x1037: "",  // ◌့ low tone — drop
	0x1038: "ḥ", // ◌း visarga / high tone
	0x1039: "",  // ◌္ virama (older) — drop
}

var punctuation = map[rune]string{
	0x104A: ",",
	0x104B: ".",
}
