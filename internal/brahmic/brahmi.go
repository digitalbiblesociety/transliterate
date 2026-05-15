package brahmic

// Brahmi (U+11000..U+1107F).
//
// Ancient Indian script and the ancestor of every script in this family.
// Used for the Aśokan edicts (3rd century BCE) and early Buddhist/Jain
// inscriptions. Codepoint romanizations follow Aksharamukha's MIT-
// licensed mappings (ISO 15919-aligned).
//
// The block layout deviates from the canonical Devanagari grid: the
// candrabindu sits at offset +0x00 instead of +0x01, so we build the
// Script explicitly rather than via newStandardScript.
//
// Known scope:
//   - Modern decimal digits at U+11066..U+1106F are mapped to 0-9.
//   - Older Brahmi number letters at U+11052..U+11065 (1, 2, 3, …, 10,
//     20, …, 100, 1000) are not mapped; they'd need place-value
//     interpretation. They drop silently.
var Brahmi = func() *Script {
	sc := &Script{
		Name:             "brah",
		BlockStart:       0x11000,
		BlockEnd:         0x1107F,
		Virama:           []rune{0x11046},
		DigitStart:       0x11066,
		IndependentVowel: map[rune]string{},
		ConsonantBase:    map[rune]string{},
		VowelSign:        map[rune]string{},
		Special:          map[rune]string{},
	}

	sc.Special[0x11000] = "m̐"  // 𑀀 candrabindu
	sc.Special[0x11001] = "ṁ"  // 𑀁 anusvara
	sc.Special[0x11002] = "ḥ"  // 𑀂 visarga
	sc.Special[0x11003] = "ẖ"  // jihvamuliya (rare)
	sc.Special[0x11004] = "ḫ"  // upadhmaniya (rare)

	sc.IndependentVowel[0x11005] = "a"
	sc.IndependentVowel[0x11006] = "ā"
	sc.IndependentVowel[0x11007] = "i"
	sc.IndependentVowel[0x11008] = "ī"
	sc.IndependentVowel[0x11009] = "u"
	sc.IndependentVowel[0x1100A] = "ū"
	sc.IndependentVowel[0x1100B] = "r̥"
	sc.IndependentVowel[0x1100C] = "r̥̄"
	sc.IndependentVowel[0x1100D] = "l̥"
	sc.IndependentVowel[0x1100E] = "l̥̄"
	sc.IndependentVowel[0x1100F] = "e"
	sc.IndependentVowel[0x11010] = "ai"
	sc.IndependentVowel[0x11011] = "o"
	sc.IndependentVowel[0x11012] = "au"

	sc.ConsonantBase[0x11013] = "k"
	sc.ConsonantBase[0x11014] = "kh"
	sc.ConsonantBase[0x11015] = "g"
	sc.ConsonantBase[0x11016] = "gh"
	sc.ConsonantBase[0x11017] = "ṅ"
	sc.ConsonantBase[0x11018] = "c"
	sc.ConsonantBase[0x11019] = "ch"
	sc.ConsonantBase[0x1101A] = "j"
	sc.ConsonantBase[0x1101B] = "jh"
	sc.ConsonantBase[0x1101C] = "ñ"
	sc.ConsonantBase[0x1101D] = "ṭ"
	sc.ConsonantBase[0x1101E] = "ṭh"
	sc.ConsonantBase[0x1101F] = "ḍ"
	sc.ConsonantBase[0x11020] = "ḍh"
	sc.ConsonantBase[0x11021] = "ṇ"
	sc.ConsonantBase[0x11022] = "t"
	sc.ConsonantBase[0x11023] = "th"
	sc.ConsonantBase[0x11024] = "d"
	sc.ConsonantBase[0x11025] = "dh"
	sc.ConsonantBase[0x11026] = "n"
	sc.ConsonantBase[0x11027] = "p"
	sc.ConsonantBase[0x11028] = "ph"
	sc.ConsonantBase[0x11029] = "b"
	sc.ConsonantBase[0x1102A] = "bh"
	sc.ConsonantBase[0x1102B] = "m"
	sc.ConsonantBase[0x1102C] = "y"
	sc.ConsonantBase[0x1102D] = "r"
	sc.ConsonantBase[0x1102E] = "l"
	sc.ConsonantBase[0x1102F] = "v"
	sc.ConsonantBase[0x11030] = "ś"
	sc.ConsonantBase[0x11031] = "ṣ"
	sc.ConsonantBase[0x11032] = "s"
	sc.ConsonantBase[0x11033] = "h"
	sc.ConsonantBase[0x11034] = "ḷ"

	sc.VowelSign[0x11038] = "ā"
	sc.VowelSign[0x11039] = "i"
	sc.VowelSign[0x1103A] = "ī"
	sc.VowelSign[0x1103B] = "u"
	sc.VowelSign[0x1103C] = "ū"
	sc.VowelSign[0x1103D] = "r̥"
	sc.VowelSign[0x1103E] = "r̥̄"
	sc.VowelSign[0x1103F] = "l̥"
	sc.VowelSign[0x11040] = "l̥̄"
	sc.VowelSign[0x11041] = "e"
	sc.VowelSign[0x11042] = "ai"
	sc.VowelSign[0x11043] = "o"
	sc.VowelSign[0x11044] = "au"

	return sc
}()
