// Package hebrew provides SBL Academic-style transliteration of Hebrew
// text to the Latin alphabet. The scheme is the de facto standard for
// Bible studies and academic OT/Tanakh work. Niqqud (vowel pointing) is
// consumed when present; the engine works on unpointed text too,
// producing consonant-only output for those passages.
//
// Cantillation marks (te'amim, the musical accents that decorate the
// Masoretic Text) are dropped — they carry recitation guidance, not
// pronunciation. Niqqud is normalized via Unicode NFD.
//
// Implemented SBL rules:
//   - Mater lectionis ligatures: hiriq-yod → î, tsere-yod → ê,
//     segol-yod → ê, holam-vav → ô, vav-with-dagesh (shuruq) → û,
//     word-final qamatz-he → â.
//   - Furtive patah: word-final patah before ḥ / ʿ / ה (with mappiq)
//     is emitted before the consonant, e.g. נֹחַ → nōaḥ.
//   - Dagesh chazaq (forte) doubles the consonant when preceded by a
//     vowel (e.g. הִנֵּה → hinnê, מַגָּל → maggāl). At word-end on a
//     consonant with silent shewa the doubling is not realised, per
//     standard SBL practice.
//   - Vocal vs silent shewa via a positional heuristic: initial → vocal;
//     after a long vowel or hataf → vocal; on a consonant with dagesh
//     forte → vocal; otherwise → silent.
//   - Qamatz-katan: U+05C7 always renders "o"; a regular qamatz is also
//     read as qatan when (a) the next cluster in the word has hataf-
//     qamatz, or (b) the word is a closed monosyllable joined by maqaf
//     (e.g. כָּל־ → kol-).
//   - Divine name (יהוה) is rendered "yhwh" regardless of pointing.
//
// Known simplifications:
//   - No accent-driven long-hiriq / long-qubuts macron (we emit "dāwid"
//     where some SBL guides give "dāwīd").
//   - No BGDKPT spirantization marking (SBL Academic default does not
//     mark it; the sblAcademicSpirantization variant which underlines
//     ḇ ḡ ḏ ḵ p̱ ṯ is not exposed).
//   - Sof passuq renders ":" rather than the SBL-default empty string,
//     because verse-end markers are useful in Bible-alignment pipelines.
package hebr

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

// Hebrew block range.
const (
	BlockStart rune = 0x0590
	BlockEnd   rune = 0x05FF
)

// Niqqud (vowel-point combining marks).
const (
	Sheva       rune = 0x05B0
	HatafSegol  rune = 0x05B1
	HatafPatah  rune = 0x05B2
	HatafQamats rune = 0x05B3
	Hiriq       rune = 0x05B4
	Tsere       rune = 0x05B5
	Segol       rune = 0x05B6
	Patah       rune = 0x05B7
	Qamats      rune = 0x05B8
	Holam       rune = 0x05B9
	HolamHaser  rune = 0x05BA
	Qubuts      rune = 0x05BB
	QamatsQatan rune = 0x05C7
)

// Non-vowel combining marks and structural punctuation.
const (
	Dagesh   rune = 0x05BC
	Meteg    rune = 0x05BD
	Maqaf    rune = 0x05BE
	Rafe     rune = 0x05BF
	Paseq    rune = 0x05C0
	ShinDot  rune = 0x05C1
	SinDot   rune = 0x05C2
	SofPasuq rune = 0x05C3
	NunHafu  rune = 0x05C6
	Geresh   rune = 0x05F3
	Gershyim rune = 0x05F4
)

// Hebrew letters.
const (
	Alef       rune = 0x05D0
	Bet        rune = 0x05D1
	Gimel      rune = 0x05D2
	Dalet      rune = 0x05D3
	He         rune = 0x05D4
	Vav        rune = 0x05D5
	Zayin      rune = 0x05D6
	Het        rune = 0x05D7
	Tet        rune = 0x05D8
	Yod        rune = 0x05D9
	FinalKaf   rune = 0x05DA
	Kaf        rune = 0x05DB
	Lamed      rune = 0x05DC
	FinalMem   rune = 0x05DD
	Mem        rune = 0x05DE
	FinalNun   rune = 0x05DF
	Nun        rune = 0x05E0
	Samekh     rune = 0x05E1
	Ayin       rune = 0x05E2
	FinalPe    rune = 0x05E3
	Pe         rune = 0x05E4
	FinalTsadi rune = 0x05E5
	Tsadi      rune = 0x05E6
	Qof        rune = 0x05E7
	Resh       rune = 0x05E8
	Shin       rune = 0x05E9
	Tav        rune = 0x05EA
)

// Transliterate returns the SBL Academic romanization of s. The input
// is NFD-decomposed so that pointed glyphs split into base consonant +
// combining marks for uniform handling.
func Transliterate(s string) string {
	if s == "" {
		return s
	}
	rs := []rune(norm.NFD.String(s))
	cs := parseClusters(rs)
	var b strings.Builder
	b.Grow(len(s) * 2)
	emit(cs, &b)
	return b.String()
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

// isHebrewLetter reports whether r is one of the 27 Hebrew base letters
// (including the five final forms).
func isHebrewLetter(r rune) bool {
	return r >= Alef && r <= Tav
}

// isVowelMark reports whether r is a niqqud combining mark (or the
// dedicated qamats-qatan U+05C7).
func isVowelMark(r rune) bool {
	switch r {
	case Sheva, HatafSegol, HatafPatah, HatafQamats,
		Hiriq, Tsere, Segol, Patah, Qamats,
		Holam, HolamHaser, Qubuts, QamatsQatan:
		return true
	}
	return false
}

// consonantLatin maps each Hebrew letter to its SBL Academic Latin
// equivalent. Shin is resolved separately via the shin/sin dot.
var consonantLatin = map[rune]string{
	Alef:       "ʾ",
	Bet:        "b",
	Gimel:      "g",
	Dalet:      "d",
	He:         "h",
	Vav:        "w",
	Zayin:      "z",
	Het:        "ḥ",
	Tet:        "ṭ",
	Yod:        "y",
	FinalKaf:   "k",
	Kaf:        "k",
	Lamed:      "l",
	FinalMem:   "m",
	Mem:        "m",
	FinalNun:   "n",
	Nun:        "n",
	Samekh:     "s",
	Ayin:       "ʿ",
	FinalPe:    "p",
	Pe:         "p",
	FinalTsadi: "ṣ",
	Tsadi:      "ṣ",
	Qof:        "q",
	Resh:       "r",
	Shin:       "š", // placeholder; resolved from shin/sin dot
	Tav:        "t",
}

// vowelLatin maps niqqud to SBL Academic vowels. Qamats may later be
// rewritten as "o" by the qamats-qatan rule; sheva may be elided by the
// silent-shewa rule.
var vowelLatin = map[rune]string{
	Sheva:       "ə",
	HatafSegol:  "ĕ",
	HatafPatah:  "ă",
	HatafQamats: "ŏ",
	Hiriq:       "i",
	Tsere:       "ē",
	Segol:       "e",
	Patah:       "a",
	Qamats:      "ā",
	Holam:       "ō",
	HolamHaser:  "ō",
	Qubuts:      "u",
	QamatsQatan: "o",
}

// punctuationLatin covers in-block separators and verse markers.
var punctuationLatin = map[rune]string{
	Maqaf:    "-",
	Paseq:    "",
	SofPasuq: ":",
	NunHafu:  "",
	Geresh:   "'",
	Gershyim: `"`,
}
