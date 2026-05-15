// Package thai provides Royal Thai General System of Transcription (RTGS)
// romanization of the Thai script (U+0E00..U+0E7F).
//
// The rule set is a Go port of the rule-based `royin` algorithm from
// PyThaiNLP (Apache 2.0, https://github.com/PyThaiNLP/pythainlp), which
// itself implements the Thai Royal Institute's official RTGS scheme.
// Output is plain ASCII — no diacritics, no tone marks, no length marks
// (Thai RTGS by design folds short/long vowel length).
//
// Phenomena handled:
//   - Pre-positioned vowels (เ แ โ ใ ไ) reordered to follow their
//     consonant in the output.
//   - Common multi-rune vowel patterns: เ◌า (ao), เ◌าะ (o), เ◌ะ (e
//     short), แ◌ะ (ae short), โ◌ะ (o short).
//   - Thanthakhat (◌์) silences the consonant beneath it.
//   - Silent leading ห before a sonorant (ง น ม ย ร ล ว) at the start
//     of a syllable.
//   - Tone marks (◌่ ◌้ ◌๊ ◌๋), maitaikhu (◌็), and various
//     punctuation/non-pronunciation marks are dropped.
//
// Known simplifications:
//   - No syllable segmentation. Thai is unsegmented in writing and many
//     RTGS rules depend on syllable boundary; we approximate.
//   - Final-position consonant forms (e.g. บ → "p" at syllable end) are
//     not applied — we always emit the initial form. PyThaiNLP itself
//     punts on edge cases here; without segmentation full accuracy is
//     infeasible.
//   - Tone-class aware vowels (e.g. ค and ก both → "k" but read with
//     different tones) are flattened to a single consonant form.
package thai

import "strings"

const (
	BlockStart rune = 0x0E00
	BlockEnd   rune = 0x0E7F

	thanthakhat rune = 0x0E4C // ◌์ silencer
	maitaikhu   rune = 0x0E47 // ◌็ vowel shortener
	hoHip       rune = 0x0E2B // ห
)

// consonants — initial-position emission. RTGS conflates the four-letter
// classes (high/mid/low/extra) at the consonant level; differences show
// in tones, which we don't represent in plain ASCII.
var consonants = map[rune]string{
	0x0E01: "k",  // ก
	0x0E02: "kh", // ข
	0x0E03: "kh", // ฃ
	0x0E04: "kh", // ค
	0x0E05: "kh", // ฅ
	0x0E06: "kh", // ฆ
	0x0E07: "ng", // ง
	0x0E08: "ch", // จ
	0x0E09: "ch", // ฉ
	0x0E0A: "ch", // ช
	0x0E0B: "s",  // ซ
	0x0E0C: "ch", // ฌ
	0x0E0D: "y",  // ญ
	0x0E0E: "d",  // ฎ
	0x0E0F: "t",  // ฏ
	0x0E10: "th", // ฐ
	0x0E11: "th", // ฑ
	0x0E12: "th", // ฒ
	0x0E13: "n",  // ณ
	0x0E14: "d",  // ด
	0x0E15: "t",  // ต
	0x0E16: "th", // ถ
	0x0E17: "th", // ท
	0x0E18: "th", // ธ
	0x0E19: "n",  // น
	0x0E1A: "b",  // บ
	0x0E1B: "p",  // ป
	0x0E1C: "ph", // ผ
	0x0E1D: "f",  // ฝ
	0x0E1E: "ph", // พ
	0x0E1F: "f",  // ฟ
	0x0E20: "ph", // ภ
	0x0E21: "m",  // ม
	0x0E22: "y",  // ย
	0x0E23: "r",  // ร
	0x0E25: "l",  // ล
	0x0E27: "w",  // ว
	0x0E28: "s",  // ศ
	0x0E29: "s",  // ษ
	0x0E2A: "s",  // ส
	0x0E2B: "h",  // ห
	0x0E2C: "l",  // ฬ
	0x0E2D: "",   // อ — glottal/vowel carrier; emits no consonant
	0x0E2E: "h",  // ฮ
}

// prepositioned vowels: visually before the consonant but pronounced
// after. Standalone Latin form is used when no consonant follows.
var prepositioned = map[rune]string{
	0x0E40: "e",  // เ
	0x0E41: "ae", // แ
	0x0E42: "o",  // โ
	0x0E43: "ai", // ใ
	0x0E44: "ai", // ไ
}

