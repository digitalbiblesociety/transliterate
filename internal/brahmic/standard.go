package brahmic

// Most Brahmic scripts (Devanagari, Bengali, Gurmukhi, Gujarati, Oriya,
// Telugu, Kannada, Malayalam) share a near-parallel block layout, with
// canonical positions for each phoneme at fixed offsets from the block
// start. Tamil and Sinhala diverge meaningfully and are described
// explicitly in their own files.
//
// The "standardOffsets" table records, for each Latin output, the offset
// within the block where the corresponding glyph lives in the "canonical"
// Brahmic layout. A nil value means the script lacks that glyph at the
// expected offset (e.g. Bengali has no short ಎ at +0x0E); per-script files
// override that entry.

// standardIndependentVowels lists vowel offsets in canonical order.
// Pairs: {blockOffset, latinForm}. A latinForm of "" means "absent at this
// offset in some scripts; skip it during base construction".
var standardIndependentVowels = []struct {
	offset rune
	latin  string
}{
	{0x05, "a"},   // ಅ-style
	{0x06, "ā"},   // ಆ
	{0x07, "i"},   // ಇ
	{0x08, "ī"},   // ಈ
	{0x09, "u"},   // ಉ
	{0x0A, "ū"},   // ಊ
	{0x0B, "r̥"},  // ಋ
	{0x0C, "l̥"},  // ಌ
	{0x0E, "e"},   // ಎ (short e — present in Kannada/Telugu/Tamil/Malayalam; absent in Devanagari/Bengali/Gurmukhi/Gujarati/Oriya)
	{0x0F, "ē"},   // ಏ (long e — Devanagari etc. use ē here as the "regular" e)
	{0x10, "ai"},  // ಐ
	{0x12, "o"},   // ಒ (short o — same north/south split as short e)
	{0x13, "ō"},   // ಓ
	{0x14, "au"},  // ಔ
	{0x60, "r̥̄"}, // ೠ extra
	{0x61, "l̥̄"}, // ೡ extra
}

var standardConsonants = []struct {
	offset rune
	latin  string
}{
	{0x15, "k"}, {0x16, "kh"}, {0x17, "g"}, {0x18, "gh"}, {0x19, "ṅ"},
	{0x1A, "c"}, {0x1B, "ch"}, {0x1C, "j"}, {0x1D, "jh"}, {0x1E, "ñ"},
	{0x1F, "ṭ"}, {0x20, "ṭh"}, {0x21, "ḍ"}, {0x22, "ḍh"}, {0x23, "ṇ"},
	{0x24, "t"}, {0x25, "th"}, {0x26, "d"}, {0x27, "dh"}, {0x28, "n"},
	{0x2A, "p"}, {0x2B, "ph"}, {0x2C, "b"}, {0x2D, "bh"}, {0x2E, "m"},
	{0x2F, "y"}, {0x30, "r"}, {0x31, "ṟ"}, // ṟ archaic — present in Kannada/Tamil/Malayalam
	{0x32, "l"}, {0x33, "ḷ"}, {0x34, "ḻ"}, // ḻ archaic
	{0x35, "v"}, {0x36, "ś"}, {0x37, "ṣ"}, {0x38, "s"}, {0x39, "h"},
}

var standardVowelSigns = []struct {
	offset rune
	latin  string
}{
	{0x3E, "ā"}, {0x3F, "i"}, {0x40, "ī"}, {0x41, "u"}, {0x42, "ū"},
	{0x43, "r̥"}, {0x44, "l̥"},
	{0x46, "e"}, {0x47, "ē"}, {0x48, "ai"},
	{0x4A, "o"}, {0x4B, "ō"}, {0x4C, "au"},
}

// standardSpecial covers anusvara (+0x02), visarga (+0x03), avagraha
// (+0x3D). Optional candrabindu (+0x01) is rendered as ṁ̐ in ISO 15919 but
// scripts that use it inconsistently make this fragile; we map it to ṁ.
var standardSpecial = []struct {
	offset rune
	latin  string
}{
	{0x01, "m̐"}, // candrabindu
	{0x02, "ṁ"}, // anusvara
	{0x03, "ḥ"}, // visarga
	{0x3D, "ʼ"}, // avagraha
}

// newStandardScript builds a Script with the canonical Brahmic offsets
// added to base. Caller may delete/override entries afterwards.
func newStandardScript(name string, base rune) *Script {
	sc := &Script{
		Name:             name,
		BlockStart:       base,
		BlockEnd:         base + 0x7F,
		Virama:           []rune{base + 0x4D},
		DigitStart:       base + 0x66,
		ConsonantBase:    map[rune]string{},
		IndependentVowel: map[rune]string{},
		VowelSign:        map[rune]string{},
		Special:          map[rune]string{},
	}
	for _, e := range standardIndependentVowels {
		sc.IndependentVowel[base+e.offset] = e.latin
	}
	for _, e := range standardConsonants {
		sc.ConsonantBase[base+e.offset] = e.latin
	}
	for _, e := range standardVowelSigns {
		sc.VowelSign[base+e.offset] = e.latin
	}
	for _, e := range standardSpecial {
		sc.Special[base+e.offset] = e.latin
	}
	return sc
}
