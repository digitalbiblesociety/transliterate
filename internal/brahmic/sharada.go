package brahmic

// Sharada (U+11180..U+111DF).
//
// Used historically for Kashmiri, especially Śaiva tantric and
// philosophical texts (8th-19th century). Limited modern use; primarily
// of scholarly interest. Codepoint romanizations follow Aksharamukha's
// MIT-licensed mappings.
//
// Layout note: vowels start at +0x03 (not +0x05 like Devanagari), so the
// Script is constructed explicitly.
var Sharada = func() *Script {
	sc := &Script{
		Name:             "shrd",
		BlockStart:       0x11180,
		BlockEnd:         0x111DF,
		Virama:           []rune{0x111C0},
		DigitStart:       0x111D0,
		IndependentVowel: map[rune]string{},
		ConsonantBase:    map[rune]string{},
		VowelSign:        map[rune]string{},
		Special:          map[rune]string{},
	}

	sc.Special[0x11180] = "m̐" // 𑆀 candrabindu
	sc.Special[0x11181] = "ṁ" // 𑆁 anusvara
	sc.Special[0x11182] = "ḥ" // 𑆂 visarga

	sc.IndependentVowel[0x11183] = "a"
	sc.IndependentVowel[0x11184] = "ā"
	sc.IndependentVowel[0x11185] = "i"
	sc.IndependentVowel[0x11186] = "ī"
	sc.IndependentVowel[0x11187] = "u"
	sc.IndependentVowel[0x11188] = "ū"
	sc.IndependentVowel[0x11189] = "r̥"
	sc.IndependentVowel[0x1118A] = "r̥̄"
	sc.IndependentVowel[0x1118B] = "l̥"
	sc.IndependentVowel[0x1118C] = "l̥̄"
	sc.IndependentVowel[0x1118D] = "e"
	sc.IndependentVowel[0x1118E] = "ai"
	sc.IndependentVowel[0x1118F] = "o"
	sc.IndependentVowel[0x11190] = "au"

	sc.ConsonantBase[0x11191] = "k"
	sc.ConsonantBase[0x11192] = "kh"
	sc.ConsonantBase[0x11193] = "g"
	sc.ConsonantBase[0x11194] = "gh"
	sc.ConsonantBase[0x11195] = "ṅ"
	sc.ConsonantBase[0x11196] = "c"
	sc.ConsonantBase[0x11197] = "ch"
	sc.ConsonantBase[0x11198] = "j"
	sc.ConsonantBase[0x11199] = "jh"
	sc.ConsonantBase[0x1119A] = "ñ"
	sc.ConsonantBase[0x1119B] = "ṭ"
	sc.ConsonantBase[0x1119C] = "ṭh"
	sc.ConsonantBase[0x1119D] = "ḍ"
	sc.ConsonantBase[0x1119E] = "ḍh"
	sc.ConsonantBase[0x1119F] = "ṇ"
	sc.ConsonantBase[0x111A0] = "t"
	sc.ConsonantBase[0x111A1] = "th"
	sc.ConsonantBase[0x111A2] = "d"
	sc.ConsonantBase[0x111A3] = "dh"
	sc.ConsonantBase[0x111A4] = "n"
	sc.ConsonantBase[0x111A5] = "p"
	sc.ConsonantBase[0x111A6] = "ph"
	sc.ConsonantBase[0x111A7] = "b"
	sc.ConsonantBase[0x111A8] = "bh"
	sc.ConsonantBase[0x111A9] = "m"
	sc.ConsonantBase[0x111AA] = "y"
	sc.ConsonantBase[0x111AB] = "r"
	sc.ConsonantBase[0x111AC] = "l"
	sc.ConsonantBase[0x111AD] = "ḷ" // LLA
	sc.ConsonantBase[0x111AE] = "v"
	sc.ConsonantBase[0x111AF] = "ś"
	sc.ConsonantBase[0x111B0] = "ṣ"
	sc.ConsonantBase[0x111B1] = "s"
	sc.ConsonantBase[0x111B2] = "h"

	sc.VowelSign[0x111B3] = "ā"
	sc.VowelSign[0x111B4] = "i"
	sc.VowelSign[0x111B5] = "ī"
	sc.VowelSign[0x111B6] = "u"
	sc.VowelSign[0x111B7] = "ū"
	sc.VowelSign[0x111B8] = "r̥"
	sc.VowelSign[0x111B9] = "r̥̄"
	sc.VowelSign[0x111BA] = "l̥"
	sc.VowelSign[0x111BB] = "e"
	sc.VowelSign[0x111BC] = "ai"
	sc.VowelSign[0x111BD] = "o"
	sc.VowelSign[0x111BE] = "au"

	return sc
}()
