// Package brahmic provides ISO 15919-style transliteration of Brahmic-family
// scripts (Devanagari, Bengali, Gurmukhi, Gujarati, Oriya, Tamil, Telugu,
// Kannada, Malayalam, Sinhala, Javanese, Sundanese, Balinese, Buginese,
// Batak) to the Latin alphabet.
//
// The transliteration algorithm is identical across scripts because Brahmic
// scripts share the same phonological model: each consonant carries an
// inherent /a/ that is suppressed by a virama (halant) or replaced by a
// vowel sign (matra). Only the codepoint tables differ from script to
// script, so each language is defined by a Script descriptor.
//
// Reference: ISO 15919 (2001) "Transliteration of Devanagari and related
// Indic scripts into Latin characters."
package brahmic

import "strings"

// Script describes one Brahmic-family writing system: the codepoint range
// it occupies, the lookup tables for its letters, and any script-specific
// quirks.
type Script struct {
	Name string

	// BlockStart and BlockEnd bound the script's Unicode block, used for
	// detection and for swallowing unhandled in-block runes (instead of
	// passing them through as garbage).
	BlockStart rune
	BlockEnd   rune

	// Virama lists the codepoint(s) that suppress the inherent /a/ of the
	// preceding consonant. Most scripts have exactly one; a few (e.g. Batak,
	// which has both pangolat and pangonangon) have alternates. Empty/nil
	// means the script has no virama at all (e.g. Buginese).
	Virama []rune

	// DigitStart is the codepoint of the script's "0" glyph; consecutive
	// codepoints 0–9 follow. Set to 0 if the script has no native digits.
	DigitStart rune

	// ConsonantBase maps a consonant rune to its Latin base WITHOUT the
	// inherent /a/. The algorithm appends "a" if no vowel sign / virama
	// follows.
	ConsonantBase map[rune]string

	// IndependentVowel maps a standalone vowel rune to its Latin form.
	IndependentVowel map[rune]string

	// VowelSign (matra) maps a dependent vowel sign to its Latin form;
	// applied to the preceding consonant in place of the inherent /a/.
	VowelSign map[rune]string

	// Special covers anusvara, visarga, avagraha, and any other in-block
	// runes that need a Latin equivalent independent of consonant context.
	Special map[rune]string

	// ConsonantBaseWithChiller maps a final-consonant form (such as
	// Malayalam chillu letters) directly to its Latin base without an
	// inherent /a/ and without consuming a following vowel-sign. Optional;
	// nil for scripts that have no such forms.
	ChilluLike map[rune]string

	// Nukta is the script's nukta combining mark (◌़ in Devanagari, ◌଼ in
	// Bengali/Oriya, etc.). Zero if the script has no nukta.
	Nukta rune

	// NuktaCombine maps a (consonant + nukta) combination to a modified
	// Latin base. Used when a base consonant is followed by Nukta in the
	// decomposed form. Pre-composed codepoints (e.g. U+0958..U+095F in
	// Devanagari) live in ConsonantBase as usual.
	NuktaCombine map[rune]string

	// Medial maps medial consonant signs (e.g. Javanese pengkal / cakra)
	// to their Latin form. Consumed greedily after a base consonant —
	// inserted between the base and any vowel sign or virama. Nil for
	// scripts without medials.
	Medial map[rune]string
}

// Transliterate returns the ISO 15919 romanization of s for the given
// script. Runes outside the script's block pass through unchanged so the
// function is safe to call on USFM lines, mixed-script text, etc.
func Transliterate(s string, sc *Script) string {
	if s == "" {
		return s
	}
	rs := []rune(s)
	var b strings.Builder
	b.Grow(len(s) * 2)

	for i := 0; i < len(rs); i++ {
		r := rs[i]

		if v, ok := sc.ChilluLike[r]; ok {
			b.WriteString(v)
			continue
		}
		if base, ok := sc.ConsonantBase[r]; ok {
			// Decomposed nukta: consonant + nukta = modified consonant.
			consumedNukta := 0
			if sc.Nukta != 0 && i+1 < len(rs) && rs[i+1] == sc.Nukta {
				if mod, ok := sc.NuktaCombine[r]; ok {
					base = mod
				}
				consumedNukta = 1
			}
			b.WriteString(base)

			// Consume zero or more medial consonants (e.g. Javanese pengkal,
			// cakra) attached to this base.
			j := i + 1 + consumedNukta
			for j < len(rs) {
				if m, ok := sc.Medial[rs[j]]; ok {
					b.WriteString(m)
					j++
					continue
				}
				break
			}

			if j < len(rs) {
				next := rs[j]
				if sc.isVirama(next) {
					i = j
					continue
				}
				if v, ok := sc.VowelSign[next]; ok {
					b.WriteString(v)
					i = j
					continue
				}
			}
			b.WriteString("a")
			i = j - 1
			continue
		}
		if v, ok := sc.IndependentVowel[r]; ok {
			// A vowel sign immediately after an independent vowel
			// overrides its inherent value — Letter A acts as a carrier
			// in scripts (e.g. Buginese) that lack independent codepoints
			// for /i/, /u/, /e/, /o/.
			if i+1 < len(rs) {
				if sign, ok := sc.VowelSign[rs[i+1]]; ok {
					b.WriteString(sign)
					i++
					continue
				}
			}
			b.WriteString(v)
			continue
		}
		if v, ok := sc.Special[r]; ok {
			b.WriteString(v)
			continue
		}
		if sc.DigitStart != 0 && r >= sc.DigitStart && r <= sc.DigitStart+9 {
			b.WriteRune('0' + (r - sc.DigitStart))
			continue
		}
		// Unhandled in-block runes (stray virama, archaic glyphs we haven't
		// mapped, formatting marks): drop. Out-of-block: passthrough.
		if sc.isVirama(r) {
			continue
		}
		if r >= sc.BlockStart && r <= sc.BlockEnd {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// Contains reports whether s has at least one rune inside sc's block.
func Contains(s string, sc *Script) bool {
	for _, r := range s {
		if r >= sc.BlockStart && r <= sc.BlockEnd {
			return true
		}
	}
	return false
}

// isVirama reports whether r is one of the script's virama codepoints.
func (sc *Script) isVirama(r rune) bool {
	for _, v := range sc.Virama {
		if r == v {
			return true
		}
	}
	return false
}
