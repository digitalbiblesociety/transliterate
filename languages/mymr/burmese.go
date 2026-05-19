// Package mymr romanizes Myanmar (Burmese) script (U+1000-U+109F).
// Default is MLCTS2 (ASCII with ":" / "." / "'" tone marks);
// -mode ipa emits IPA, -mode mlcts emits orthographic MLCTS.
package mymr

import (
	"sort"
	"strings"
)

const (
	BlockStart rune = 0x1000
	BlockEnd   rune = 0x109F

	asat     rune = 0x103A
	virama   rune = 0x1039
	anusvara rune = 0x1036
	aukmyit  rune = 0x1037
	visarga  rune = 0x1038

	yaPin  rune = 0x103B
	yaYit  rune = 0x103C
	waHswe rune = 0x103D
	haTho  rune = 0x103E
)

// mode indexes the [3]string romanization tables: 0=MLCTS2, 1=IPA, 2=MLCTS.
type mode int

const (
	modeMLCTS2 mode = iota
	modeIPA
	modeMLCTS
)

// Transliterate returns the MLCTS2 romanization of s. Non-Burmese runes
// pass through unchanged.
func Transliterate(s string) string { return run(s, modeMLCTS2) }

// TransliterateIPA returns the IPA transcription of s.
func TransliterateIPA(s string) string { return run(s, modeIPA) }

// TransliterateMLCTS returns the orthographic MLCTS transliteration of s.
func TransliterateMLCTS(s string) string { return run(s, modeMLCTS) }

// Contains reports whether s has any Myanmar-block rune.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}

func run(s string, m mode) string {
	if s == "" {
		return s
	}
	rs := normalizeOrder([]rune(s))
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
		if r >= 0x1040 && r <= 0x1049 { // Myanmar digits
			b.WriteRune('0' + (r - 0x1040))
			i++
			continue
		}
		if p, ok := punctuation[r]; ok {
			b.WriteString(p)
			i++
			continue
		}

		if _, ok := consonants[r]; ok {
			out, next := emitSyllable(rs, i, m)
			b.WriteString(out)
			i = next
			continue
		}
		if v, ok := independentVowels[r]; ok {
			out, next := emitVowelSyllable(rs, i, m, v)
			b.WriteString(out)
			i = next
			continue
		}
		i++
	}
	return b.String()
}

// normalizeOrder swaps asat-aukmyit to aukmyit-asat (rhyme keys expect
// that order) and injects an implicit asat before a bare virama after
// a consonant so Pali stacks parse as a single closed syllable.
func normalizeOrder(rs []rune) []rune {
	for i := 0; i+1 < len(rs); i++ {
		if rs[i] == asat && rs[i+1] == aukmyit {
			rs[i], rs[i+1] = rs[i+1], rs[i]
		}
	}
	out := make([]rune, 0, len(rs)+4)
	for i, r := range rs {
		if r == virama && i > 0 {
			prev := rs[i-1]
			if _, isCons := consonants[prev]; isCons {
				out = append(out, asat)
			}
		}
		out = append(out, r)
	}
	return out
}

func emitSyllable(rs []rune, i int, m mode) (string, int) {
	onsetR := rs[i]
	i++

	var medials []rune
	for i < len(rs) && isMedial(rs[i]) {
		medials = append(medials, rs[i])
		i++
	}

	rhymeStart := i
	for i < len(rs) && isRhymeBody(rs[i]) {
		i++
	}
	for {
		next := tryFinalGroup(rs, i)
		if next == i {
			break
		}
		i = next
	}
	rhyme := string(rs[rhymeStart:i])

	if i < len(rs) && rs[i] == virama {
		i++
	}

	return composeSyllable(onsetR, medials, rhyme, m), i
}

func emitVowelSyllable(rs []rune, i int, m mode, base [3]string) (string, int) {
	out := base[m]
	i++
	if i < len(rs) {
		switch rs[i] {
		case visarga:
			out = applyTone(out, ":", "́", m)
			i++
		case aukmyit:
			out = applyTone(out, ".", "̰", m)
			i++
		}
	}
	return out, i
}

