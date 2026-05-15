package arab

import "strings"

// Arabic combining marks (tashkeel) used by the rule engine.
const (
	arFatha      rune = 0x064E // ◌َ
	arDamma      rune = 0x064F // ◌ُ
	arKasra      rune = 0x0650 // ◌ِ
	arSukun      rune = 0x0652 // ◌ْ
	arShadda     rune = 0x0651 // ◌ّ
	arFathatan   rune = 0x064B // ◌ً
	arDammatan   rune = 0x064C // ◌ٌ
	arKasratan   rune = 0x064D // ◌ٍ
	arMaddah     rune = 0x0653 // ◌ٓ
	arHamzaAbove rune = 0x0654 // ◌ٔ
	arHamzaBelow rune = 0x0655 // ◌ٕ
	arSuperAlef  rune = 0x0670 // ◌ٰ
	arTatweel    rune = 0x0640 // ـ
)

// arabicConsBase maps an Arabic letter to its Latin base for tashkeel-aware
// transliteration. The scheme follows the conventions Google Translate's
// dt=rm endpoint emits: ع→e, ق→q, ث→th, etc. — phonemic, no diacritics.
var arabicConsBase = map[rune]string{
	'ب': "b",
	'ت': "t",
	'ث': "th",
	'ج': "j",
	'ح': "h",
	'خ': "kh",
	'د': "d",
	'ذ': "dh",
	'ر': "r",
	'ز': "z",
	'س': "s",
	'ش': "sh",
	'ص': "s",
	'ض': "d",
	'ط': "t",
	'ظ': "z",
	'ع': "e",
	'غ': "gh",
	'ف': "f",
	'ق': "q",
	'ك': "k",
	'ل': "l",
	'م': "m",
	'ن': "n",
	'ه': "h",
	'و': "w",
	'ي': "y",
	'ى': "a", // alef maksura
	// Bare alef and alef-wasla are silent length markers; hamza-bearing
	// alef variants are inconsistent in Google's output. Empirically the
	// apostrophe-emitting rune is إ (alef with hamza below) and ء (bare
	// hamza); the others (أ آ ؤ ئ) more often appear unmarked.
	'ا': "a",
	'ٱ': "a",
	'أ': "a",
	'إ': "'i",
	'آ': "a",
	'ؤ': "w",
	'ئ': "y",
	'ء': "'",
}

// arabicVowelLatin maps short vowel diacritics to their Latin form.
var arabicVowelLatin = map[rune]string{
	arFatha: "a",
	arDamma: "u",
	arKasra: "i",
}

// arabicTanweenLatin maps tanween (nunation) diacritics — appearing in
// indefinite words — to their Latin form.
var arabicTanweenLatin = map[rune]string{
	arFathatan: "an",
	arDammatan: "un",
	arKasratan: "in",
}

// isArabicMark reports whether r is a combining diacritic that follows an
// Arabic base letter.
func isArabicMark(r rune) bool {
	switch r {
	case arFatha, arDamma, arKasra, arSukun, arShadda,
		arFathatan, arDammatan, arKasratan,
		arMaddah, arHamzaAbove, arHamzaBelow, arSuperAlef:
		return true
	}
	return false
}

// isArabicConsonant reports whether r is an Arabic-script letter (not a mark).
func isArabicConsonant(r rune) bool {
	if r >= 0x0621 && r <= 0x063A {
		return true
	}
	if r >= 0x0641 && r <= 0x064A {
		return true
	}
	// Extended Arabic letters used in some scripts.
	if r >= 0x0671 && r <= 0x06D3 {
		return true
	}
	return false
}

// TransliterateTashkeel converts fully vocalized Arabic text to a
// Google-style romanization. Unvocalized words still produce reasonable
// output (consonants without vowels), but accuracy is highest on text
// with full tashkeel like the ARBERV Bible.
//
// Rules applied:
//   - Mater lectionis collapse: ا after fatha, ي after kasra, و after
//     damma are silent (the short vowel already gives the "long" sound).
//   - Shadda doubles the preceding consonant.
//   - Sukun produces no vowel.
//   - ة (ta marbuta) renders as "t" when followed by a vowel mark, "h"
//     when at word-end without vowel mark.
//   - The final short vowel of a word is dropped (pause state) so
//     "yaequb" not "yaequba".
//   - Tatweel (ـ) and various tone-like marks are stripped.
func TransliterateTashkeel(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) * 2)
	// Process word by word so we can detect end-of-word and drop the
	// final pause-state vowel.
	start := 0
	rs := []rune(s)
	for i := 0; i <= len(rs); i++ {
		if i == len(rs) || isWordSeparator(rs[i]) {
			if start < i {
				b.WriteString(transliterateArabicWord(rs[start:i]))
			}
			if i < len(rs) {
				r := rs[i]
				switch {
				case r >= 0x0660 && r <= 0x0669:
					b.WriteRune('0' + (r - 0x0660))
				case r >= 0x06F0 && r <= 0x06F9:
					b.WriteRune('0' + (r - 0x06F0))
				default:
					if mapped, ok := arabicPunctMap[r]; ok {
						b.WriteRune(mapped)
					} else {
						b.WriteRune(r)
					}
				}
			}
			start = i + 1
		}
	}
	return b.String()
}

func isWordSeparator(r rune) bool {
	if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
		return true
	}
	// Treat any non-Arabic, non-mark rune as a separator.
	return !isArabicConsonant(r) && !isArabicMark(r) && r != arTatweel
}

