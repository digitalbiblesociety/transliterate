// Package japanese provides Hepburn romanization of Japanese kana
// (Hiragana U+3040-U+309F and Katakana U+30A0-U+30FF).
//
// Kanji (CJK Unified Ideographs) are NOT transliterated by this engine
// — doing so accurately requires morphological analysis (MeCab,
// Sudachi, etc.) since each kanji has multiple readings selected by
// context. Kanji pass through unchanged so the output is mixed-script:
// transliterated kana with kanji left as-is. This is honest about the
// engine's scope.
//
// Known simplifications:
//   - Small tsu (gemination marker) is dropped (no consonant doubling).
//   - Long-vowel chōonpu (ー) is dropped (no macron on the previous vowel).
package jpan

import "strings"

const (
	HiraStart rune = 0x3040
	HiraEnd   rune = 0x309F
	KataStart rune = 0x30A0
	KataEnd   rune = 0x30FF
)

// Transliterate returns the Hepburn romanization of kana in s. Non-kana
// runes (including kanji) pass through unchanged.
func Transliterate(s string) string {
	if s == "" {
		return s
	}
	rs := []rune(s)
	var b strings.Builder
	b.Grow(len(s) * 2)
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		// Check for digraph (small ya/yu/yo combinations).
		if i+1 < len(rs) {
			if v, ok := digraphs[[2]rune{r, rs[i+1]}]; ok {
				b.WriteString(v)
				i++
				continue
			}
		}
		if v, ok := kana[r]; ok {
			b.WriteString(v)
			continue
		}
		// In kana block but unmapped (small kana that didn't form a
		// digraph, chōonpu, etc.): drop.
		if r >= HiraStart && r <= KataEnd {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// Contains reports whether s has any Hiragana or Katakana rune.
func Contains(s string) bool {
	for _, r := range s {
		if (r >= HiraStart && r <= HiraEnd) || (r >= KataStart && r <= KataEnd) {
			return true
		}
	}
	return false
}

// kana maps each Hiragana and Katakana character to Hepburn.
var kana = map[rune]string{
	// Hiragana basic vowels.
	0x3042: "a", 0x3044: "i", 0x3046: "u", 0x3048: "e", 0x304A: "o",
	// Hiragana k-series.
	0x304B: "ka", 0x304D: "ki", 0x304F: "ku", 0x3051: "ke", 0x3053: "ko",
	0x304C: "ga", 0x304E: "gi", 0x3050: "gu", 0x3052: "ge", 0x3054: "go",
	// s-series.
	0x3055: "sa", 0x3057: "shi", 0x3059: "su", 0x305B: "se", 0x305D: "so",
	0x3056: "za", 0x3058: "ji", 0x305A: "zu", 0x305C: "ze", 0x305E: "zo",
	// t-series.
	0x305F: "ta", 0x3061: "chi", 0x3064: "tsu", 0x3066: "te", 0x3068: "to",
	0x3060: "da", 0x3062: "ji", 0x3065: "zu", 0x3067: "de", 0x3069: "do",
	// n-series.
	0x306A: "na", 0x306B: "ni", 0x306C: "nu", 0x306D: "ne", 0x306E: "no",
	// h-series.
	0x306F: "ha", 0x3072: "hi", 0x3075: "fu", 0x3078: "he", 0x307B: "ho",
	0x3070: "ba", 0x3073: "bi", 0x3076: "bu", 0x3079: "be", 0x307C: "bo",
	0x3071: "pa", 0x3074: "pi", 0x3077: "pu", 0x307A: "pe", 0x307D: "po",
	// m-series.
	0x307E: "ma", 0x307F: "mi", 0x3080: "mu", 0x3081: "me", 0x3082: "mo",
	// y-series.
	0x3084: "ya", 0x3086: "yu", 0x3088: "yo",
	// r-series.
	0x3089: "ra", 0x308A: "ri", 0x308B: "ru", 0x308C: "re", 0x308D: "ro",
	// w-series + n.
	0x308F: "wa", 0x3090: "wi", 0x3091: "we", 0x3092: "wo", 0x3093: "n",
	// Hiragana small a/i/u/e/o (rarely standalone, but possible).
	0x3041: "a", 0x3043: "i", 0x3045: "u", 0x3047: "e", 0x3049: "o",

	// Katakana basic vowels.
	0x30A2: "a", 0x30A4: "i", 0x30A6: "u", 0x30A8: "e", 0x30AA: "o",
	0x30AB: "ka", 0x30AD: "ki", 0x30AF: "ku", 0x30B1: "ke", 0x30B3: "ko",
	0x30AC: "ga", 0x30AE: "gi", 0x30B0: "gu", 0x30B2: "ge", 0x30B4: "go",
	0x30B5: "sa", 0x30B7: "shi", 0x30B9: "su", 0x30BB: "se", 0x30BD: "so",
	0x30B6: "za", 0x30B8: "ji", 0x30BA: "zu", 0x30BC: "ze", 0x30BE: "zo",
	0x30BF: "ta", 0x30C1: "chi", 0x30C4: "tsu", 0x30C6: "te", 0x30C8: "to",
	0x30C0: "da", 0x30C2: "ji", 0x30C5: "zu", 0x30C7: "de", 0x30C9: "do",
	0x30CA: "na", 0x30CB: "ni", 0x30CC: "nu", 0x30CD: "ne", 0x30CE: "no",
	0x30CF: "ha", 0x30D2: "hi", 0x30D5: "fu", 0x30D8: "he", 0x30DB: "ho",
	0x30D0: "ba", 0x30D3: "bi", 0x30D6: "bu", 0x30D9: "be", 0x30DC: "bo",
	0x30D1: "pa", 0x30D4: "pi", 0x30D7: "pu", 0x30DA: "pe", 0x30DD: "po",
	0x30DE: "ma", 0x30DF: "mi", 0x30E0: "mu", 0x30E1: "me", 0x30E2: "mo",
	0x30E4: "ya", 0x30E6: "yu", 0x30E8: "yo",
	0x30E9: "ra", 0x30EA: "ri", 0x30EB: "ru", 0x30EC: "re", 0x30ED: "ro",
	0x30EF: "wa", 0x30F0: "wi", 0x30F1: "we", 0x30F2: "wo", 0x30F3: "n",
	0x30A1: "a", 0x30A3: "i", 0x30A5: "u", 0x30A7: "e", 0x30A9: "o",
}

// digraphs covers small-ya/yu/yo combinations that produce palatalized
// syllables. The key is the {base, small-y-} pair.
var digraphs = map[[2]rune]string{
	// Hiragana small ya/yu/yo are at 0x3083 / 0x3085 / 0x3087.
	{0x304D, 0x3083}: "kya", {0x304D, 0x3085}: "kyu", {0x304D, 0x3087}: "kyo",
	{0x304E, 0x3083}: "gya", {0x304E, 0x3085}: "gyu", {0x304E, 0x3087}: "gyo",
	{0x3057, 0x3083}: "sha", {0x3057, 0x3085}: "shu", {0x3057, 0x3087}: "sho",
	{0x3058, 0x3083}: "ja", {0x3058, 0x3085}: "ju", {0x3058, 0x3087}: "jo",
	{0x3061, 0x3083}: "cha", {0x3061, 0x3085}: "chu", {0x3061, 0x3087}: "cho",
	{0x306B, 0x3083}: "nya", {0x306B, 0x3085}: "nyu", {0x306B, 0x3087}: "nyo",
	{0x3072, 0x3083}: "hya", {0x3072, 0x3085}: "hyu", {0x3072, 0x3087}: "hyo",
	{0x3073, 0x3083}: "bya", {0x3073, 0x3085}: "byu", {0x3073, 0x3087}: "byo",
	{0x3074, 0x3083}: "pya", {0x3074, 0x3085}: "pyu", {0x3074, 0x3087}: "pyo",
	{0x307F, 0x3083}: "mya", {0x307F, 0x3085}: "myu", {0x307F, 0x3087}: "myo",
	{0x308A, 0x3083}: "rya", {0x308A, 0x3085}: "ryu", {0x308A, 0x3087}: "ryo",
	// Katakana small ya/yu/yo are at 0x30E3 / 0x30E5 / 0x30E7.
	{0x30AD, 0x30E3}: "kya", {0x30AD, 0x30E5}: "kyu", {0x30AD, 0x30E7}: "kyo",
	{0x30B7, 0x30E3}: "sha", {0x30B7, 0x30E5}: "shu", {0x30B7, 0x30E7}: "sho",
	{0x30B8, 0x30E3}: "ja", {0x30B8, 0x30E5}: "ju", {0x30B8, 0x30E7}: "jo",
	{0x30C1, 0x30E3}: "cha", {0x30C1, 0x30E5}: "chu", {0x30C1, 0x30E7}: "cho",
	{0x30CB, 0x30E3}: "nya", {0x30CB, 0x30E5}: "nyu", {0x30CB, 0x30E7}: "nyo",
	{0x30D2, 0x30E3}: "hya", {0x30D2, 0x30E5}: "hyu", {0x30D2, 0x30E7}: "hyo",
	{0x30DF, 0x30E3}: "mya", {0x30DF, 0x30E5}: "myu", {0x30DF, 0x30E7}: "myo",
	{0x30EA, 0x30E3}: "rya", {0x30EA, 0x30E5}: "ryu", {0x30EA, 0x30E7}: "ryo",
}
