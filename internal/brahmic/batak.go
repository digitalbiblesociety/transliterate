package brahmic

// Batak (U+1BC0..U+1BFF).
//
// Covers the surat Batak family used for Toba, Karo, Mandailing, Pakpak/
// Dairi, and Simalungun Batak languages of North Sumatra. Unicode encodes
// many language-specific letter variants (e.g. "Karo ba", "Simalungun ha")
// as distinct codepoints; we map each variant to the same Latin base so
// downstream consumers see consistent output regardless of dialect glyph.
//
// Two virama-like "killer" marks both suppress the inherent /a/: pangolat
// (U+1BF2, mostly Toba) and pangonangon (U+1BF3, mostly the other
// varieties).
var Batak = func() *Script {
	sc := &Script{
		Name:             "batak",
		BlockStart:       0x1BC0,
		BlockEnd:         0x1BFF,
		Virama:           []rune{0x1BF2, 0x1BF3}, // pangolat, pangonangon
		IndependentVowel: map[rune]string{},
		ConsonantBase:    map[rune]string{},
		VowelSign:        map[rune]string{},
		Special:          map[rune]string{},
	}

	sc.IndependentVowel[0x1BC0] = "a"
	sc.IndependentVowel[0x1BC1] = "a" // simalungun a
	sc.IndependentVowel[0x1BE4] = "i"
	sc.IndependentVowel[0x1BE5] = "u"

	// Consonants — variant letters collapse to the same Latin base.
	sc.ConsonantBase[0x1BC2] = "h"
	sc.ConsonantBase[0x1BC3] = "h" // simalungun ha
	sc.ConsonantBase[0x1BC4] = "h" // mandailing ha
	sc.ConsonantBase[0x1BC5] = "b"
	sc.ConsonantBase[0x1BC6] = "b" // karo ba
	sc.ConsonantBase[0x1BC7] = "p"
	sc.ConsonantBase[0x1BC8] = "p" // simalungun pa
	sc.ConsonantBase[0x1BC9] = "n"
	sc.ConsonantBase[0x1BCA] = "n" // mandailing na
	sc.ConsonantBase[0x1BCB] = "w"
	sc.ConsonantBase[0x1BCC] = "w" // simalungun wa
	sc.ConsonantBase[0x1BCD] = "w" // pakpak wa
	sc.ConsonantBase[0x1BCE] = "g"
	sc.ConsonantBase[0x1BCF] = "g" // simalungun ga
	sc.ConsonantBase[0x1BD0] = "j"
	sc.ConsonantBase[0x1BD1] = "d"
	sc.ConsonantBase[0x1BD2] = "r"
	sc.ConsonantBase[0x1BD3] = "r" // simalungun ra
	sc.ConsonantBase[0x1BD4] = "m"
	sc.ConsonantBase[0x1BD5] = "m" // simalungun ma
	sc.ConsonantBase[0x1BD6] = "t" // southern ta
	sc.ConsonantBase[0x1BD7] = "t" // northern ta
	sc.ConsonantBase[0x1BD8] = "s"
	sc.ConsonantBase[0x1BD9] = "s" // simalungun sa
	sc.ConsonantBase[0x1BDA] = "s" // mandailing sa
	sc.ConsonantBase[0x1BDB] = "y"
	sc.ConsonantBase[0x1BDC] = "y" // simalungun ya
	sc.ConsonantBase[0x1BDD] = "ṅ" // nga
	sc.ConsonantBase[0x1BDE] = "l"
	sc.ConsonantBase[0x1BDF] = "l" // simalungun la
	sc.ConsonantBase[0x1BE0] = "ñ" // nya
	sc.ConsonantBase[0x1BE1] = "c"
	sc.ConsonantBase[0x1BE2] = "nd" // nda (prenasalized)
	sc.ConsonantBase[0x1BE3] = "mb" // mba (prenasalized)

	sc.VowelSign[0x1BE7] = "e"
	sc.VowelSign[0x1BE8] = "e" // pakpak e
	sc.VowelSign[0x1BE9] = "ē" // ee
	sc.VowelSign[0x1BEA] = "i"
	sc.VowelSign[0x1BEB] = "i" // karo i
	sc.VowelSign[0x1BEC] = "o"
	sc.VowelSign[0x1BED] = "o" // karo o
	sc.VowelSign[0x1BEE] = "u"
	sc.VowelSign[0x1BEF] = "u" // u for simalungun sa

	sc.Special[0x1BE6] = ""  // tompi (modifier, dropped)
	sc.Special[0x1BF0] = "ŋ" // consonant sign ng
	sc.Special[0x1BF1] = "h" // consonant sign h
	return sc
}()
