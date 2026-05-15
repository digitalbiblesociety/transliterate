// Package tibt provides Wylie transliteration of Tibetan script
// (U+0F00..U+0FFF) to ASCII, plus TransliteratePhonetic for an
// approximate THL Simplified rendering.
//
// Each tsek-delimited syllable carries exactly one implicit /a/,
// placed on the root letter. The root is identified by the stack
// carrying the vowel sign or, in vowel-less syllables, by stripping
// orthographic suffixes off the right.
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
	return cleanPunctuation(transliterateWylie(s))
}

func transliterateWylie(s string) string {
	rs := []rune(s)
	var b strings.Builder
	b.Grow(len(s) * 2)

	i := 0
	for i < len(rs) {
		r := rs[i]
		if r < BlockStart || r > BlockEnd {
			b.WriteRune(r)
			i++
			continue
		}
		if r >= 0x0F20 && r <= 0x0F29 {
			b.WriteRune('0' + (r - 0x0F20))
			i++
			continue
		}
		if r == 0x0F0B || r == 0x0F0C {
			b.WriteByte(' ')
			i++
			continue
		}
		if p, ok := punctuation[r]; ok {
			b.WriteString(p)
			i++
			continue
		}
		if t, ok := trailingSign[r]; ok {
			b.WriteString(t)
			i++
			continue
		}
		if _, ok := headConsonant[r]; ok {
			i = emitSyllable(rs, i, &b)
			continue
		}
		i++
	}
	return b.String()
}

