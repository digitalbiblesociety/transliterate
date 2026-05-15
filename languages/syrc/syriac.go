// Package syriac provides ISO 233-3 transliteration of the Syriac
// script (U+0700-U+074F) to the Latin alphabet. Syriac is an Aramaic-
// family abjad used for Eastern Christian liturgical and biblical
// texts (Peshitta, NT in Syriac).
package syrc

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

const (
	BlockStart rune = 0x0700
	BlockEnd   rune = 0x074F
)

// Transliterate returns the ISO 233-3 romanization of s. Vowel
// pointing (when present) is consumed; combining marks not in our
// table are dropped. Non-Syriac runes pass through unchanged.
func Transliterate(s string) string {
	if s == "" {
		return s
	}
	rs := []rune(norm.NFD.String(s))
	var b strings.Builder
	b.Grow(len(s) * 2)
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		if v, ok := consonants[r]; ok {
			b.WriteString(v)
			// Append vowel marks in the trailing cluster.
			for i+1 < len(rs) && unicode.Is(unicode.Mn, rs[i+1]) {
				i++
				if v, ok := vowelMarks[rs[i]]; ok {
					b.WriteString(v)
				}
			}
			continue
		}
		if v, ok := punctuation[r]; ok {
			b.WriteString(v)
			continue
		}
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		if r >= BlockStart && r <= BlockEnd {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// Contains reports whether s has any Syriac-block rune.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}

var consonants = map[rune]string{
	0x0710: "ʾ", // ܐ alaph
	0x0712: "b", // ܒ beth
	0x0713: "g", // ܓ gamal
	0x0714: "g̱", // ܔ Garshuni gamal
	0x0715: "d", // ܕ dalath
	0x0716: "d̥", // ܖ dotless dalath
	0x0717: "h", // ܗ he
	0x0718: "w", // ܘ waw
	0x0719: "z", // ܙ zayn
	0x071A: "ḥ", // ܚ heth
	0x071B: "ṭ", // ܛ teth
	0x071C: "ṭ̄", // ܜ Garshuni teth
	0x071D: "y", // ܝ yodh
	0x071E: "ȳ", // ܞ yodh-yodh
	0x071F: "k", // ܟ kaph
	0x0720: "l", // ܠ lamadh
	0x0721: "m", // ܡ mim
	0x0722: "n", // ܢ nun
	0x0723: "s", // ܣ semkath
	0x0724: "s̄", // ܤ final semkath
	0x0725: "ʿ", // ܥ ʿe
	0x0726: "p", // ܦ pe
	0x0727: "p̄", // ܧ reversed pe
	0x0728: "ṣ", // ܨ sadhe
	0x0729: "q", // ܩ qaph
	0x072A: "r", // ܪ rish
	0x072B: "š", // ܫ shin
	0x072C: "t", // ܬ taw
	0x072D: "b̥", // ܭ Persian bheth
	0x072E: "g̱", // ܮ Persian ghamal
	0x072F: "d̥", // ܯ Persian dhalath
}

// vowelMarks — Syriac vowel pointing (East and West Syrian systems).
var vowelMarks = map[rune]string{
	0x0730: "a",  // ◌ܰ pthaha
	0x0731: "ā",  // ◌ܱ pthaha esasa
	0x0732: "ā",  // ◌ܲ pthaha dotted
	0x0733: "ā",  // ◌ܳ zqapha
	0x0734: "ā",  // ◌ܴ zqapha esasa
	0x0735: "ā",  // ◌ܵ zqapha dotted
	0x0736: "e",  // ◌ܶ rbasa
	0x0737: "e",  // ◌ܷ rbasa esasa
	0x0738: "ə",  // ◌ܸ dotted zlama horizontal
	0x0739: "ē",  // ◌ܹ dotted zlama angular
	0x073A: "i",  // ◌ܺ hbasa above
	0x073B: "i",  // ◌ܻ hbasa below
	0x073C: "ū",  // ◌ܼ hbasa-esasa dotted
	0x073D: "u",  // ◌ܽ esasa above
	0x073E: "u",  // ◌ܾ esasa below
}

var punctuation = map[rune]string{
	0x0700: ".",   // ܀ end of paragraph
	0x0701: ":",   // ܁ supralinear full stop
	0x0702: ":",   // ܂ sublinear full stop
	0x0703: ":",   // ܃ supralinear colon
	0x0704: ":",   // ܄ sublinear colon
	0x0705: ",",   // ܅
	0x0706: ",",   // ܆
	0x0707: "?",   // ܇
	0x0708: "?",   // ܈
	0x0709: ".",   // ܉
	0x070A: "+",   // ܊ contraction
	0x070C: "-",   // ܌ harklean
}
