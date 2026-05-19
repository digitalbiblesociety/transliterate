// Package laoo romanizes Lao script (U+0E80-U+0EFF) using the
// BGN/PCGN 1966 system. ASCII with macrons for long vowels; tone
// marks are dropped.
package laoo

import "strings"

const (
	BlockStart rune = 0x0E80
	BlockEnd   rune = 0x0EFF

	hoSung    rune = 0x0EAB
	cancel    rune = 0x0ECC
	niggahita rune = 0x0ECD
	maiKan    rune = 0x0EB1
	maiKong   rune = 0x0EBB
	saraAa    rune = 0x0EB2
)

// Transliterate returns the BGN/PCGN romanization of s. Non-Lao runes
// pass through unchanged.
func Transliterate(s string) string {
	if s == "" {
		return s
	}
	rs := []rune(s)

	silent := make([]bool, len(rs))
	for i, r := range rs {
		if r == hoSung && i+1 < len(rs) && sonorantAfterHo[rs[i+1]] {
			silent[i] = true
		}
	}
	for i := 0; i+1 < len(rs); i++ {
		if rs[i+1] == cancel {
			if _, ok := consonants[rs[i]]; ok {
				silent[i] = true
			}
		}
	}

	var b strings.Builder
	b.Grow(len(s) * 2)

	i := 0
	for i < len(rs) {
		r := rs[i]

		if silent[i] {
			i++
			continue
		}

		if r < BlockStart || r > BlockEnd {
			b.WriteRune(r)
			i++
			continue
		}

		if r >= 0x0ED0 && r <= 0x0ED9 {
			b.WriteRune('0' + (r - 0x0ED0))
			i++
			continue
		}
		if p, ok := punctuation[r]; ok {
			b.WriteString(p)
			i++
			continue
		}

		if pv, ok := prepositioned[r]; ok {
			j := i + 1
			for j < len(rs) && silent[j] {
				j++
			}
			if j >= len(rs) {
				b.WriteString(pv)
				i++
				continue
			}
			cons, hasC := lookupConsonant(rs, j, silent)
			if !hasC {
				b.WriteString(pv)
				i++
				continue
			}
			b.WriteString(cons.text)
			i = applyVowelPrefix(rs, &b, r, cons.next)
			continue
		}

		if _, ok := toneMarks[r]; ok {
			i++
			continue
		}
		if r == cancel {
			i++
			continue
		}

		if c, ok := consonants[r]; ok {
			if v := medialVowel(rs, i, silent); v != "" {
				b.WriteString(v)
				i++
				continue
			}
			if isFinalPosition(rs, i, silent) {
				if c.final != "" {
					b.WriteString(c.final)
				} else {
					b.WriteString(c.initial)
				}
			} else {
				b.WriteString(c.initial)
			}
			i++
			continue
		}

		if r == niggahita {
			if i+1 < len(rs) && rs[i+1] == saraAa {
				b.WriteString("am")
				i += 2
				continue
			}
			b.WriteString("ǭ")
			i++
			continue
		}

		if v, ok := vowelSigns[r]; ok {
			next, out, ok := vowelMultiRune(rs, i)
			if ok {
				b.WriteString(out)
				i = next
				continue
			}
			b.WriteString(v)
			i++
			continue
		}

		if r == 0x0EC6 {
			b.WriteString("-")
			i++
			continue
		}

		i++
	}
	return b.String()
}

// Contains reports whether s has any Lao-block rune.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}

type consonantLookup struct {
	text string
	next int
}

// lookupConsonant resolves the onset consonant at rs[j], transparently
// skipping a silent ho-sung in a ◌ຫ-sonorant compound.
func lookupConsonant(rs []rune, j int, silent []bool) (consonantLookup, bool) {
	if j >= len(rs) {
		return consonantLookup{}, false
	}
	if silent[j] && j+1 < len(rs) {
		if c, ok := consonants[rs[j+1]]; ok {
			return consonantLookup{text: c.initial, next: j + 2}, true
		}
		return consonantLookup{}, false
	}
	if c, ok := consonants[rs[j]]; ok {
		return consonantLookup{text: c.initial, next: j + 1}, true
	}
	return consonantLookup{}, false
}

