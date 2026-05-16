// Package jpan provides Hepburn romanization of Japanese text. Kanji
// readings come from the kagome morphological analyzer (vendored under
// internal/kagome, MIT, with the IPA dictionary derived from
// mecab-ipadic-2.7.0-20070801). Kana (Hiragana U+3040-U+309F, Katakana
// U+30A0-U+30FF) is mapped directly via the table at the bottom of this
// file.
//
// Phenomena handled:
//   - Kanji is run through kagome to obtain a katakana reading, which is
//     then put through the same kana mapper as kana surfaces. Compound
//     reading selection is whatever IPA's primary reading is.
//   - Small tsu (っ/ッ) doubles the following consonant (a → kk, ch → tch).
//   - Chōonpu (ー) lengthens the preceding vowel to its macron form.
//   - Syllabic n (ん/ン) before b/m/p is emitted as "m" per traditional Hepburn.
//
// Known simplifications:
//   - Single-reading lookup: kagome's IPA dict gives one primary reading
//     per surface form. Names and rare words may not get their intended
//     reading.
//   - No long-o disambiguation: the reading オウ is rendered "ō" whether
//     it comes from a phonetic long /oː/ (e.g. 東京 → トーキョー) or a
//     historical /ou/ that's still pronounced as two morae (very rare).
//   - Particles は (ha→wa) and へ (he→e) are emitted by their spelling,
//     not their pronunciation. IPA's reading is the kana letter itself.
package jpan

import (
	"strings"
	"sync"

	"github.com/digitalbiblesociety/transliterate/internal/kagome/ipa"
	"github.com/digitalbiblesociety/transliterate/internal/kagome/tokenizer"
)

const (
	HiraStart rune = 0x3040
	HiraEnd   rune = 0x309F
	KataStart rune = 0x30A0
	KataEnd   rune = 0x30FF

	hiraSmallTsu rune = 0x3063
	kataSmallTsu rune = 0x30C3
	choonpu      rune = 0x30FC
	hiraN        rune = 0x3093
	kataN        rune = 0x30F3
)

// Transliterate returns the Hepburn romanization of s. The input is run
// through kagome's morphological analyzer; tokens whose surface lies in
// CJK / kana ranges are emitted via their katakana reading, everything
// else (Latin, digits, punctuation, USFM markup) passes through.
func Transliterate(s string) string {
	if s == "" {
		return s
	}
	if !needsTokenizer(s) {
		return transliterateKana(s)
	}
	tk := jpTokenizer()
	var b strings.Builder
	b.Grow(len(s) * 2)
	for _, tok := range tk.Tokenize(s) {
		if tok.Class == tokenizer.DUMMY {
			continue
		}
		if tok.Class == tokenizer.UNKNOWN {
			b.WriteString(tok.Surface)
			continue
		}
		if reading, ok := tok.Reading(); ok && reading != "*" {
			b.WriteString(transliterateKana(reading))
			continue
		}
		b.WriteString(transliterateKana(tok.Surface))
	}
	return b.String()
}

// Contains reports whether s has any Hiragana, Katakana, or CJK Unified
// Ideograph rune (the three ranges kagome can analyse).
func Contains(s string) bool {
	for _, r := range s {
		if (r >= HiraStart && r <= HiraEnd) || (r >= KataStart && r <= KataEnd) || (r >= 0x4E00 && r <= 0x9FFF) {
			return true
		}
	}
	return false
}

// needsTokenizer is true if the input has any kanji. Pure-kana strings
// take the cheap path and skip the tokenizer entirely.
func needsTokenizer(s string) bool {
	for _, r := range s {
		if r >= 0x4E00 && r <= 0x9FFF {
			return true
		}
		if r >= 0x3400 && r <= 0x4DBF { // CJK Extension A
			return true
		}
	}
	return false
}

