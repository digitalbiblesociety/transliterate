// Package greek provides SBL-style transliteration of Greek text to the
// Latin alphabet. The scheme is the de facto standard for Bible studies
// and academic NT/Koine work: η→ē, ω→ō, θ→th, φ→ph, χ→ch, ψ→ps, etc.
//
// Polytonic Greek (combining breathings, accents, iota subscript) is
// supported via Unicode NFD decomposition. Rough breathing prefixes "h"
// to the affected vowel; smooth breathing and accents are dropped.
//
// Known simplifications (acceptable for v1; can be tightened later):
//   - No nasal-velar assimilation: γγ→gg (not "ng"), γκ→gk, γχ→gch.
//   - υ always → y (no diphthong-aware "u" alternation).
//   - Initial ρ does not get "rh" prefix (the rough breathing on rho is
//     handled if present, but unmarked ῥ may slip through as "r").
//   - Iota subscript is dropped (some schemes render it as "i").
package grek

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

const (
	BlockStart   rune = 0x0370
	BlockEnd     rune = 0x03FF
	ExtStart     rune = 0x1F00 // Greek Extended (polytonic)
	ExtEnd       rune = 0x1FFF
	RoughMark    rune = 0x0314 // combining reversed comma above (rough breathing)
	SmoothMark   rune = 0x0313 // combining comma above (smooth breathing)
	IotaSubMark  rune = 0x0345 // combining Greek ypogegrammeni (iota subscript)
)

// Transliterate returns the SBL-style romanization of s. Non-Greek runes
// pass through unchanged. The input is first NFD-decomposed so that
// polytonic precomposed glyphs split into base letter + combining marks
// for uniform handling.
func Transliterate(s string) string {
	if s == "" {
		return s
	}
	rs := []rune(norm.NFD.String(s))
	var b strings.Builder
	b.Grow(len(s) * 2)

	var prevGreek rune
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		if v, ok := letterMap[r]; ok {
			// υ/Υ following a, e, o, or η in the same word forms a
			// diphthong and is conventionally rendered "u/U", not y/Y.
			if (r == 0x03C5 || r == 0x03A5) && isDiphthongTrigger(prevGreek) {
				v = caseLike("u", r)
			}
			// Vowel + rough breathing in the same combining-mark
			// cluster → prefix "h" (case-aware) to the vowel. When the
			// base letter is uppercase Greek we lowercase the Latin
			// vowel after the leading "H" so we get "Hē..." not "HĒ".
			if vowels[r] && hasRoughBreathing(rs, i) {
				if r >= 0x0391 && r <= 0x03A9 {
					v = "H" + strings.ToLower(v)
				} else {
					v = "h" + v
				}
			}
			b.WriteString(v)
			prevGreek = r
			// Skip combining marks that belong to this base letter.
			for i+1 < len(rs) && unicode.Is(unicode.Mn, rs[i+1]) {
				i++
			}
			continue
		}
		// Combining marks reached without a preceding base in our table:
		// drop silently.
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		// In-block but unmapped (rare): drop.
		if inGreek(r) {
			continue
		}
		// Anything else: passthrough, and break the diphthong chain.
		prevGreek = 0
		b.WriteRune(r)
	}
	return b.String()
}

// isDiphthongTrigger reports whether r is one of the Greek vowels that
// forms a diphthong with following υ/Υ: α, ε, ο, η (and uppercase forms).
func isDiphthongTrigger(r rune) bool {
	switch r {
	case 0x03B1, 0x0391, // α Α
		0x03B5, 0x0395, // ε Ε
		0x03BF, 0x039F, // ο Ο
		0x03B7, 0x0397: // η Η
		return true
	}
	return false
}

// caseLike returns base with the same case as the Greek letter at r.
// Used to keep "u" matched to "υ" vs "U" matched to "Υ" etc.
func caseLike(base string, r rune) string {
	if r >= 0x0391 && r <= 0x03A9 {
		return strings.ToUpper(base[:1]) + base[1:]
	}
	return base
}

