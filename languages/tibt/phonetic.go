package tibt

import "strings"

// TransliteratePhonetic returns an approximate THL Simplified Phonetic
// rendering of s — Lhasa Central Tibetan, no tone or length marking.
// Silent prefixes/superscripts drop, subjoined clusters phonologise,
// suffix d/n/l/s fronts the vowel (a→e, o→ö, u→ü), final g/b devoice.
// For reversible transcription use Transliterate.
func TransliteratePhonetic(s string) string {
	if s == "" {
		return s
	}
	return cleanPunctuation(transliteratePhonetic(s))
}

func transliteratePhonetic(s string) string {
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
		if t, ok := phoneticTrailingMap[r]; ok {
			b.WriteString(t)
			i++
			continue
		}
		if _, ok := headConsonant[r]; ok {
			i = emitSyllablePhonetic(rs, i, &b)
			continue
		}
		i++
	}
	return b.String()
}

func emitSyllablePhonetic(rs []rune, i int, b *strings.Builder) int {
	stacks, trailing, next := parseSyllable(rs, i)
	if len(stacks) == 0 {
		for _, r := range trailing {
			b.WriteString(phoneticTrailingMap[r])
		}
		return next
	}

	rootIdx := findRoot(stacks)

	// Pre-root stacks: a real Tibetan prefix drops; anything else (a
	// Sanskrit-style leading consonant like the "p" in "padme") emits
	// with its own /a/.
	for idx := range rootIdx {
		s := stacks[idx]
		if s.singleC && phoneticPrefixSilent[s.head] {
			continue
		}
		b.WriteString(initialFor(s))
		b.WriteByte('a')
	}

	root := stacks[rootIdx]
	initial := initialFor(root)
	vowel := vowelFor(root)

	// Suffix in slot rootIdx+1; rootIdx+2 (yang-'jug) is silent.
	// Achung-particle suffix carries its own vowel: "pa'i" → "pai".
	finalStr := ""
	if rootIdx+1 < len(stacks) {
		sx := stacks[rootIdx+1]
		if sx.singleC {
			if sx.vowel != "" {
				finalStr = phoneticInitialMap[sx.head] + vowelFor(sx)
			} else {
				if phoneticUmlautSuffix[sx.head] {
					vowel = umlautFor(vowel)
				}
				finalStr = phoneticFinalMap[sx.head]
			}
		}
	}

	b.WriteString(initial)
	b.WriteString(vowel)
	b.WriteString(finalStr)

	for _, r := range trailing {
		b.WriteString(phoneticTrailingMap[r])
	}
	return next
}

// initialFor returns the phonetic onset of a stack. r/s/l above a
// "real" consonant is a silent superscript; cluster lookup goes by the
// last subjoined modifier, falling back to head+subjoined concatenated.
func initialFor(s stack) string {
	head := s.head
	sub := s.sub
	if len(sub) > 0 && isSuperscript(head, sub[0]) {
		head = sub[0] - 0x50
		sub = sub[1:]
	}
	base := phoneticInitialMap[head]
	if len(sub) == 0 {
		return base
	}
	last := sub[len(sub)-1]
	if v, ok := phoneticCluster[[2]rune{head, last}]; ok {
		return v
	}
	var out strings.Builder
	out.WriteString(base)
	for _, sR := range sub {
		out.WriteString(phoneticInitialMap[sR-0x50])
	}
	return out.String()
}

func vowelFor(s stack) string {
	if s.vowelR == 0 {
		return "a"
	}
	if v, ok := phoneticVowelMap[s.vowelR]; ok {
		return v
	}
	return strings.ToLower(s.vowel)
}

func umlautFor(vowel string) string {
	switch vowel {
	case "a":
		return "e"
	case "o":
		return "ö"
	case "u":
		return "ü"
	}
	return vowel
}

// Only r/s/l can be superscripts, and only over a "real" consonant —
// when the first subjoined is a y/r/l/w/h modifier, headR is the root
// (e.g. ལྷ is l + subjoined h, not l-superscript + h).
func isSuperscript(headR, firstSub rune) bool {
	if headR != 0x0F62 && headR != 0x0F66 && headR != 0x0F63 {
		return false
	}
	switch firstSub {
	case 0x0FB1, 0x0FB2, 0x0FB3, 0x0FAD, 0x0FB7:
		return false
	}
	return true
}