// transliterateKana converts a kana string (typically the katakana
// reading from kagome) into Hepburn. The input must not contain kanji.
func transliterateKana(s string) string {
	rs := []rune(s)
	var b strings.Builder
	b.Grow(len(s) * 2)
	var lastVowel byte
	for i := 0; i < len(rs); i++ {
		r := rs[i]

		// Small tsu doubles the next syllable's first consonant.
		if r == hiraSmallTsu || r == kataSmallTsu {
			if next, j := nextSyllable(rs, i+1); next != "" {
				switch {
				case strings.HasPrefix(next, "ch"):
					b.WriteByte('t')
				default:
					b.WriteByte(next[0])
				}
				b.WriteString(next)
				lastVowel = vowelOf(next)
				i = j
			}
			continue
		}

		// Chōonpu lengthens the preceding vowel.
		if r == choonpu {
			if m, ok := macron[lastVowel]; ok {
				if cur := b.String(); len(cur) > 0 && cur[len(cur)-1] == lastVowel {
					// Replace the trailing ASCII vowel with its macron form.
					b.Reset()
					b.WriteString(cur[:len(cur)-1])
					b.WriteString(m)
					lastVowel = 0
				}
			}
			continue
		}

		// Syllabic n: emit "m" before b/m/p (traditional Hepburn), else "n".
		// Insert apostrophe between n and a following vowel/y to avoid
		// ambiguity (e.g. 単位 たんい → tan'i).
		if r == hiraN || r == kataN {
			if next, _ := nextSyllable(rs, i+1); next != "" {
				switch next[0] {
				case 'b', 'm', 'p':
					b.WriteByte('m')
				default:
					b.WriteByte('n')
					if next[0] == 'y' || isVowel(next[0]) {
						b.WriteByte('\'')
					}
				}
			} else {
				b.WriteByte('n')
			}
			lastVowel = 0
			continue
		}

		if i+1 < len(rs) {
			if v, ok := digraphs[[2]rune{r, rs[i+1]}]; ok {
				b.WriteString(v)
				lastVowel = vowelOf(v)
				i++
				continue
			}
		}
		if v, ok := kana[r]; ok {
			b.WriteString(v)
			lastVowel = vowelOf(v)
			continue
		}
		// In-block but unmapped (rare small kana not in a digraph,
		// archaic glyphs): drop. Out-of-block: passthrough.
		if r >= HiraStart && r <= KataEnd {
			continue
		}
		b.WriteRune(r)
		lastVowel = 0
	}
	return b.String()
}

// nextSyllable returns the Hepburn syllable starting at rs[i] and the
// index of the last rune consumed. Used by small-tsu and syllabic-n to
// peek without advancing the outer loop.
func nextSyllable(rs []rune, i int) (string, int) {
	if i >= len(rs) {
		return "", i - 1
	}
	if i+1 < len(rs) {
		if v, ok := digraphs[[2]rune{rs[i], rs[i+1]}]; ok {
			return v, i + 1
		}
	}
	if v, ok := kana[rs[i]]; ok {
		return v, i
	}
	return "", i - 1
}

func vowelOf(s string) byte {
	if s == "" {
		return 0
	}
	c := s[len(s)-1]
	if isVowel(c) {
		return c
	}
	return 0
}

func isVowel(c byte) bool {
	switch c {
	case 'a', 'i', 'u', 'e', 'o':
		return true
	}
	return false
}

var macron = map[byte]string{
	'a': "ā",
	'i': "ī",
	'u': "ū",
	'e': "ē",
	'o': "ō",
}

var (
	tkOnce sync.Once
	tk     *tokenizer.Tokenizer
	tkErr  error
)

func jpTokenizer() *tokenizer.Tokenizer {
	tkOnce.Do(func() {
		tk, tkErr = tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	})
	if tkErr != nil {
		// Dictionary failed to load — should never happen because the
		// dict is embedded and validated at compile time. Panicking
		// instead of silently degrading to kana-only keeps the failure
		// visible to callers.
		panic(tkErr)
	}
	return tk
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
	0x308F: "wa", 0x3090: "wi", 0x3091: "we", 0x3092: "wo",
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
	0x30EF: "wa", 0x30F0: "wi", 0x30F1: "we", 0x30F2: "wo",
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
	{0x30AE, 0x30E3}: "gya", {0x30AE, 0x30E5}: "gyu", {0x30AE, 0x30E7}: "gyo",
	{0x30B7, 0x30E3}: "sha", {0x30B7, 0x30E5}: "shu", {0x30B7, 0x30E7}: "sho",
	{0x30B8, 0x30E3}: "ja", {0x30B8, 0x30E5}: "ju", {0x30B8, 0x30E7}: "jo",
	{0x30C1, 0x30E3}: "cha", {0x30C1, 0x30E5}: "chu", {0x30C1, 0x30E7}: "cho",
	{0x30CB, 0x30E3}: "nya", {0x30CB, 0x30E5}: "nyu", {0x30CB, 0x30E7}: "nyo",
	{0x30D2, 0x30E3}: "hya", {0x30D2, 0x30E5}: "hyu", {0x30D2, 0x30E7}: "hyo",
	{0x30D3, 0x30E3}: "bya", {0x30D3, 0x30E5}: "byu", {0x30D3, 0x30E7}: "byo",
	{0x30D4, 0x30E3}: "pya", {0x30D4, 0x30E5}: "pyu", {0x30D4, 0x30E7}: "pyo",
	{0x30DF, 0x30E3}: "mya", {0x30DF, 0x30E5}: "myu", {0x30DF, 0x30E7}: "myo",
	{0x30EA, 0x30E3}: "rya", {0x30EA, 0x30E5}: "ryu", {0x30EA, 0x30E7}: "ryo",
}
