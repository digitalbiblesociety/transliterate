// Package korean provides Revised Romanization (RR) of Korean Hangul
// (U+AC00-U+D7AF) to the Latin alphabet. RR is the official South Korean
// transliteration scheme.
//
// Hangul syllable blocks are mathematically decomposable into three
// jamo components (initial consonant, medial vowel, optional final
// consonant). We decompose, look up each component, and concatenate.
//
// Known simplifications:
//   - No context-sensitive g/k, d/t, b/p alternation: we use the
//     "after-vowel" form for initials and the "syllable-final" form
//     for finals, which produces consistent if slightly non-standard
//     spelling.
//   - No glottal assimilation between adjacent syllables.
package hang

import "strings"

const (
	BlockStart  rune = 0xAC00
	BlockEnd    rune = 0xD7A3
	SyllableMod rune = 28           // finals per medial
	MedialMod   rune = 21 * 28      // medials × finals per initial
	BlockSize   rune = 19 * 21 * 28 // 11172 syllables
)

// Transliterate returns the RR romanization of s. Non-Korean runes
// pass through unchanged.
func Transliterate(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) * 3)
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			idx := r - BlockStart
			initIdx := idx / MedialMod
			medIdx := (idx % MedialMod) / SyllableMod
			finIdx := idx % SyllableMod
			b.WriteString(initials[initIdx])
			b.WriteString(medials[medIdx])
			b.WriteString(finals[finIdx])
			continue
		}
		// Compatibility jamo / standalone jamo block — emit best-effort.
		if v, ok := compatJamo[r]; ok {
			b.WriteString(v)
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// Contains reports whether s has any Hangul syllable rune.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}

// initials — 19 leading consonants, ordered by jamo index 0..18.
var initials = [19]string{
	"g", "kk", "n", "d", "tt", "r", "m", "b", "pp",
	"s", "ss", "", "j", "jj", "ch", "k", "t", "p", "h",
}

// medials — 21 vowels, ordered by jamo index 0..20.
var medials = [21]string{
	"a", "ae", "ya", "yae", "eo", "e", "yeo", "ye",
	"o", "wa", "wae", "oe", "yo", "u", "wo", "we",
	"wi", "yu", "eu", "ui", "i",
}

// finals — 28 trailing consonants, ordered by jamo index 0..27.
// Index 0 = no final consonant.
var finals = [28]string{
	"", "k", "k", "ks", "n", "nj", "nh", "t", "l",
	"lk", "lm", "lb", "ls", "lt", "lp", "lh", "m",
	"p", "ps", "t", "t", "ng", "t", "t", "k", "t", "p", "t",
}

// compatJamo — Hangul Compatibility Jamo (U+3131..U+318E). Sometimes
// appears in dictionaries / phonetic notes.
var compatJamo = map[rune]string{
	0x3131: "g", 0x3132: "kk", 0x3134: "n", 0x3137: "d",
	0x3138: "tt", 0x3139: "r", 0x3141: "m", 0x3142: "b",
	0x3143: "pp", 0x3145: "s", 0x3146: "ss", 0x3147: "ng",
	0x3148: "j", 0x3149: "jj", 0x314A: "ch", 0x314B: "k",
	0x314C: "t", 0x314D: "p", 0x314E: "h",
	0x314F: "a", 0x3150: "ae", 0x3151: "ya", 0x3152: "yae",
	0x3153: "eo", 0x3154: "e", 0x3155: "yeo", 0x3156: "ye",
	0x3157: "o", 0x3158: "wa", 0x3159: "wae", 0x315A: "oe",
	0x315B: "yo", 0x315C: "u", 0x315D: "wo", 0x315E: "we",
	0x315F: "wi", 0x3160: "yu", 0x3161: "eu", 0x3162: "ui",
	0x3163: "i",
}
