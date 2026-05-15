package brahmic

// Tirhuta (U+11480..U+114DF).
//
// Used historically for Maithili (~50M speakers, Bihar and Nepal) — the
// language has a written tradition extending back to the 14th century,
// and the script is being revived for cultural / educational use.
// Codepoints derived from Aksharamukha (MIT).
//
// Layout note: vowels begin at +0x01 (U+11480 itself is reserved). The
// canonical consonant grid runs U+1148F..U+114AF.
var Tirhuta = func() *Script {
	sc := &Script{
		Name:             "tirh",
		BlockStart:       0x11480,
		BlockEnd:         0x114DF,
		Virama:           []rune{0x114C2},
		DigitStart:       0x114D0,
		IndependentVowel: map[rune]string{},
		ConsonantBase:    map[rune]string{},
		VowelSign:        map[rune]string{},
		Special:          map[rune]string{},
	}

	sc.IndependentVowel[0x11481] = "a"
	sc.IndependentVowel[0x11482] = "ā"
	sc.IndependentVowel[0x11483] = "i"
	sc.IndependentVowel[0x11484] = "ī"
	sc.IndependentVowel[0x11485] = "u"
	sc.IndependentVowel[0x11486] = "ū"
	sc.IndependentVowel[0x11487] = "r̥"
	sc.IndependentVowel[0x11488] = "r̥̄"
	sc.IndependentVowel[0x11489] = "l̥"
	sc.IndependentVowel[0x1148A] = "l̥̄"
	sc.IndependentVowel[0x1148B] = "e"
	sc.IndependentVowel[0x1148C] = "ai"
	sc.IndependentVowel[0x1148D] = "o"
	sc.IndependentVowel[0x1148E] = "au"

	consonants := []string{
		"k", "kh", "g", "gh", "ṅ",
		"c", "ch", "j", "jh", "ñ",
		"ṭ", "ṭh", "ḍ", "ḍh", "ṇ",
		"t", "th", "d", "dh", "n",
		"p", "ph", "b", "bh", "m",
		"y", "r", "l", "v",
		"ś", "ṣ", "s", "h",
	}
	for i, v := range consonants {
		sc.ConsonantBase[0x1148F+rune(i)] = v
	}

	// Vowel signs. AA..LL run consecutively, then a gap at 0x114BA / 0x114BD.
	sc.VowelSign[0x114B0] = "ā"
	sc.VowelSign[0x114B1] = "i"
	sc.VowelSign[0x114B2] = "ī"
	sc.VowelSign[0x114B3] = "u"
	sc.VowelSign[0x114B4] = "ū"
	sc.VowelSign[0x114B5] = "r̥"
	sc.VowelSign[0x114B6] = "r̥̄"
	sc.VowelSign[0x114B7] = "l̥"
	sc.VowelSign[0x114B8] = "l̥̄"
	sc.VowelSign[0x114B9] = "e"
	sc.VowelSign[0x114BB] = "ai"
	sc.VowelSign[0x114BC] = "o"
	sc.VowelSign[0x114BE] = "au"

	sc.Special[0x114BF] = "m̐" // candrabindu
	sc.Special[0x114C0] = "ṁ" // anusvara
	sc.Special[0x114C1] = "ḥ" // visarga
	sc.Special[0x114C4] = "ʼ" // avagraha-like

	return sc
}()
