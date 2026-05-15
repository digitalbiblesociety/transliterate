package brahmic

// Balinese (U+1B00..U+1B7F).
//
// An Indic-derived Abugida used for the Balinese language and Kawi
// (Old Javanese) liturgical/literary texts. The block deviates from the
// canonical Devanagari layout enough that we describe it explicitly.
//
// Known simplifications:
//   - Sign Ulu Ricem (U+1B00) and other supra-segmental modifiers are
//     dropped; they carry tone-like nuance that doesn't survive Latin.
//   - Pepet vowel (ə) and its lengthened form are both rendered "ə".
var Balinese = func() *Script {
	sc := &Script{
		Name:             "balinese",
		BlockStart:       0x1B00,
		BlockEnd:         0x1B7F,
		Virama:           []rune{0x1B44}, // ◌᭄ adeg-adeg
		DigitStart:       0x1B50,
		IndependentVowel: map[rune]string{},
		ConsonantBase:    map[rune]string{},
		VowelSign:        map[rune]string{},
		Special:          map[rune]string{},
	}

	// Independent vowels.
	sc.IndependentVowel[0x1B05] = "a"
	sc.IndependentVowel[0x1B06] = "ā"
	sc.IndependentVowel[0x1B07] = "i"
	sc.IndependentVowel[0x1B08] = "ī"
	sc.IndependentVowel[0x1B09] = "u"
	sc.IndependentVowel[0x1B0A] = "ū"
	sc.IndependentVowel[0x1B0B] = "r̥"
	sc.IndependentVowel[0x1B0C] = "r̥̄"
	sc.IndependentVowel[0x1B0D] = "l̥"
	sc.IndependentVowel[0x1B0E] = "l̥̄"
	sc.IndependentVowel[0x1B0F] = "e"
	sc.IndependentVowel[0x1B10] = "ai"
	sc.IndependentVowel[0x1B11] = "o"
	sc.IndependentVowel[0x1B12] = "au"

	// Consonants — canonical Brahmic order.
	sc.ConsonantBase[0x1B13] = "k"
	sc.ConsonantBase[0x1B14] = "kh"
	sc.ConsonantBase[0x1B15] = "g"
	sc.ConsonantBase[0x1B16] = "gh"
	sc.ConsonantBase[0x1B17] = "ṅ"
	sc.ConsonantBase[0x1B18] = "c"
	sc.ConsonantBase[0x1B19] = "ch"
	sc.ConsonantBase[0x1B1A] = "j"
	sc.ConsonantBase[0x1B1B] = "jh"
	sc.ConsonantBase[0x1B1C] = "ñ"
	sc.ConsonantBase[0x1B1D] = "ṭ"
	sc.ConsonantBase[0x1B1E] = "ṭh"
	sc.ConsonantBase[0x1B1F] = "ḍ"
	sc.ConsonantBase[0x1B20] = "ḍh"
	sc.ConsonantBase[0x1B21] = "ṇ"
	sc.ConsonantBase[0x1B22] = "t"
	sc.ConsonantBase[0x1B23] = "th"
	sc.ConsonantBase[0x1B24] = "d"
	sc.ConsonantBase[0x1B25] = "dh"
	sc.ConsonantBase[0x1B26] = "n"
	sc.ConsonantBase[0x1B27] = "p"
	sc.ConsonantBase[0x1B28] = "ph"
	sc.ConsonantBase[0x1B29] = "b"
	sc.ConsonantBase[0x1B2A] = "bh"
	sc.ConsonantBase[0x1B2B] = "m"
	sc.ConsonantBase[0x1B2C] = "y"
	sc.ConsonantBase[0x1B2D] = "r"
	sc.ConsonantBase[0x1B2E] = "l"
	sc.ConsonantBase[0x1B2F] = "w"
	sc.ConsonantBase[0x1B30] = "ś"
	sc.ConsonantBase[0x1B31] = "ṣ"
	sc.ConsonantBase[0x1B32] = "s"
	sc.ConsonantBase[0x1B33] = "h"

	// Vowel signs (matras).
	sc.VowelSign[0x1B35] = "ā"
	sc.VowelSign[0x1B36] = "i"
	sc.VowelSign[0x1B37] = "ī"
	sc.VowelSign[0x1B38] = "u"
	sc.VowelSign[0x1B39] = "ū"
	sc.VowelSign[0x1B3A] = "r̥"
	sc.VowelSign[0x1B3B] = "r̥̄"
	sc.VowelSign[0x1B3C] = "l̥"
	sc.VowelSign[0x1B3D] = "l̥̄"
	sc.VowelSign[0x1B3E] = "e"
	sc.VowelSign[0x1B3F] = "ai"
	sc.VowelSign[0x1B40] = "o"
	sc.VowelSign[0x1B41] = "au"
	sc.VowelSign[0x1B42] = "ə" // pepet

	// Special / anusvara-family.
	sc.Special[0x1B01] = "m̐" // sign ulu candra (chandrabindu)
	sc.Special[0x1B02] = "ṁ" // sign cecek (anusvara)
	sc.Special[0x1B03] = "r" // sign surang (r-final)
	sc.Special[0x1B04] = "ḥ" // sign bisah (visarga)
	return sc
}()
