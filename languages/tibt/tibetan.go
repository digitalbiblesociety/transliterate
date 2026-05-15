// Package tibetan provides Wylie transliteration of Tibetan script
// (U+0F00-U+0FFF) to ASCII. Wylie was designed for 1:1 reversible
// ASCII encoding of Tibetan; we implement the standard scheme.
//
// The algorithm: each Tibetan letter has an inherent /a/ unless a
// vowel sign or virama follows. Subjoined consonants (U+0F90..U+0FBC)
// are emitted directly without an "a" since they form clusters.
package tibt

import "strings"

const (
	BlockStart rune = 0x0F00
	BlockEnd   rune = 0x0FFF
	Virama     rune = 0x0F84 // ◌྄ Tibetan halant
)

// Transliterate returns the Wylie romanization of s. Non-Tibetan runes
// pass through unchanged.
func Transliterate(s string) string {
	if s == "" {
		return s
	}
	rs := []rune(s)
	var b strings.Builder
	b.Grow(len(s) * 2)
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		// Head consonant: emit base + inherent /a/ unless next is a
		// vowel sign, virama, or subjoined consonant.
		if base, ok := headConsonant[r]; ok {
			b.WriteString(base)
			// Lookahead for vowel sign / virama / subjoined cluster.
			vowelSeen := false
			for i+1 < len(rs) {
				next := rs[i+1]
				if v, ok := vowelSign[next]; ok {
					b.WriteString(v)
					vowelSeen = true
					i++
					continue
				}
				if next == Virama {
					vowelSeen = true // virama suppresses inherent /a/
					i++
					continue
				}
				if sub, ok := subjoined[next]; ok {
					b.WriteString(sub)
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
		if v, ok := vowelSign[r]; ok {
			// Stray vowel sign (shouldn't happen in well-formed text).
			b.WriteString(v)
			continue
		}
		if v, ok := special[r]; ok {
			b.WriteString(v)
			continue
		}
		if r >= 0x0F20 && r <= 0x0F29 { // Tibetan digits
			b.WriteRune('0' + (r - 0x0F20))
			continue
		}
		if r >= BlockStart && r <= BlockEnd {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// Contains reports whether s has any Tibetan-block rune.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}

// headConsonant — head (full-width) Tibetan consonants.
var headConsonant = map[rune]string{
	0x0F40: "k",   0x0F41: "kh",  0x0F42: "g",   0x0F43: "gh",
	0x0F44: "ng",  0x0F45: "c",   0x0F46: "ch",  0x0F47: "j",
	0x0F49: "ny",  0x0F4A: "T",   0x0F4B: "Th",  0x0F4C: "D",
	0x0F4D: "Dh",  0x0F4E: "N",   0x0F4F: "t",   0x0F50: "th",
	0x0F51: "d",   0x0F52: "dh",  0x0F53: "n",   0x0F54: "p",
	0x0F55: "ph",  0x0F56: "b",   0x0F57: "bh",  0x0F58: "m",
	0x0F59: "ts",  0x0F5A: "tsh", 0x0F5B: "dz",  0x0F5C: "dzh",
	0x0F5D: "w",   0x0F5E: "zh",  0x0F5F: "z",   0x0F60: "'",
	0x0F61: "y",   0x0F62: "r",   0x0F63: "l",   0x0F64: "sh",
	0x0F65: "Sh",  0x0F66: "s",   0x0F67: "h",   0x0F68: "a",
	0x0F69: "ksh",
}

// subjoined — subjoined consonants (combined below a head). Same names
// as the head set but at +0x50 in the block. Wylie writes them inline.
var subjoined = map[rune]string{}

func init() {
	// Mirror head consonants into subjoined slots.
	for r, v := range headConsonant {
		subjoined[r+0x50] = v
	}
}

// vowelSign — dependent vowel signs.
var vowelSign = map[rune]string{
	0x0F71: "A",   // ◌ྰ long a marker (rare, mostly Sanskrit)
	0x0F72: "i",   // ◌ི
	0x0F73: "I",   // ◌ཱི long i
	0x0F74: "u",   // ◌ུ
	0x0F75: "U",   // ◌ཱུ long u
	0x0F7A: "e",   // ◌ེ
	0x0F7B: "ai",  // ◌ཻ
	0x0F7C: "o",   // ◌ོ
	0x0F7D: "au",  // ◌ཽ
}

// special — anusvara, visarga, punctuation, etc.
var special = map[rune]string{
	0x0F7E: "M",  // ◌ཾ anusvara
	0x0F7F: "H",  // ◌ཿ visarga
	0x0F0B: " ",  // ་ tsek (word separator)
	0x0F0C: " ",  // ༌ no-break tsek
	0x0F0D: "/",  // ། shad
	0x0F0E: "//", // ༎ double shad
	0x0F0F: "/!", // ༏ rin chen spungs shad
	0x0F11: "/!", // ༑
	0x0F14: "?",  // ༔ tsa-phru
}