// applyVowelPrefix emits the trailing vowel pattern paired with a
// pre-positioned vowel. The consonant has already been emitted; j is
// the index past it. Returns the new outer-loop index.
func applyVowelPrefix(rs []rune, b *strings.Builder, pre rune, j int) int {
	k := j
	for k < len(rs) && toneMarks[rs[k]] {
		k++
	}

	switch pre {
	case 0x0EC0: // ເ
		if k < len(rs) && rs[k] == maiKong {
			m := k + 1
			for m < len(rs) && toneMarks[rs[m]] {
				m++
			}
			if m < len(rs) && rs[m] == saraAa {
				b.WriteString("ao")
				return m + 1
			}
		}
		if k+1 < len(rs) && rs[k] == saraAa && rs[k+1] == 0x0EB0 {
			b.WriteString("o")
			return k + 2
		}
		if k < len(rs) && rs[k] == 0x0EB4 {
			b.WriteString("œ")
			return k + 1
		}
		if k < len(rs) && rs[k] == 0x0EB5 {
			b.WriteString("œ̄")
			return k + 1
		}
		if k+1 < len(rs) && rs[k] == 0x0EB6 && rs[k+1] == 0x0EAD {
			b.WriteString("ua")
			return k + 2
		}
		if k+1 < len(rs) && rs[k] == 0x0EB7 && rs[k+1] == 0x0EAD {
			b.WriteString("ūa")
			return k + 2
		}
		if k+1 < len(rs) && rs[k] == maiKan && rs[k+1] == 0x0EBD {
			b.WriteString("ia")
			return k + 2
		}
		if k < len(rs) && rs[k] == 0x0EBD {
			b.WriteString("īa")
			return k + 1
		}
		if k < len(rs) && rs[k] == 0x0EB0 {
			b.WriteString("e")
			return k + 1
		}
		if k < len(rs) && rs[k] == maiKan {
			b.WriteString("e")
			return k + 1
		}
		b.WriteString("ē")
		return j
	case 0x0EC1: // ແ
		if k < len(rs) && (rs[k] == 0x0EB0 || rs[k] == maiKan) {
			b.WriteString("ǣ")
			return k + 1
		}
		b.WriteString("ǣ")
		return j
	case 0x0EC2: // ໂ
		if k < len(rs) && rs[k] == 0x0EB0 {
			b.WriteString("o")
			return k + 1
		}
		b.WriteString("ō")
		return j
	case 0x0EC3, 0x0EC4: // ໃ ໄ
		b.WriteString("ai")
		return j
	}
	return j
}

func vowelMultiRune(rs []rune, i int) (int, string, bool) {
	r := rs[i]
	switch r {
	case maiKong:
		if i+2 < len(rs) && rs[i+1] == 0x0EA7 && rs[i+2] == 0x0EB0 {
			return i + 3, "ua", true
		}
		if i+1 < len(rs) && rs[i+1] == 0x0EA7 {
			return i + 2, "ua", true
		}
		return i + 1, "o", true
	case maiKan:
		if i+1 < len(rs) {
			switch rs[i+1] {
			case 0x0EBD:
				return i + 2, "ia", true
			case 0x0E8D:
				return i + 2, "ai", true
			}
		}
		return i + 1, "a", true
	case 0x0EB3:
		return i + 1, "am", true
	}
	return 0, "", false
}

// medialVowel returns a vowel-medial emission when ◌ວ or ◌ອ sits
// between two consonants with no intervening vowel sign: ◌ວ → "ua",
// ◌ອ → "ǭ". Subscript la (◌ຼ) and silent ho are transparent.
func medialVowel(rs []rune, i int, silent []bool) string {
	if !bracketedByConsonants(rs, i, silent) {
		return ""
	}
	switch rs[i] {
	case 0x0EA7:
		return "ua"
	case 0x0EAD:
		return "ǭ"
	}
	return ""
}

func bracketedByConsonants(rs []rune, i int, silent []bool) bool {
	prevCons := false
	for k := i - 1; k >= 0; k-- {
		if silent[k] {
			continue
		}
		r := rs[k]
		if r < BlockStart || r > BlockEnd {
			return false
		}
		if r == 0x0EBC { // ◌ຼ subscript la — counts as a consonant
			prevCons = true
			break
		}
		if _, ok := consonants[r]; ok {
			prevCons = true
			break
		}
		return false
	}
	if !prevCons {
		return false
	}
	for k := i + 1; k < len(rs); k++ {
		if silent[k] {
			continue
		}
		r := rs[k]
		if _, ok := toneMarks[r]; ok {
			continue
		}
		if _, ok := consonants[r]; ok {
			return true
		}
		return false
	}
	return false
}

