package brahmic

// Newa / Prachalit Nepal (U+11400..U+1147F).
//
// Used for Nepal Bhasa (Newari, ~1M speakers in the Kathmandu Valley)
// and historically for Classical Newar literature and Buddhist /
// Hindu manuscripts. A small Christian community uses translated
// materials. Codepoints derived from Aksharamukha (MIT).
//
// Layout note: vowels begin at +0x00 like Modi. The consonant grid is
// non-consecutive — six Newa-specific additional letters (NHA, etc.)
// sit between the canonical positions; those are left unmapped here
// and drop silently as in-block unhandled runes.
var Newa = func() *Script {
	sc := &Script{
		Name:             "newa",
		BlockStart:       0x11400,
		BlockEnd:         0x1147F,
		Virama:           []rune{0x11442},
		DigitStart:       0x11450,
		IndependentVowel: map[rune]string{},
		ConsonantBase:    map[rune]string{},
		VowelSign:        map[rune]string{},
		Special:          map[rune]string{},
	}

	sc.IndependentVowel[0x11400] = "a"
	sc.IndependentVowel[0x11401] = "ā"
	sc.IndependentVowel[0x11402] = "i"
	sc.IndependentVowel[0x11403] = "ī"
	sc.IndependentVowel[0x11404] = "u"
	sc.IndependentVowel[0x11405] = "ū"
	sc.IndependentVowel[0x11406] = "r̥"
	sc.IndependentVowel[0x11407] = "r̥̄"
	sc.IndependentVowel[0x11408] = "l̥"
	sc.IndependentVowel[0x11409] = "l̥̄"
	sc.IndependentVowel[0x1140A] = "e"
	sc.IndependentVowel[0x1140B] = "ai"
	sc.IndependentVowel[0x1140C] = "o"
	sc.IndependentVowel[0x1140D] = "au"

	// Consonants — canonical Brahmic order, but non-consecutive in the
	// block (positions 0x11413, 0x11419, 0x11424, 0x1142A, 0x1142D,
	// 0x1142F are Newa-specific extras that drop silently).
	cons := map[rune]string{
		0x1140E: "k", 0x1140F: "kh", 0x11410: "g", 0x11411: "gh", 0x11412: "ṅ",
		0x11414: "c", 0x11415: "ch", 0x11416: "j", 0x11417: "jh", 0x11418: "ñ",
		0x1141A: "ṭ", 0x1141B: "ṭh", 0x1141C: "ḍ", 0x1141D: "ḍh", 0x1141E: "ṇ",
		0x1141F: "t", 0x11420: "th", 0x11421: "d", 0x11422: "dh", 0x11423: "n",
		0x11425: "p", 0x11426: "ph", 0x11427: "b", 0x11428: "bh", 0x11429: "m",
		0x1142B: "y", 0x1142C: "r",
		0x1142E: "l",
		0x11430: "v",
		0x11431: "ś", 0x11432: "ṣ", 0x11433: "s", 0x11434: "h",
	}
	for r, v := range cons {
		sc.ConsonantBase[r] = v
	}

	sc.VowelSign[0x11435] = "ā"
	sc.VowelSign[0x11436] = "i"
	sc.VowelSign[0x11437] = "ī"
	sc.VowelSign[0x11438] = "u"
	sc.VowelSign[0x11439] = "ū"
	sc.VowelSign[0x1143A] = "r̥"
	sc.VowelSign[0x1143B] = "r̥̄"
	sc.VowelSign[0x1143C] = "l̥"
	sc.VowelSign[0x1143D] = "l̥̄"
	sc.VowelSign[0x1143E] = "e"
	sc.VowelSign[0x1143F] = "ai"
	sc.VowelSign[0x11440] = "o"
	sc.VowelSign[0x11441] = "au"

	sc.Special[0x11443] = "m̐" // candrabindu
	sc.Special[0x11444] = "ṁ" // anusvara
	sc.Special[0x11445] = "ḥ" // visarga

	return sc
}()