// Post-consonant vowel signs. Plain ASCII; length distinctions are dropped.
var vowelSigns = map[rune]string{
	0x0E30: "a",  // ◌ะ
	0x0E31: "a",  // ◌ั
	0x0E32: "a",  // า
	0x0E33: "am", // ◌ำ
	0x0E34: "i",  // ◌ิ
	0x0E35: "i",  // ◌ี
	0x0E36: "ue", // ◌ึ
	0x0E37: "ue", // ◌ื
	0x0E38: "u",  // ◌ุ
	0x0E39: "u",  // ◌ู
}

// sonorant: consonants that silence a preceding ห at syllable start.
var sonorant = map[rune]bool{
	0x0E07: true, // ง
	0x0E19: true, // น
	0x0E21: true, // ม
	0x0E22: true, // ย
	0x0E23: true, // ร
	0x0E25: true, // ล
	0x0E27: true, // ว
}

// toneMarks are dropped on sight and skipped over when matching trailing
// vowel patterns (so e.g. โต๊ะ still resolves the โCะ → "o" pattern).
var toneMarks = map[rune]bool{
	0x0E48: true, // tone mark 1
	0x0E49: true, // tone mark 2
	0x0E4A: true, // tone mark 3
	0x0E4B: true, // tone mark 4
}

// ignored is dropped on sight in the main loop. Includes tone marks plus
// maitaikhu (the vowel shortener — pattern detection handles it
// explicitly, so when it shows up here it had no pattern context) and
// other formatting marks. Native pause markers (ฯ ๏ ๚ ๛) are NOT in
// this set — they map to Western punctuation via pauseMarkers so the
// downstream alignment tooling (e.g. Aeneas) sees explicit phrase
// boundaries.
var ignored = map[rune]bool{
	0x0E3A: true, // phinthu
	0x0E47: true, // maitaikhu (pattern-detected as part of เC็ short e)
	0x0E48: true, 0x0E49: true, 0x0E4A: true, 0x0E4B: true, // tone marks
	0x0E4D: true, // niggahita
	0x0E4E: true, // yamakkan
}

// pauseMarkers maps Thai punctuation that signals an audible pause to
// the Western punctuation of equivalent strength. The mapping intent:
//
//	ฯ paiyannoi → ","  medium pause / abbreviation / soft sentence end
//	๏ fongman   → "."  section/paragraph opener (boundary precedes it)
//	๚ angkhankhu → "." strong: end of stanza or chapter
//	๛ khomut    → "."  very strong: end of text / chapter
//
// Aeneas (and other forced aligners) consume these as fragment splits.
var pauseMarkers = map[rune]string{
	0x0E2F: ",", // ฯ
	0x0E4F: ".", // ๏
	0x0E5A: ".", // ๚
	0x0E5B: ".", // ๛
}