func composeSyllable(onsetR rune, medials []rune, rhyme string, m mode) string {
	onset := consonants[onsetR][m]

	// Devoicing (◌ှ) attaches to the onset; j-glide and w-glide follow.
	// ASCII modes use a "DV" marker that devoiceShift later converts to
	// an h-prefix; IPA uses combining ring below.
	var devoice, hasYaPin, hasYaYit, hasWa bool
	for _, med := range medials {
		switch med {
		case haTho:
			devoice = true
		case yaPin:
			hasYaPin = true
		case yaYit:
			hasYaYit = true
		case waHswe:
			hasWa = true
		}
	}

	var sb strings.Builder
	sb.WriteString(onset)
	if devoice {
		if m == modeIPA {
			sb.WriteString("̥")
		} else {
			sb.WriteString("DV")
		}
	}
	if hasYaPin {
		writeYaMedial(&sb, true, m)
	}
	if hasYaYit {
		writeYaMedial(&sb, false, m)
	}
	if hasWa {
		sb.WriteString("w")
	}

	sb.WriteString(lookupRhyme(rhyme, onsetR, m))
	return soundChanges(sb.String(), m)
}

func writeYaMedial(sb *strings.Builder, isYaPin bool, m mode) {
	switch m {
	case modeIPA, modeMLCTS2:
		sb.WriteString("j")
	case modeMLCTS:
		if isYaPin {
			sb.WriteString("y")
		} else {
			sb.WriteString("r")
		}
	}
}

// lookupRhyme returns the rhyme romanization. An empty rhyme on a
// regular consonant emits the inherent /a/; on an independent vowel it
// emits nothing (the vowel already supplied itself). Unknown rhymes
// fall back rune-by-rune to handle Pali stack patterns not in the table.
func lookupRhyme(rhyme string, onsetR rune, m mode) string {
	if rhyme == "" {
		if _, ok := independentVowels[onsetR]; ok {
			return ""
		}
		return inherentVowel[m]
	}
	for _, entry := range rhymeTable {
		if entry.key == rhyme {
			return entry.out[m]
		}
	}
	var sb strings.Builder
	for _, r := range rhyme {
		if v, ok := vowelFallback[r]; ok {
			sb.WriteString(v[m])
			continue
		}
		if t, ok := toneFallback[r]; ok {
			sb.WriteString(t[m])
			continue
		}
		if c, ok := consonants[r]; ok {
			sb.WriteString(c[m])
		}
	}
	return sb.String()
}

// applyTone appends an ASCII tone marker or, in IPA mode, attaches a
// combining diacritic to the syllable's final vowel.
func applyTone(syl, ascii, ipa string, m mode) string {
	if m == modeIPA {
		return attachIPATone(syl, ipa)
	}
	return syl + ascii
}

func attachIPATone(syl, mark string) string {
	rs := []rune(syl)
	for i := len(rs) - 1; i >= 0; i-- {
		if isIPAVowel(rs[i]) {
			out := make([]rune, 0, len(rs)+1)
			out = append(out, rs[:i+1]...)
			out = append(out, []rune(mark)...)
			out = append(out, rs[i+1:]...)
			return string(out)
		}
	}
	return syl + mark
}

func isIPAVowel(r rune) bool {
	switch r {
	case 'a', 'e', 'i', 'o', 'u', 'ɛ', 'ɔ', 'ə', 'ɪ', 'ʊ',
		'à', 'á', 'è', 'é', 'ì', 'í', 'ò', 'ó', 'ù', 'ú':
		return true
	}
	return false
}

func isMedial(r rune) bool {
	switch r {
	case yaPin, yaYit, waHswe, haTho:
		return true
	}
	return false
}

func isRhymeBody(r rune) bool {
	switch r {
	case anusvara, aukmyit, visarga:
		return true
	}
	if _, ok := vowelFallback[r]; ok {
		return true
	}
	return false
}

// tryFinalGroup consumes one (final-consonant + anusvara? + asat/tone+)
// group starting at rs[i]. Returns i unchanged if nothing matches.
func tryFinalGroup(rs []rune, i int) int {
	if i >= len(rs) {
		return i
	}
	start := i
	if _, ok := consonants[rs[i]]; ok {
		i++
	} else if rs[i] != anusvara {
		return start
	}
	if i < len(rs) && rs[i] == anusvara {
		i++
	}
	consumed := false
	for i < len(rs) {
		switch rs[i] {
		case asat, aukmyit, visarga:
			consumed = true
			i++
			continue
		}
		break
	}
	if !consumed {
		return start
	}
	return i
}