// cleanPunctuation promotes "།།" (two adjacent shads, the source's
// sentence-end pattern) to a period, then collapses any other doubled
// punctuation and trims spaces before "," / ".". Project rule:
// periods are for hard pauses only.
func cleanPunctuation(s string) string {
	if s == "" {
		return s
	}
	for {
		prev := s
		s = strings.ReplaceAll(s, " ,", ",")
		s = strings.ReplaceAll(s, " .", ".")
		s = strings.ReplaceAll(s, ",,", ".")
		s = strings.ReplaceAll(s, "..", ".")
		s = strings.ReplaceAll(s, "  ", " ")
		if s == prev {
			return s
		}
	}
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

// SplitPhrases breaks s at shad/whitespace/Latin-punctuation boundaries
// for forced alignment. Tsek is too granular to be a boundary.
func SplitPhrases(s string) []string {
	var out []string
	var b strings.Builder
	flush := func() {
		if b.Len() == 0 {
			return
		}
		if frag := strings.TrimSpace(b.String()); frag != "" {
			out = append(out, frag)
		}
		b.Reset()
	}
	for _, r := range s {
		if isTibetanPhraseBoundary(r) {
			flush()
			continue
		}
		b.WriteRune(r)
	}
	flush()
	return out
}

func isTibetanPhraseBoundary(r rune) bool {
	switch r {
	case 0x0F0D, 0x0F0E, 0x0F0F, 0x0F11, 0x0F14: // shad and variants
		return true
	case ' ', '\t', '\n', '\r', '\v', '\f', '\u00A0', '\u2028', '\u2029':
		return true
	case '.', ',', ';', ':', '!', '?':
		return true
	}
	return false
}

// stack is a head consonant with its subjoined letters and optional
// vowel sign. singleC (no subjoined letters) is what lets a stack act
// as a prefix or suffix.
type stack struct {
	head    rune
	sub     []rune
	letters string // head + subjoined concatenated as Wylie letters
	vowel   string // Wylie value of the vowel sign
	vowelR  rune   // raw vowel-sign rune; F71+F72/F74 canonicalised to F73/F75
	singleC bool
}

func (s stack) eligibleSuffix() bool    { return s.singleC && suffixSet[s.head] }
func (s stack) eligible2ndSuffix() bool { return s.singleC && secondSuffixSet[s.head] }

func parseSyllable(rs []rune, i int) ([]stack, []rune, int) {
	var stacks []stack
	var trailing []rune

	for i < len(rs) {
		r := rs[i]

		if h, ok := headConsonant[r]; ok {
			cur := stack{head: r, letters: h, singleC: true}
			i++
			for i < len(rs) {
				if sub, ok := subjoined[rs[i]]; ok {
					cur.letters += sub
					cur.sub = append(cur.sub, rs[i])
					cur.singleC = false
					i++
					continue
				}
				break
			}
			// F71 may combine with a following F72/F74 to form the
			// decomposed long-i / long-u; canonicalise to F73/F75.
			if i < len(rs) {
				switch rs[i] {
				case 0x0F71:
					cur.vowelR = 0x0F71
					if i+1 < len(rs) && rs[i+1] == 0x0F72 {
						cur.vowel = "I"
						cur.vowelR = 0x0F73
						i += 2
					} else if i+1 < len(rs) && rs[i+1] == 0x0F74 {
						cur.vowel = "U"
						cur.vowelR = 0x0F75
						i += 2
					} else {
						cur.vowel = "A"
						i++
					}
				default:
					if v, ok := vowelSign[rs[i]]; ok {
						cur.vowel = v
						cur.vowelR = rs[i]
						i++
					}
				}
			}
			stacks = append(stacks, cur)
			continue
		}
		if _, ok := trailingSign[r]; ok {
			trailing = append(trailing, r)
			i++
			continue
		}
		if r == Virama {
			i++
			continue
		}
		break
	}
	return stacks, trailing, i
}

func emitSyllable(rs []rune, i int, b *strings.Builder) int {
	stacks, trailing, next := parseSyllable(rs, i)
	if len(stacks) == 0 {
		for _, r := range trailing {
			b.WriteString(trailingSign[r])
		}
		return next
	}

	rootIdx := findRoot(stacks)

	for idx, s := range stacks {
		// Stack-disambiguator only matters at the prefix→root boundary;
		// past the root the explicit /a/ already keeps "bar" from
		// collapsing to "bra". We emit "+" not EWTS's "." so the
		// downstream TTS doesn't read a sentence break.
		if idx > 0 && idx == rootIdx && needsWylieDot(stacks[idx-1], s) {
			b.WriteByte('+')
		}
		b.WriteString(s.letters)
		if s.vowel != "" {
			b.WriteString(s.vowel)
		} else if idx == rootIdx {
			b.WriteByte('a')
		}
	}
	for _, r := range trailing {
		b.WriteString(trailingSign[r])
	}
	return next
}

func needsWylieDot(prev, curr stack) bool {
	if !prev.singleC {
		return false
	}
	return wylieAmbiguous[[2]rune{prev.head, curr.head}]
}

// findRoot returns the root-stack index. Achung-with-vowel in
// non-initial position is treated as the genitive/instrumental
// particle ('i, 'u, ...) — the preceding stack becomes the root so
// པའི renders as "pa'i", not "p'i".
func findRoot(stacks []stack) int {
	for i, s := range stacks {
		if s.vowel != "" {
			if i > 0 && s.head == 0x0F60 {
				return i - 1
			}
			return i
		}
	}
	n := len(stacks)
	if n <= 1 {
		return 0
	}
	i := n - 1
	if stacks[i].eligibleSuffix() {
		if stacks[i].eligible2ndSuffix() && i > 1 && stacks[i-1].eligibleSuffix() {
			return i - 2
		}
		return i - 1
	}
	return i
}

// 0x0F68 (achen) is empty so the vowel-sign / inherent-/a/ logic
// fills it in. 0x0F60 (achung) keeps "'" so it survives as the
// genitive particle marker.
var headConsonant = map[rune]string{
	0x0F40: "k", 0x0F41: "kh", 0x0F42: "g", 0x0F43: "gh",
	0x0F44: "ng", 0x0F45: "c", 0x0F46: "ch", 0x0F47: "j",
	0x0F49: "ny", 0x0F4A: "T", 0x0F4B: "Th", 0x0F4C: "D",
	0x0F4D: "Dh", 0x0F4E: "N", 0x0F4F: "t", 0x0F50: "th",
	0x0F51: "d", 0x0F52: "dh", 0x0F53: "n", 0x0F54: "p",
	0x0F55: "ph", 0x0F56: "b", 0x0F57: "bh", 0x0F58: "m",
	0x0F59: "ts", 0x0F5A: "tsh", 0x0F5B: "dz", 0x0F5C: "dzh",
	0x0F5D: "w", 0x0F5E: "zh", 0x0F5F: "z", 0x0F60: "'",
	0x0F61: "y", 0x0F62: "r", 0x0F63: "l", 0x0F64: "sh",
	0x0F65: "Sh", 0x0F66: "s", 0x0F67: "h", 0x0F68: "",
	0x0F69: "kSh",
}

// Subjoined consonants mirror head consonants at +0x50.
var subjoined = map[rune]string{}

func init() {
	for r, v := range headConsonant {
		subjoined[r+0x50] = v
	}
}

var vowelSign = map[rune]string{
	0x0F71: "A",  // long a
	0x0F72: "i",
	0x0F73: "I",  // precomposed long i
	0x0F74: "u",
	0x0F75: "U",  // precomposed long u
	0x0F76: "Ri", // vocalic R (Sanskrit)
	0x0F77: "RI", // vocalic Rr (Sanskrit)
	0x0F78: "Li", // vocalic L (Sanskrit)
	0x0F79: "LI", // vocalic Ll (Sanskrit)
	0x0F7A: "e",
	0x0F7B: "ai",
	0x0F7C: "o",
	0x0F7D: "au",
	0x0F80: "-i", // reversed i (Sanskrit)
	0x0F81: "-I", // reversed long i
}

var trailingSign = map[rune]string{
	0x0F7E: "M", // anusvara
	0x0F7F: "H", // visarga
	0x0F82: "M", // candrabindu-like
	0x0F83: "M", // sna ldan
}

// English-style: periods only for hard sentence-ending pauses.
var punctuation = map[rune]string{
	0x0F0D: ",", // shad
	0x0F0E: ".", // double shad
	0x0F0F: ".", // rin chen spungs shad
	0x0F11: ".", // alt shad
	0x0F14: ",", // tsa-phru
}

// rjes-'jug. Prefix set isn't enumerated — right-to-left stripping
// only needs to recognise suffixes; whatever's left of the root is
// the prefix.
var suffixSet = map[rune]bool{
	0x0F42: true, // g
	0x0F44: true, // ng
	0x0F51: true, // d
	0x0F53: true, // n
	0x0F56: true, // b
	0x0F58: true, // m
	0x0F60: true, // '
	0x0F62: true, // r
	0x0F63: true, // l
	0x0F66: true, // s
}

// yang-'jug. "d" is Old Tibetan only.
var secondSuffixSet = map[rune]bool{
	0x0F66: true, // s
	0x0F51: true, // d
}

// (prefix-letter, subjoinable-letter) pairs where un-dotted output
// would collide with a single-stack subjoined cluster.
var wylieAmbiguous = map[[2]rune]bool{
	{0x0F42, 0x0F61}: true, // g.y vs gy
	{0x0F42, 0x0F62}: true, // g.r vs gr
	{0x0F42, 0x0F63}: true, // g.l vs gl
	{0x0F42, 0x0F5D}: true, // g.w vs gw
	{0x0F51, 0x0F62}: true, // d.r vs dr
	{0x0F51, 0x0F5D}: true, // d.w vs dw
	{0x0F56, 0x0F61}: true, // b.y vs by
	{0x0F56, 0x0F62}: true, // b.r vs br
	{0x0F56, 0x0F63}: true, // b.l vs bl
	{0x0F58, 0x0F61}: true, // m.y vs my
	{0x0F58, 0x0F62}: true, // m.r vs mr
}
