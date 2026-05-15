// Package hebrew provides SBL-style transliteration of Hebrew text to
// the Latin alphabet. The scheme is the de facto standard for Bible
// studies and academic OT/Tanakh work. Niqqud (vowel pointing) is
// consumed when present; the engine works on unpointed text too,
// producing consonant-only output for those passages.
//
// Cantillation marks (te'amim, the musical accents that decorate the
// Masoretic Text) are dropped — they carry recitation guidance, not
// pronunciation. Niqqud is normalized via Unicode NFD.
//
// Known simplifications (acceptable for v1):
//   - Kamatz (ָ) always renders as "ā" (no kamatz-katan / "o" disambiguation).
//   - Dagesh is dropped (consonants not doubled, BGD KPT not toggled).
//   - Vav-shuruq (וּ) renders as "wu" not "ū".
//   - Yod-with-hiriq (יִ) renders as "yi" not "î".
package hebr

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

const (
	BlockStart rune = 0x0590
	BlockEnd   rune = 0x05FF
	ShinDot    rune = 0x05C1 // ◌ׁ — right-side dot, marks shin
	SinDot     rune = 0x05C2 // ◌ׂ — left-side dot, marks sin
	Shin       rune = 0x05E9 // ש
)

// Transliterate returns the SBL-style romanization of s. The input is
// NFD-decomposed so that pointed glyphs split into base consonant +
// combining marks for uniform handling.
func Transliterate(s string) string {
	if s == "" {
		return s
	}
	rs := []rune(norm.NFD.String(s))
	var b strings.Builder
	b.Grow(len(s) * 2)

	for i := 0; i < len(rs); i++ {
		r := rs[i]
		// Consonants.
		if v, ok := consonants[r]; ok {
			// Shin needs the dot to disambiguate š vs ś.
			if r == Shin {
				v = shinSinFromMarks(rs, i)
			}
			b.WriteString(v)
			// Append any niqqud marks that follow.
			for j := i + 1; j < len(rs) && unicode.Is(unicode.Mn, rs[j]); j++ {
				if v, ok := niqqud[rs[j]]; ok {
					b.WriteString(v)
				}
				// Cantillation marks and shin/sin dots are silently
				// consumed without emitting anything.
			}
			// Skip over the combining-mark cluster.
			for i+1 < len(rs) && unicode.Is(unicode.Mn, rs[i+1]) {
				i++
			}
			continue
		}
		// Punctuation.
		if v, ok := punctuation[r]; ok {
			b.WriteString(v)
			continue
		}
		// Stray combining marks (no base in our consonant table): drop.
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		// In-block but unmapped: drop.
		if r >= BlockStart && r <= BlockEnd {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// shinSinFromMarks inspects the combining marks immediately after a
// shin (rs[i] == 0x05E9) and returns "š" if right-side dot follows, "ś"
// if left-side dot follows, "š" otherwise (most-common case).
func shinSinFromMarks(rs []rune, i int) string {
	for j := i + 1; j < len(rs) && unicode.Is(unicode.Mn, rs[j]); j++ {
		switch rs[j] {
		case ShinDot:
			return "š"
		case SinDot:
			return "ś"
		}
	}
	return "š"
}

// Contains reports whether s has at least one Hebrew-block rune.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}

// consonants maps Hebrew letters (including final forms) to their SBL
// Latin equivalent.
var consonants = map[rune]string{
	0x05D0: "ʾ", // א alef
	0x05D1: "b", // ב bet
	0x05D2: "g", // ג gimel
	0x05D3: "d", // ד dalet
	0x05D4: "h", // ה he
	0x05D5: "w", // ו vav
	0x05D6: "z", // ז zayin
	0x05D7: "ḥ", // ח chet
	0x05D8: "ṭ", // ט tet
	0x05D9: "y", // י yod
	0x05DA: "k", // ך final kaf
	0x05DB: "k", // כ kaf
	0x05DC: "l", // ל lamed
	0x05DD: "m", // ם final mem
	0x05DE: "m", // מ mem
	0x05DF: "n", // ן final nun
	0x05E0: "n", // נ nun
	0x05E1: "s", // ס samekh
	0x05E2: "ʿ", // ע ayin
	0x05E3: "p", // ף final pe
	0x05E4: "p", // פ pe
	0x05E5: "ṣ", // ץ final tsadi
	0x05E6: "ṣ", // צ tsadi
	0x05E7: "q", // ק qof
	0x05E8: "r", // ר resh
	0x05E9: "š", // ש shin (placeholder; shin/sin chosen by mark)
	0x05EA: "t", // ת tav
}

// niqqud maps vowel-point combining marks to Latin vowels.
var niqqud = map[rune]string{
	0x05B0: "ə", // ְ shewa (silent / vocal — we always render)
	0x05B1: "ě", // ֱ hataf segol
	0x05B2: "ă", // ֲ hataf patah
	0x05B3: "ŏ", // ֳ hataf qamatz
	0x05B4: "i", // ִ hiriq
	0x05B5: "ē", // ֵ tsere
	0x05B6: "e", // ֶ segol
	0x05B7: "a", // ַ patah
	0x05B8: "ā", // ָ qamatz
	0x05B9: "ō", // ֹ holam
	0x05BA: "ō", // ֺ holam haser for vav
	0x05BB: "u", // ֻ qubuts
	// 0x05BC: dagesh — dropped (would double consonant or alter b/g/d/k/p/t).
	// 0x05BD: meteg — dropped (cantillation/rhythm mark).
	// 0x05BE: maqaf — handled in punctuation.
	// 0x05BF: rafe — dropped.
}

// punctuation covers separators and verse markers.
var punctuation = map[rune]string{
	0x05BE: "-", // ־ maqaf (hyphen joiner)
	0x05C0: "|", // ׀ paseq
	0x05C3: ":", // ׃ sof passuq (end-of-verse)
	0x05C6: "",  // ׆ nun hafukha (drop)
	0x05F3: "'", // ׳ geresh
	0x05F4: `"`, // ״ gershayim
}