// out indices: 0=MLCTS2, 1=IPA, 2=MLCTS.
type rhymeEntry struct {
	key string
	out [3]string
}

var rhymeTable = []rhymeEntry{
	{"ိုင်း", [3]string{"ain:", "áɪɴ", "uing:"}},
	{"ိုင့်", [3]string{"ain.", "a̰ɪɴ", "uing."}},
	{"ေါင်း", [3]string{"aun:", "áʊɴ", "aung:"}},
	{"ေါင့်", [3]string{"aun.", "a̰ʊɴ", "aung."}},
	{"ောင်း", [3]string{"aun:", "áʊɴ", "aung:"}},
	{"ောင့်", [3]string{"aun.", "a̰ʊɴ", "aung."}},
	{"ိုက်", [3]string{"ai'", "aɪʔ", "uik"}},
	{"ေါက်", [3]string{"au'", "aʊʔ", "auk"}},
	{"ောက်", [3]string{"au'", "aʊʔ", "auk"}},
	{"ိုင်", [3]string{"ain", "àɪɴ", "uing"}},
	{"ေါင်", [3]string{"aun", "àʊɴ", "aung"}},
	{"ောင်", [3]string{"aun", "àʊɴ", "aung"}},
	{"ုမ်း", [3]string{"oun:", "óʊɴ", "um:"}},
	{"ုမ့်", [3]string{"oun.", "o̰ʊɴ", "um."}},
	{"ုန်း", [3]string{"oun:", "óʊɴ", "un:"}},
	{"ုန့်", [3]string{"oun.", "o̰ʊɴ", "un."}},
	{"ိမ်း", [3]string{"ein:", "éɪɴ", "im:"}},
	{"ိမ့်", [3]string{"ein.", "ḛɪɴ", "im."}},
	{"ိန်း", [3]string{"ein:", "éɪɴ", "in:"}},
	{"ိန့်", [3]string{"ein.", "ḛɪɴ", "in."}},
	{"ိုး", [3]string{"ou:", "ó", "ui:"}},
	{"ို့", [3]string{"ou.", "o̰", "ui."}},
	{"ေါ့", [3]string{"o.", "ɔ̰", "au."}},
	{"ေါ်", [3]string{"o", "ɔ̀", "au"}},
	{"ော့", [3]string{"o.", "ɔ̰", "au."}},
	{"ော်", [3]string{"o", "ɔ̀", "au"}},

	{"ုံး", [3]string{"oun:", "óʊɴ", "um:"}},
	{"ုံ့", [3]string{"oun.", "o̰ʊɴ", "um."}},
	{"ုဏ်", [3]string{"oun", "òʊɴ", "un"}},
	{"ုမ်", [3]string{"oun", "òʊɴ", "um"}},
	{"ုပ်", [3]string{"ou", "oʊʔ", "up"}},
	{"ုန်", [3]string{"oun", "òʊɴ", "un"}},
	{"ုတ်", [3]string{"ou", "oʊʔ", "ut"}},
	{"ိံး", [3]string{"ein:", "éɪɴ", "im:"}},
	{"ိံ့", [3]string{"ein.", "ḛɪɴ", "im."}},
	{"ိမ်", [3]string{"ein", "èɪɴ", "im"}},
	{"ိပ်", [3]string{"ei'", "eɪʔ", "ip'"}},
	{"ိန်", [3]string{"ein", "èɪɴ", "in"}},
	{"ိတ်", [3]string{"ei'", "eɪʔ", "it"}},
	{"မ်း", [3]string{"an:", "áɴ", "am:"}},
	{"မ့်", [3]string{"an.", "a̰ɴ", "am."}},
	{"န်း", [3]string{"an:", "áɴ", "an:"}},
	{"န့်", [3]string{"an.", "a̰ɴ", "an."}},
	{"ဉ်း", [3]string{"in:", "ɪ́ɴ", "any:"}},
	{"ည်း", [3]string{"i:", "í", "any:"}},
	{"ဉ့်", [3]string{"in.", "ɪ̰ɴ", "any."}},
	{"ည့်", [3]string{"i.", "ḭ", "any."}},
	{"င်း", [3]string{"in:", "ɪ́ɴ", "ang:"}},
	{"င့်", [3]string{"in.", "ɪ̰ɴ", "ang."}},
	{"ား", [3]string{"a:", "á", "a:"}},
	{"ါး", [3]string{"a:", "á", "a:"}},
	{"ာ့", [3]string{"a.", "a̰", "a."}},
	{"ါ့", [3]string{"a.", "a̰", "a."}},
	{"ေး", [3]string{"ei:", "é", "e:"}},
	{"ေ့", [3]string{"ei.", "ḛ", "e."}},
	{"ူး", [3]string{"u:", "ú", "u:"}},
	{"ီး", [3]string{"i:", "í", "i:"}},
	{"ဲ့", [3]string{"e.", "ɛ̰", "ai."}},
	{"ို", [3]string{"ou", "ò", "ui"}},
	{"ေါ", [3]string{"o:", "ɔ́", "au:"}},
	{"ော", [3]string{"o:", "ɔ́", "au:"}},
	{"ုံ", [3]string{"oun", "òʊɴ", "um"}},
	{"ိံ", [3]string{"ein", "èɪɴ", "im"}},
	{"ံး", [3]string{"an:", "áɴ", "am:"}},
	{"ံ့", [3]string{"an.", "a̰ɴ", "am."}},

	{"ယ်", [3]string{"e", "ɛ̀", "ai"}},
	{"မ်", [3]string{"an", "àɴ", "am"}},
	{"ပ်", [3]string{"a'", "aʔ", "ap"}},
	{"န်", [3]string{"an", "àɴ", "an"}},
	{"တ်", [3]string{"a'", "aʔ", "at"}},
	{"ဉ်", [3]string{"in", "ɪ̀ɴ", "any"}},
	{"ည်", [3]string{"i", "ì", "any"}},
	{"စ်", [3]string{"i'", "ɪʔ", "ac"}},
	{"င်", [3]string{"in", "ɪ̀ɴ", "ang"}},
	{"ဲ", [3]string{"e:", "ɛ́", "ai:"}},
	{"ေ", [3]string{"ei", "è", "e"}},
	{"ူ", [3]string{"u", "ù", "u"}},
	{"ု", [3]string{"u.", "ṵ", "u."}},
	{"ီ", [3]string{"i", "ì", "i"}},
	{"ိ", [3]string{"i.", "ḭ", "i."}},
	{"ံ", [3]string{"an", "àɴ", "am"}},
	{"်", [3]string{"e'", "ɛʔ", "ak"}},
	{"ါ", [3]string{"a", "à", "a"}},
	{"ာ", [3]string{"a", "à", "a"}},
	{"း", [3]string{"a:", "á", "a:"}},
	{"့", [3]string{"a.", "a̰", "a."}},
}