// arabicPunctMap normalizes Arabic punctuation to its ASCII equivalent
// — matching the output Google emits for the same input characters.
var arabicPunctMap = map[rune]rune{
	'،': ',',
	'؛': ';',
	'؟': '?',
	'٪': '%',
	'٫': '.', // decimal separator
	'٬': ',', // thousands separator
}

// transliterateArabicWord handles one whitespace-delimited Arabic word
// applying all the tashkeel rules.
func transliterateArabicWord(rs []rune) string {
	var b strings.Builder
	b.Grow(len(rs) * 2)

	// Collect emissions as a slice of "segments" so we can drop the
	// final pause-state vowel before joining.
	type seg struct {
		text    string
		vowel   rune // the short vowel emitted for this segment (0 if none)
		matered bool // a mater letter (ا/ي/و) consumed this segment's vowel
	}
	var segs []seg

	var lastVowel rune
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		if r == arTatweel {
			continue
		}
		// Stray combining mark with no preceding base — drop.
		if isArabicMark(r) {
			continue
		}
		// Mater lectionis: silent alef/ya/waw after matching short vowel,
		// provided the mater itself carries no vowel of its own. Mark the
		// preceding segment so pause-state stripping leaves its vowel
		// alone — the mater "claimed" it.
		if r == 'ا' && lastVowel == arFatha {
			i = skipMarks(rs, i)
			if n := len(segs); n > 0 {
				segs[n-1].matered = true
			}
			lastVowel = 0
			continue
		}
		if r == 'ي' && lastVowel == arKasra && !hasOwnVowel(rs, i) {
			i = skipMarks(rs, i)
			if n := len(segs); n > 0 {
				segs[n-1].matered = true
			}
			lastVowel = 0
			continue
		}
		if r == 'و' && lastVowel == arDamma && !hasOwnVowel(rs, i) {
			i = skipMarks(rs, i)
			if n := len(segs); n > 0 {
				segs[n-1].matered = true
			}
			lastVowel = 0
			continue
		}

		// Ta marbuta: t if followed by vowel mark, h if word-end with no marks.
		if r == 'ة' {
			marks := collectMarks(rs, i)
			i = marks.endIdx
			if marks.vowel != 0 || marks.tanween != 0 {
				segs = append(segs, seg{text: "t", vowel: marks.vowel})
				lastVowel = marks.vowel
			} else {
				segs = append(segs, seg{text: "h"})
				lastVowel = 0
			}
			continue
		}

		// Consonants and other letters.
		if base, ok := arabicConsBase[r]; ok {
			marks := collectMarks(rs, i)
			i = marks.endIdx
			emit := base
			if marks.shadda {
				emit = base + base
			}
			vowelLatin := ""
			if marks.vowel != 0 {
				vowelLatin = arabicVowelLatin[marks.vowel]
			} else if marks.tanween != 0 {
				vowelLatin = arabicTanweenLatin[marks.tanween]
			}
			segs = append(segs, seg{text: emit + vowelLatin, vowel: marks.vowel})
			if marks.vowel != 0 {
				lastVowel = marks.vowel
			} else if marks.tanween != 0 {
				lastVowel = marks.tanween
			} else {
				lastVowel = 0
			}
			continue
		}

		// Anything else (rare extended letters etc.): drop.
		lastVowel = 0
	}

	// Pause-state: drop the final short vowel of the last segment if it
	// has one. Skip if the segment's vowel was consumed by a mater (the
	// vowel is part of a long vowel and not a pause-droppable case mark).
	if len(segs) > 0 {
		last := &segs[len(segs)-1]
		if !last.matered {
			// Strip trailing short vowel (Google drops case markers and
			// tanween in pause state).
			for _, v := range []string{"an", "un", "in", "a", "u", "i"} {
				if strings.HasSuffix(last.text, v) && len(last.text) > len(v) {
					last.text = last.text[:len(last.text)-len(v)]
					break
				}
			}
		}
	}

	for _, s := range segs {
		b.WriteString(s.text)
	}
	return b.String()
}

// marksInfo captures the relevant diacritics that follow a base letter.
type marksInfo struct {
	shadda  bool
	vowel   rune // fatha/damma/kasra/sukun (0 if none)
	tanween rune // fathatan/dammatan/kasratan
	endIdx  int  // last consumed index
}

// collectMarks scans the combining marks immediately following the base
// letter at rs[i] and returns what it found.
func collectMarks(rs []rune, i int) marksInfo {
	out := marksInfo{endIdx: i}
	for j := i + 1; j < len(rs); j++ {
		m := rs[j]
		switch m {
		case arShadda:
			out.shadda = true
		case arFatha, arDamma, arKasra:
			out.vowel = m
		case arSukun:
			out.vowel = arSukun
		case arFathatan, arDammatan, arKasratan:
			out.tanween = m
		case arMaddah, arHamzaAbove, arHamzaBelow, arSuperAlef:
			// recognized but no Latin emission
		default:
			return out
		}
		out.endIdx = j
	}
	return out
}

// skipMarks returns the index of the last combining mark following rs[i],
// or i if no marks follow. Used for silent characters.
func skipMarks(rs []rune, i int) int {
	for j := i + 1; j < len(rs); j++ {
		if !isArabicMark(rs[j]) {
			return j - 1
		}
	}
	return len(rs) - 1
}

// hasOwnVowel reports whether the letter at rs[i] carries any vowel mark
// of its own (used to distinguish mater lectionis from consonantal use).
func hasOwnVowel(rs []rune, i int) bool {
	for j := i + 1; j < len(rs); j++ {
		m := rs[j]
		if !isArabicMark(m) {
			return false
		}
		switch m {
		case arFatha, arDamma, arKasra, arFathatan, arDammatan, arKasratan:
			return true
		}
	}
	return false
}