// Transliterate returns the RTGS romanization of s. Non-Thai runes pass
// through unchanged.
func Transliterate(s string) string {
	if s == "" {
		return s
	}
	rs := []rune(s)

	silent := make([]bool, len(rs))

	// thanthakhat: silence the consonant immediately before it (and the
	// thanthakhat itself). When that silenced consonant is a sonorant
	// (ร or ล), it's typically the second member of a Sanskrit-derived
	// cluster — the preceding consonant is also silent (e.g. ทร์ in
	// จันทร์ → "chan" not "chanth"). Other consonants get a single-rune
	// silencing only (e.g. ต์ in เซ็นต์ silences just ต).
	for i, r := range rs {
		if r != thanthakhat || i == 0 {
			continue
		}
		silent[i] = true
		silent[i-1] = true
		if (rs[i-1] == 0x0E23 || rs[i-1] == 0x0E25) && i >= 2 {
			if _, isCons := consonants[rs[i-2]]; isCons {
				silent[i-2] = true
			}
		}
	}

	// Silent ห before a sonorant — heuristic: any ห immediately followed
	// by a sonorant consonant is treated as silent. PyThaiNLP refines
	// this with syllable analysis; the approximation catches most cases
	// at the cost of occasional false positives mid-cluster.
	for i, r := range rs {
		if silent[i] {
			continue
		}
		if r == hoHip && i+1 < len(rs) && sonorant[rs[i+1]] {
			silent[i] = true
		}
	}

	var b strings.Builder
	b.Grow(len(s) * 2)
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		if silent[i] || ignored[r] {
			continue
		}

		if v, ok := prepositioned[r]; ok {
			j := i + 1
			for j < len(rs) && silent[j] {
				j++
			}
			if j >= len(rs) {
				b.WriteString(v)
				continue
			}
			cons, hasC := consonants[rs[j]]
			if !hasC {
				b.WriteString(v)
				continue
			}
			b.WriteString(cons)
			i = applyTrailingPattern(rs, &b, r, j, v)
			continue
		}

		if v, ok := consonants[r]; ok {
			b.WriteString(v)
			continue
		}
		if v, ok := vowelSigns[r]; ok {
			b.WriteString(v)
			continue
		}
		if v, ok := pauseMarkers[r]; ok {
			b.WriteString(v)
			continue
		}
		if r >= 0x0E50 && r <= 0x0E59 { // Thai digits
			b.WriteRune('0' + (r - 0x0E50))
			continue
		}
		switch r {
		case 0x0E24: // ฤ
			b.WriteString("rue")
			continue
		case 0x0E26: // ฦ
			b.WriteString("lue")
			continue
		case 0x0E45: // ๅ — vowel-length marker that pairs with ฤ/ฦ; absorbed.
			continue
		case 0x0E46: // ๆ — repetition mark; emit nothing.
			continue
		}
		if r >= BlockStart && r <= BlockEnd {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// applyTrailingPattern handles multi-rune vowel patterns that follow the
// consonant at j. Tone marks and other ignored runes between the
// consonant and the trailing vowel are skipped. Returns the new
// outer-loop index (the for-loop will then i++ past it).
func applyTrailingPattern(rs []rune, b *strings.Builder, pre rune, j int, fallback string) int {
	k := j + 1
	// Skip tone marks but NOT maitaikhu — maitaikhu is itself a pattern
	// trigger (เC็ → short e).
	for k < len(rs) && toneMarks[rs[k]] {
		k++
	}
	if pre == 0x0E40 { // เ
		// เCือ → Cuea (e.g. เมือง "mueang")
		if k+1 < len(rs) && rs[k] == 0x0E37 && rs[k+1] == 0x0E2D {
			b.WriteString("uea")
			return k + 1
		}
		// เCาะ → Co (short o)
		if k+1 < len(rs) && rs[k] == 0x0E32 && rs[k+1] == 0x0E30 {
			b.WriteString("o")
			return k + 1
		}
		if k < len(rs) {
			switch rs[k] {
			case 0x0E32:
				b.WriteString("ao") // เCา
				return k
			case 0x0E30:
				b.WriteString("e") // เCะ (short e)
				return k
			case maitaikhu:
				b.WriteString("e") // เC็ (short e via maitaikhu)
				return k
			}
		}
	}
	if pre == 0x0E41 && k < len(rs) && rs[k] == 0x0E30 {
		b.WriteString("ae") // แCะ
		return k
	}
	if pre == 0x0E42 && k < len(rs) && rs[k] == 0x0E30 {
		b.WriteString("o") // โCะ
		return k
	}
	b.WriteString(fallback)
	return j
}

// Contains reports whether s has any Thai-block rune.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}

// SplitPhrases breaks s into phrase-level fragments suitable for forced
// alignment (e.g. with Aeneas — https://github.com/readbeyond/aeneas).
// Thai script doesn't use spaces between words, so word-level segmentation
// would require a dictionary; but phrase boundaries are signalled by
// structural markers that ARE present in well-edited text:
//
//   - Unicode whitespace (Thai writers insert spaces at phrase breaks)
//   - Thai native pause marks: ฯ ๏ ๚ ๛
//   - Western sentence/clause terminators: . , ; : ! ?
//
// The boundary characters themselves are not included in any fragment.
// Empty fragments are dropped. Word segmentation is not attempted.
func SplitPhrases(s string) []string {
	var out []string
	var b strings.Builder
	flush := func() {
		if b.Len() == 0 {
			return
		}
		frag := strings.TrimSpace(b.String())
		if frag != "" {
			out = append(out, frag)
		}
		b.Reset()
	}
	for _, r := range s {
		if isPhraseBoundary(r) {
			flush()
			continue
		}
		b.WriteRune(r)
	}
	flush()
	return out
}

func isPhraseBoundary(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r', '\v', '\f', '\u00A0', '\u2028', '\u2029':
		return true
	case 0x0E2F, 0x0E4F, 0x0E5A, 0x0E5B: // ฯ ๏ ๚ ๛
		return true
	case '.', ',', ';', ':', '!', '?':
		return true
	}
	return false
}