// inherentVowel is the abugida default /a/ (creaky-toned).
var inherentVowel = [3]string{"a.", "a̰", "a."}

var vowelFallback = map[rune][3]string{
	0x102B: {"a", "à", "a"},
	0x102C: {"a", "à", "a"},
	0x102D: {"i", "ḭ", "i"},
	0x102E: {"i", "ì", "i"},
	0x102F: {"u", "ṵ", "u"},
	0x1030: {"u", "ù", "u"},
	0x1031: {"ei", "è", "e"},
	0x1032: {"e", "ɛ̀", "ai"},
}

var toneFallback = map[rune][3]string{
	anusvara: {"an", "àɴ", "am"},
	aukmyit:  {".", "̰", "."},
	visarga:  {":", "́", ":"},
	asat:     {"", "", ""},
}

func soundChanges(s string, m mode) string {
	switch m {
	case modeIPA:
		s = strings.ReplaceAll(s, "kʰj", "tɕʰ")
		s = strings.ReplaceAll(s, "kj", "tɕ")
		s = strings.ReplaceAll(s, "ɡj", "dʑ")
		s = strings.ReplaceAll(s, "ŋj", "ɲ")
		s = strings.ReplaceAll(s, "jw", "w")
		s = strings.ReplaceAll(s, "j̥", "ʃ")
		s = strings.ReplaceAll(s, "ŋ̥", "ŋ̊")
		s = strings.ReplaceAll(s, "w̥", "ʍ")
		s = strings.ReplaceAll(s, "l̥j", "ʃ")
		s = strings.ReplaceAll(s, "θ̥j", "ʃ")
	case modeMLCTS:
		s = devoiceShift(s)
	case modeMLCTS2:
		s = strings.ReplaceAll(s, "khj", "ch")
		s = devoiceShift(s)
		s = strings.ReplaceAll(s, "ngj", "nj")
		s = strings.ReplaceAll(s, "hj", "sh")
		s = strings.ReplaceAll(s, "wu", "u")
		s = strings.ReplaceAll(s, "ww", "w")
	}
	return s
}