// isFinalPosition is true when rs[i] (a consonant) acts as a coda.
// Lookback finds a preceding vowel in the same syllable; lookahead
// confirms the next non-tone-mark rune isn't a vowel/medial that would
// make rs[i] the onset of a new syllable.
func isFinalPosition(rs []rune, i int, silent []bool) bool {
	hadVowel := false
	for k := i - 1; k >= 0; k-- {
		if silent[k] {
			continue
		}
		r := rs[k]
		if r < BlockStart || r > BlockEnd {
			break
		}
		if _, ok := vowelSigns[r]; ok {
			hadVowel = true
			break
		}
		if _, ok := prepositioned[r]; ok {
			hadVowel = true
			break
		}
		if r == niggahita || r == 0x0EB3 {
			hadVowel = true
			break
		}
		if _, ok := consonants[r]; ok {
			// Pre-vowels live before their consonant in storage order;
			// peek one further back so CC patterns like ໂ◌C◌C see the
			// vowel.
			if k-1 >= 0 {
				if _, ok := prepositioned[rs[k-1]]; ok {
					hadVowel = true
				}
			}
			break
		}
	}
	if !hadVowel {
		return false
	}
	for k := i + 1; k < len(rs); k++ {
		if silent[k] {
			continue
		}
		r := rs[k]
		if _, ok := toneMarks[r]; ok {
			continue
		}
		if _, ok := vowelSigns[r]; ok {
			return false
		}
		if _, ok := prepositioned[r]; ok {
			return true
		}
		if r == niggahita || r == 0x0EB3 {
			return false
		}
		if r == 0x0EBC {
			return false
		}
		// ◌ວ/◌ອ followed by another consonant is a vowel-medial; the
		// current consonant is the onset, not a coda.
		if r == 0x0EA7 || r == 0x0EAD {
			m := k + 1
			for m < len(rs) && toneMarks[rs[m]] {
				m++
			}
			if m < len(rs) {
				if _, ok := consonants[rs[m]]; ok {
					return false
				}
			}
			return true
		}
		return true
	}
	return true
}

// consonantEntry pairs the initial and final-position romanizations.
// final == "" means BGN/PCGN drops this consonant as a coda.
type consonantEntry struct {
	initial string
	final   string
}

var consonants = map[rune]consonantEntry{
	0x0E81: {"k", "k"},
	0x0E82: {"kh", ""},
	0x0E84: {"kh", ""},
	0x0E87: {"ng", "ng"},
	0x0E88: {"ch", ""},
	0x0E8A: {"s", ""},
	0x0E8D: {"ny", "y"},
	0x0E94: {"d", "t"},
	0x0E95: {"t", ""},
	0x0E96: {"th", ""},
	0x0E97: {"th", ""},
	0x0E99: {"n", "n"},
	0x0E9A: {"b", "p"},
	0x0E9B: {"p", ""},
	0x0E9C: {"ph", ""},
	0x0E9D: {"f", ""},
	0x0E9E: {"ph", ""},
	0x0E9F: {"f", ""},
	0x0EA1: {"m", "m"},
	0x0EA2: {"y", ""},
	0x0EA3: {"r", "n"},
	0x0EA5: {"l", ""},
	0x0EA7: {"v", "o"},
	0x0EAA: {"s", ""},
	0x0EAB: {"h", ""},
	0x0EAD: {"", ""},
	0x0EAE: {"h", ""},
	0x0EDC: {"n", ""},
	0x0EDD: {"m", ""},
	0x0EBC: {"l", ""},
}

// sonorantAfterHo: sonorants that ◌ຫ silently precedes to form a
// high-class compound. ໜ and ໝ are pre-composed and not listed.
var sonorantAfterHo = map[rune]bool{
	0x0E87: true,
	0x0E8D: true,
	0x0E99: true,
	0x0EA1: true,
	0x0EA5: true,
	0x0EBC: true,
	0x0EA3: true,
	0x0EA7: true,
}

var prepositioned = map[rune]string{
	0x0EC0: "ē",
	0x0EC1: "ǣ",
	0x0EC2: "ō",
	0x0EC3: "ai",
	0x0EC4: "ai",
}

var vowelSigns = map[rune]string{
	0x0EB0: "a",
	0x0EB1: "a",
	0x0EB2: "ā",
	0x0EB3: "am",
	0x0EB4: "i",
	0x0EB5: "ī",
	0x0EB6: "ư",
	0x0EB7: "ư̄",
	0x0EB8: "u",
	0x0EB9: "ū",
	0x0EBB: "o",
	0x0EBD: "ia",
}

var toneMarks = map[rune]bool{
	0x0EC8: true,
	0x0EC9: true,
	0x0ECA: true,
	0x0ECB: true,
}

var punctuation = map[rune]string{
	0x0EAF: "...",
}
