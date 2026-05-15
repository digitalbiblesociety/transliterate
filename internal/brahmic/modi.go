package brahmic

// Modi (U+11600..U+1165F).
//
// Used historically for Marathi business correspondence and
// administrative records from the 17th to mid-20th century. Some
// Christian missionary materials exist in this script. Codepoints
// derived from Aksharamukha (MIT).
//
// Layout note: vowels begin at +0x00 (not +0x05 like Devanagari).
var Modi = func() *Script {
	sc := &Script{
		Name:             "modi",
		BlockStart:       0x11600,
		BlockEnd:         0x1165F,
		Virama:           []rune{0x1163F},
		DigitStart:       0x11650,
		IndependentVowel: map[rune]string{},
		ConsonantBase:    map[rune]string{},
		VowelSign:        map[rune]string{},
		Special:          map[rune]string{},
	}

	// Independent vowels (canonical order A, AA, I, II, U, UU, R, RR, L, LL, E, AI, O, AU).
	sc.IndependentVowel[0x11600] = "a"
	sc.IndependentVowel[0x11601] = "ā"
	sc.IndependentVowel[0x11602] = "i"
	sc.IndependentVowel[0x11603] = "ī"
	sc.IndependentVowel[0x11604] = "u"
	sc.IndependentVowel[0x11605] = "ū"
	sc.IndependentVowel[0x11606] = "r̥"
	sc.IndependentVowel[0x11607] = "r̥̄"
	sc.IndependentVowel[0x11608] = "l̥"
	sc.IndependentVowel[0x11609] = "l̥̄"
	sc.IndependentVowel[0x1160A] = "e"
	sc.IndependentVowel[0x1160B] = "ai"
	sc.IndependentVowel[0x1160C] = "o"
	sc.IndependentVowel[0x1160D] = "au"

	// Consonants 0x1160E..0x1162E (33 in canonical order) + LLA at 0x1162F.
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
		sc.ConsonantBase[0x1160E+rune(i)] = v
	}
	sc.ConsonantBase[0x1162F] = "ḷ" // LLA

	// Vowel signs (matras) at 0x11630..0x1163C.
	sc.VowelSign[0x11630] = "ā"
	sc.VowelSign[0x11631] = "i"
	sc.VowelSign[0x11632] = "ī"
	sc.VowelSign[0x11633] = "u"
	sc.VowelSign[0x11634] = "ū"
	sc.VowelSign[0x11635] = "r̥"
	sc.VowelSign[0x11636] = "r̥̄"
	sc.VowelSign[0x11637] = "l̥"
	sc.VowelSign[0x11638] = "l̥̄"
	sc.VowelSign[0x11639] = "e"
	sc.VowelSign[0x1163A] = "ai"
	sc.VowelSign[0x1163B] = "o"
	sc.VowelSign[0x1163C] = "au"

	sc.Special[0x1163D] = "ṁ" // anusvara
	sc.Special[0x1163E] = "ḥ" // visarga
	sc.Special[0x11640] = "m̐" // ardhacandra

	return sc
}()