// devoiceShift moves the "DV" marker to an h-prefix on the preceding
// consonant. "[C]DV..." → "h[C]..."; two-byte digraphs move as a unit.
func devoiceShift(s string) string {
	for {
		i := strings.Index(s, "DV")
		if i < 0 {
			return s
		}
		if i == 0 {
			s = "h" + s[2:]
			continue
		}
		j := i - 1
		if i >= 2 {
			prev2 := s[i-2 : i]
			switch prev2 {
			case "th", "ng", "kh", "hp", "ph", "ht", "hs":
				j = i - 2
			}
		}
		s = s[:j] + "h" + s[j:i] + s[i+2:]
	}
}

var consonants = map[rune][3]string{
	0x1000: {"k", "k", "k"},
	0x1001: {"kh", "kʰ", "hk"},
	0x1002: {"g", "ɡ", "g"},
	0x1003: {"g", "ɡ", "gh"},
	0x1004: {"ng", "ŋ", "ng"},
	0x1005: {"s", "s", "c"},
	0x1006: {"hs", "sʰ", "hc"},
	0x1007: {"z", "z", "j"},
	0x1008: {"z", "z", "jh"},
	0x1009: {"nj", "ɲ", "ny"},
	0x100A: {"nj", "ɲ", "ny"},
	0x100B: {"t", "t", "t"},
	0x100C: {"ht", "tʰ", "ht"},
	0x100D: {"d", "d", "d"},
	0x100E: {"d", "d", "dh"},
	0x100F: {"n", "n", "n"},
	0x1010: {"t", "t", "t"},
	0x1011: {"ht", "tʰ", "ht"},
	0x1012: {"d", "d", "d"},
	0x1013: {"d", "d", "dh"},
	0x1014: {"n", "n", "n"},
	0x1015: {"p", "p", "p"},
	0x1016: {"hp", "pʰ", "hp"},
	0x1017: {"b", "b", "b"},
	0x1018: {"b", "b", "bh"},
	0x1019: {"m", "m", "m"},
	0x101A: {"j", "j", "y"},
	0x101B: {"j", "j", "r"},
	0x101C: {"l", "l", "l"},
	0x101D: {"w", "w", "w"},
	0x101E: {"th", "θ", "s"},
	0x101F: {"h", "h", "h"},
	0x1020: {"l", "l", "l"},
	0x1021: {"", "ʔ", ""},
	0x103F: {"th", "θ", "s"},
}

var independentVowels = map[rune][3]string{
	0x1023: {"i.", "ʔḭ", "i."},
	0x1024: {"i", "ʔì", "i"},
	0x1025: {"u.", "ʔṵ", "u."},
	0x1026: {"u", "ʔù", "u"},
	0x1027: {"ei", "ʔè", "e"},
	0x1029: {"o:", "ʔɔ́", "au:"},
	0x102A: {"o", "ʔɔ̀", "au"},
}

var punctuation = map[rune]string{
	0x104A: ",",
	0x104B: ".",
}

// Sort longest-first so linear scans in lookupRhyme are greedy.
func init() {
	sort.SliceStable(rhymeTable, func(i, j int) bool {
		ri := []rune(rhymeTable[i].key)
		rj := []rune(rhymeTable[j].key)
		if len(ri) != len(rj) {
			return len(ri) > len(rj)
		}
		return rhymeTable[i].key < rhymeTable[j].key
	})
}
