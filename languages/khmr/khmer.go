// Package khmer provides simplified transliteration of the Khmer
// script (U+1780-U+17FF) to the Latin alphabet.
//
// Khmer is an Abugida with two "registers" (a-series and o-series)
// that determine the inherent vowel of each consonant. For Bible-
// transliteration purposes we elide the two-register distinction and
// use a single inherent /a/, producing readable but not phonemically
// perfect output.
//
// Known simplifications:
//   - Inherent vowel always rendered "a" (not o-series "o").
//   - COENG (U+17D2) subjoining produces inline cluster without
//     stacking indication.
//   - Various diacritics (nikahit, reahmuk, etc.) handled minimally.
package khmr

import "strings"

const (
	BlockStart rune = 0x1780
	BlockEnd   rune = 0x17FF
	Coeng      rune = 0x17D2 // subscript marker
)

// Transliterate returns the romanization of s. Non-Khmer runes pass
// through unchanged.
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
			vowelSeen := false
			for i+1 < len(rs) {
				next := rs[i+1]
				if next == Coeng {
					// Subjoined consonant: emit its base without inherent /a/.
					if i+2 < len(rs) {
						if sub, ok := consonants[rs[i+2]]; ok {
							b.WriteString(sub)
							i += 2
							continue
						}
					}
					i++ // skip stray coeng
					continue
				}
				if v, ok := vowelSigns[next]; ok {
					b.WriteString(v)
					vowelSeen = true
					i++
					continue
				}
				if v, ok := signs[next]; ok {
					b.WriteString(v)
					i++
					continue
				}
				break
			}
			if !vowelSeen {
				b.WriteString("a")
			}
			continue
		}
		if v, ok := independentVowels[r]; ok {
			b.WriteString(v)
			continue
		}
		if r >= 0x17E0 && r <= 0x17E9 { // Khmer digits
			b.WriteRune('0' + (r - 0x17E0))
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

// Contains reports whether s has any Khmer-block rune.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}

var consonants = map[rune]string{
	0x1780: "k",  // ក
	0x1781: "kh", // ខ
	0x1782: "g",  // គ
	0x1783: "gh", // ឃ
	0x1784: "ng", // ង
	0x1785: "c",  // ច
	0x1786: "ch", // ឆ
	0x1787: "j",  // ជ
	0x1788: "jh", // ឈ
	0x1789: "ny", // ញ
	0x178A: "ṭ",  // ដ
	0x178B: "ṭh", // ឋ
	0x178C: "ḍ",  // ឌ
	0x178D: "ḍh", // ឍ
	0x178E: "ṇ",  // ណ
	0x178F: "t",  // ត
	0x1790: "th", // ថ
	0x1791: "d",  // ទ
	0x1792: "dh", // ធ
	0x1793: "n",  // ន
	0x1794: "p",  // ប
	0x1795: "ph", // ផ
	0x1796: "b",  // ព
	0x1797: "bh", // ភ
	0x1798: "m",  // ម
	0x1799: "y",  // យ
	0x179A: "r",  // រ
	0x179B: "l",  // ល
	0x179C: "v",  // វ
	0x179D: "ś",  // ឝ
	0x179E: "ṣ",  // ឞ
	0x179F: "s",  // ស
	0x17A0: "h",  // ហ
	0x17A1: "ḷ",  // ឡ
	0x17A2: "ʾ",  // អ
}

var independentVowels = map[rune]string{
	0x17A3: "a",  // ឣ
	0x17A4: "ā",  // ឤ
	0x17A5: "i",  // ឥ
	0x17A6: "ī",  // ឦ
	0x17A7: "u",  // ឧ
	0x17A8: "ū",  // ឨ
	0x17A9: "ū",  // ឩ
	0x17AA: "ūv", // ឪ
	0x17AB: "r̥",  // ឫ
	0x17AC: "r̥̄", // ឬ
	0x17AD: "l̥",  // ឭ
	0x17AE: "l̥̄", // ឮ
	0x17AF: "e",  // ឯ
	0x17B0: "ai", // ឰ
	0x17B1: "o",  // ឱ
	0x17B2: "o",  // ឲ
	0x17B3: "au", // ឳ
}

var vowelSigns = map[rune]string{
	0x17B6: "ā",  // ◌ា
	0x17B7: "i",  // ◌ិ
	0x17B8: "ī",  // ◌ី
	0x17B9: "ɨ",  // ◌ឹ
	0x17BA: "ɨ̄",  // ◌ឺ
	0x17BB: "u",  // ◌ុ
	0x17BC: "ū",  // ◌ូ
	0x17BD: "ua", // ◌ួ
	0x17BE: "œ",  // ◌ើ
	0x17BF: "ɨə", // ◌ឿ
	0x17C0: "iə", // ◌ៀ
	0x17C1: "e",  // ◌េ
	0x17C2: "ae", // ◌ែ
	0x17C3: "ai", // ◌ៃ
	0x17C4: "o",  // ◌ោ
	0x17C5: "au", // ◌ៅ
}

var signs = map[rune]string{
	0x17C6: "ṁ", // ◌ំ nikahit (anusvara)
	0x17C7: "ḥ", // ◌ះ reahmuk (visarga)
	0x17C8: "",  // ◌ៈ yuukaleapintu
	0x17C9: "",  // ◌៉ muusikatoan
	0x17CA: "",  // ◌៊ triisap
	0x17CB: "",  // ◌់ bantoc
	0x17CC: "",  // ◌៌ robat
	0x17CD: "",  // ◌៍ toandakhiat (kills consonant)
	0x17CE: "",  // ◌៎ kakabat
	0x17CF: "",  // ◌៏ ahsda
	0x17D0: "",  // ◌័ samyok
}

var punctuation = map[rune]string{
	0x17D4: ".",  // ។
	0x17D5: ":",  // ៕
	0x17D6: ":",  // ៖
	0x17D8: "",   // ៘
	0x17D9: "",   // ៙
	0x17DA: "",   // ៚
	0x17DB: "$",  // ៛ riel sign
	0x17DC: "ʾ",  // ៜ avakrahasanya
}