// hasRoughBreathing peeks at the combining marks immediately after a
// base letter at rs[i] and reports whether U+0314 is among them.
func hasRoughBreathing(rs []rune, i int) bool {
	for j := i + 1; j < len(rs) && unicode.Is(unicode.Mn, rs[j]); j++ {
		if rs[j] == RoughMark {
			return true
		}
	}
	return false
}

func inGreek(r rune) bool {
	return (r >= BlockStart && r <= BlockEnd) || (r >= ExtStart && r <= ExtEnd)
}

// Contains reports whether s has at least one Greek-block (basic or
// extended) rune.
func Contains(s string) bool {
	for _, r := range s {
		if inGreek(r) {
			return true
		}
	}
	return false
}

// vowels marks which mapped Greek letters are vowels (subject to
// rough-breathing prefix). Includes both basic block and extended
// precomposed forms via NFD base letters.
var vowels = map[rune]bool{
	0x0391: true, 0x03B1: true, // Α α
	0x0395: true, 0x03B5: true, // Ε ε
	0x0397: true, 0x03B7: true, // Η η
	0x0399: true, 0x03B9: true, // Ι ι
	0x039F: true, 0x03BF: true, // Ο ο
	0x03A5: true, 0x03C5: true, // Υ υ
	0x03A9: true, 0x03C9: true, // Ω ω
	// Rho (ρ) traditionally carries rough breathing; treat as vowel-like
	// for the breathing prefix.
	0x03A1: true, 0x03C1: true,
}

// letterMap is the base SBL romanization for each Greek letter.
// Combining marks are handled separately by Transliterate.
var letterMap = map[rune]string{
	// Uppercase basic block.
	0x0391: "A",  // Α
	0x0392: "B",  // Β
	0x0393: "G",  // Γ
	0x0394: "D",  // Δ
	0x0395: "E",  // Ε
	0x0396: "Z",  // Ζ
	0x0397: "Ē",  // Η  (note: H-prefix case is handled by roughPrefix,
	//                    so capital Η with rough breathing yields "Hē"
	//                    not "HĒ" because we lowercase the rest.)
	0x0398: "Th", // Θ
	0x0399: "I",  // Ι
	0x039A: "K",  // Κ
	0x039B: "L",  // Λ
	0x039C: "M",  // Μ
	0x039D: "N",  // Ν
	0x039E: "X",  // Ξ
	0x039F: "O",  // Ο
	0x03A0: "P",  // Π
	0x03A1: "R",  // Ρ
	0x03A3: "S",  // Σ
	0x03A4: "T",  // Τ
	0x03A5: "Y",  // Υ
	0x03A6: "Ph", // Φ
	0x03A7: "Ch", // Χ
	0x03A8: "Ps", // Ψ
	0x03A9: "Ō",  // Ω
	// Lowercase basic block.
	0x03B1: "a",  // α
	0x03B2: "b",  // β
	0x03B3: "g",  // γ
	0x03B4: "d",  // δ
	0x03B5: "e",  // ε
	0x03B6: "z",  // ζ
	0x03B7: "ē",  // η
	0x03B8: "th", // θ
	0x03B9: "i",  // ι
	0x03BA: "k",  // κ
	0x03BB: "l",  // λ
	0x03BC: "m",  // μ
	0x03BD: "n",  // ν
	0x03BE: "x",  // ξ
	0x03BF: "o",  // ο
	0x03C0: "p",  // π
	0x03C1: "r",  // ρ
	0x03C2: "s",  // ς final sigma
	0x03C3: "s",  // σ
	0x03C4: "t",  // τ
	0x03C5: "y",  // υ
	0x03C6: "ph", // φ
	0x03C7: "ch", // χ
	0x03C8: "ps", // ψ
	0x03C9: "ō",  // ω
	// Archaic / variant letters occasionally seen in biblical texts.
	0x03DC: "F",  // Ϝ digamma
	0x03DD: "f",  // ϝ
}
