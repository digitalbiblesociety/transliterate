// Package ethiopic provides BGN/PCGN-style transliteration of the
// Ethiopic (Ge'ez) script to the Latin alphabet. The script is used for
// Amharic, Tigrinya, Ge'ez (liturgical), Tigre, and several other
// languages of the Horn of Africa.
//
// Ethiopic is a syllabary: each glyph encodes a consonant + one of seven
// vowel "orders." Regular consonant rows are 8 codepoints wide; cells
// 0..6 hold orders 1..7 (the 7 standard vowels) and cell 7 is usually
// reserved. Labialized rows (QWA, KWA, GWA, etc.) live at separate
// codepoint ranges with their own irregular layouts and are handled
// individually.
package ethi

import "strings"

const (
	BlockStart rune = 0x1200
	BlockEnd   rune = 0x137F
)

// vowelOrders is the suffix appended to a consonant base for each of the
// seven Ethiopic vowel orders.
//
//	order 0 ("ä")  — 1st order
//	order 1 ("u")  — 2nd order
//	order 2 ("i")  — 3rd order
//	order 3 ("a")  — 4th order
//	order 4 ("e")  — 5th order
//	order 5 ("ə")  — 6th order (silent at word boundary)
//	order 6 ("o")  — 7th order
var vowelOrders = [7]string{"ä", "u", "i", "a", "e", "ə", "o"}

// consonantRows maps the first-codepoint of each 8-cell row to its
// consonant base. Codepoints verified against Unicode 15.1 names.
var consonantRows = map[rune]string{
	0x1200: "h",  // ሀ HA
	0x1208: "l",  // ለ LA
	0x1210: "ḥ",  // ሐ HHA (pharyngeal h)
	0x1218: "m",  // መ MA
	0x1220: "ś",  // ሠ SZA
	0x1228: "r",  // ረ RA
	0x1230: "s",  // ሰ SA
	0x1238: "š",  // ሸ SHA
	0x1240: "q",  // ቀ QA
	0x1250: "q̱",  // ቐ QHA (Tigrinya)
	0x1260: "b",  // በ BA
	0x1268: "v",  // ቨ VA
	0x1270: "t",  // ተ TA
	0x1278: "č",  // ቸ CA
	0x1280: "ḫ",  // ኀ XA
	0x1290: "n",  // ነ NA
	0x1298: "ñ",  // ኘ NYA
	0x12A0: "ʾ",  // አ GLOTTAL A (alef)
	0x12A8: "k",  // ከ KA
	0x12B8: "x",  // ኸ KXA (Tigrinya)
	0x12C8: "w",  // ወ WA
	0x12D0: "ʿ",  // ዐ PHARYNGEAL A (ayn)
	0x12D8: "z",  // ዘ ZA
	0x12E0: "ž",  // ዠ ZHA
	0x12E8: "y",  // የ YA
	0x12F0: "d",  // ደ DA
	0x12F8: "ḏ",  // ዸ DDA
	0x1300: "ǧ",  // ጀ JA
	0x1308: "g",  // ገ GA
	0x1318: "ġ",  // ጘ GGA
	0x1320: "ṭ",  // ጠ THA (emphatic t)
	0x1328: "č̣",  // ጨ CHA (emphatic ch)
	0x1330: "ṗ",  // ጰ PHA (emphatic p)
	0x1338: "ṣ",  // ጸ TSA
	0x1340: "ḍ",  // ፀ TZA
	0x1348: "f",  // ፈ FA
	0x1350: "p",  // ፐ PA
}

// labializedTable holds the labialized (w-glide) syllable rows whose
// codepoints don't follow the regular 8-cell grid. Each row has 5
// glyphs at offsets +0, +2, +3, +4, +5 (cells 1, 6, 7 are reserved).
var labializedTable = map[rune]string{
	// QWA series (after QA).
	0x1248: "qwä", 0x124A: "qwi", 0x124B: "qwa", 0x124C: "qwe", 0x124D: "qwə",
	// QHWA series (after QHA).
	0x1258: "q̱wä", 0x125A: "q̱wi", 0x125B: "q̱wa", 0x125C: "q̱we", 0x125D: "q̱wə",
	// XWA series (after XA).
	0x1288: "ḫwä", 0x128A: "ḫwi", 0x128B: "ḫwa", 0x128C: "ḫwe", 0x128D: "ḫwə",
	// KWA series (after KA).
	0x12B0: "kwä", 0x12B2: "kwi", 0x12B3: "kwa", 0x12B4: "kwe", 0x12B5: "kwə",
	// KXWA series (after KXA).
	0x12C0: "xwä", 0x12C2: "xwi", 0x12C3: "xwa", 0x12C4: "xwe", 0x12C5: "xwə",
	// GWA series (after GA).
	0x1310: "gwä", 0x1312: "gwi", 0x1313: "gwa", 0x1314: "gwe", 0x1315: "gwə",
}

// punctuation covers Ethiopic separators and end-of-sentence marks.
var punctuation = map[rune]string{
	0x1361: " ",   // ፡ word-space (between words)
	0x1362: ".",   // ።
	0x1363: ",",   // ፣
	0x1364: ";",   // ፤
	0x1365: ":",   // ፥
	0x1366: ":-",  // ፦
	0x1367: "?",   // ፧
	0x1368: "",    // ፨ paragraph separator
}

// Transliterate returns the BGN/PCGN romanization of s. Runes outside
// the Ethiopic block pass through unchanged. The 6th-order vowel (ə) is
// dropped when it falls at a word boundary, matching conventional
// pronunciation: ሰላም → "sälam" not "sälamə", but ምድር → "mədər" because
// the medial 6th-order is preserved.
func Transliterate(s string) string {
	if s == "" {
		return s
	}
	rs := []rune(s)
	var b strings.Builder
	b.Grow(len(s) * 2)
	for i, r := range rs {
		// Labialized syllable.
		if v, ok := labializedTable[r]; ok {
			b.WriteString(v)
			continue
		}
		// Punctuation.
		if v, ok := punctuation[r]; ok {
			b.WriteString(v)
			continue
		}
		// Digits 1-9 (Ethiopic positional 10/20/... numerals dropped).
		if r >= 0x1369 && r <= 0x1371 {
			b.WriteRune('1' + (r - 0x1369))
			continue
		}
		// Regular 7-order consonant grid.
		if r >= BlockStart && r <= BlockEnd {
			order := (r - 0x1200) % 8
			rowStart := r - order
			if base, ok := consonantRows[rowStart]; ok && order <= 6 {
				if order == 5 && isWordEnd(rs, i) {
					b.WriteString(base)
				} else {
					b.WriteString(base)
					b.WriteString(vowelOrders[order])
				}
				continue
			}
			// In-block but unmapped: positional numerals (10/20/...
			// 10000), archaic glyphs, reserved cells. Drop rather than
			// emit raw Ethiopic.
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// isWordEnd reports whether the glyph at rs[i] is the last Ethiopic
// glyph of its word — i.e., the next rune is not an Ethiopic
// consonant-vowel glyph. Used to drop the silent 6th-order ə.
func isWordEnd(rs []rune, i int) bool {
	if i+1 >= len(rs) {
		return true
	}
	next := rs[i+1]
	if next < BlockStart || next > BlockEnd {
		return true
	}
	// Ethiopic punctuation marks count as word boundaries.
	if next >= 0x1361 && next <= 0x1368 {
		return true
	}
	return false
}

// Contains reports whether s has at least one Ethiopic-block rune.
func Contains(s string) bool {
	for _, r := range s {
		if r >= BlockStart && r <= BlockEnd {
			return true
		}
	}
	return false
}