// Voicing kept for reader clarity even where Lhasa devoices in
// practice. Achen (F68) and achung (F60) are vowel carriers → "".
var phoneticInitialMap = map[rune]string{
	0x0F40: "k", 0x0F41: "kh", 0x0F42: "g", 0x0F43: "g",
	0x0F44: "ng", 0x0F45: "c", 0x0F46: "ch", 0x0F47: "j",
	0x0F49: "ny", 0x0F4A: "t", 0x0F4B: "th", 0x0F4C: "d",
	0x0F4D: "d", 0x0F4E: "n", 0x0F4F: "t", 0x0F50: "th",
	0x0F51: "d", 0x0F52: "d", 0x0F53: "n", 0x0F54: "p",
	0x0F55: "ph", 0x0F56: "b", 0x0F57: "b", 0x0F58: "m",
	0x0F59: "ts", 0x0F5A: "tsh", 0x0F5B: "dz", 0x0F5C: "dz",
	0x0F5D: "w", 0x0F5E: "zh", 0x0F5F: "z", 0x0F60: "",
	0x0F61: "y", 0x0F62: "r", 0x0F63: "l", 0x0F64: "sh",
	0x0F65: "sh", 0x0F66: "s", 0x0F67: "h", 0x0F68: "",
	0x0F69: "ksh",
}

// (head, last-subjoined-modifier) → cluster onset. ya-btags
// palatalises labials, ra-btags is retroflex-leaning, la-btags
// simplifies, lh keeps as written. wa-zur falls through (silent w).
var phoneticCluster = map[[2]rune]string{
	{0x0F40, 0x0FB1}: "ky",
	{0x0F41, 0x0FB1}: "khy",
	{0x0F42, 0x0FB1}: "gy",
	{0x0F54, 0x0FB1}: "ch",
	{0x0F55, 0x0FB1}: "ch",
	{0x0F56, 0x0FB1}: "j",
	{0x0F58, 0x0FB1}: "ny",
	{0x0F67, 0x0FB1}: "hy",
	{0x0F40, 0x0FB2}: "tr",
	{0x0F41, 0x0FB2}: "thr",
	{0x0F42, 0x0FB2}: "dr",
	{0x0F4F, 0x0FB2}: "tr",
	{0x0F50, 0x0FB2}: "thr",
	{0x0F51, 0x0FB2}: "dr",
	{0x0F54, 0x0FB2}: "tr",
	{0x0F55, 0x0FB2}: "thr",
	{0x0F56, 0x0FB2}: "dr",
	{0x0F58, 0x0FB2}: "m",
	{0x0F64, 0x0FB2}: "sh",
	{0x0F66, 0x0FB2}: "s",
	{0x0F67, 0x0FB2}: "hr",
	{0x0F40, 0x0FB3}: "l",
	{0x0F42, 0x0FB3}: "l",
	{0x0F56, 0x0FB3}: "l",
	{0x0F66, 0x0FB3}: "l",
	{0x0F5F, 0x0FB3}: "d",
	{0x0F62, 0x0FB3}: "l",
	{0x0F63, 0x0FB7}: "lh",
}

// Length not marked: long ī/ū collapse to i/u.
var phoneticVowelMap = map[rune]string{
	0x0F71: "a",
	0x0F72: "i",
	0x0F73: "i",
	0x0F74: "u",
	0x0F75: "u",
	0x0F76: "i", // vocalic R (Sanskrit) → /i/ in Lhasa
	0x0F77: "i",
	0x0F78: "i",
	0x0F79: "i",
	0x0F7A: "e",
	0x0F7B: "e",
	0x0F7C: "o",
	0x0F7D: "o",
	0x0F80: "i",
	0x0F81: "i",
}

// d/'/s leave nothing — the umlaut has already shifted the vowel.
var phoneticFinalMap = map[rune]string{
	0x0F42: "k",
	0x0F44: "ng",
	0x0F51: "",
	0x0F53: "n",
	0x0F56: "p",
	0x0F58: "m",
	0x0F60: "",
	0x0F62: "r",
	0x0F63: "l",
	0x0F66: "",
}

var phoneticUmlautSuffix = map[rune]bool{
	0x0F51: true, // d
	0x0F53: true, // n
	0x0F63: true, // l
	0x0F66: true, // s
}

var phoneticPrefixSilent = map[rune]bool{
	0x0F42: true, // g
	0x0F51: true, // d
	0x0F56: true, // b
	0x0F58: true, // m
	0x0F60: true, // '
}

// Lowercase here; the Wylie path uses uppercase M/H.
var phoneticTrailingMap = map[rune]string{
	0x0F7E: "m",
	0x0F7F: "h",
	0x0F82: "m",
	0x0F83: "m",
}
