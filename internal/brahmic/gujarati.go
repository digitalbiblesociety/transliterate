package brahmic

// Gujarati (U+0A80..U+0AFF).
//
// Layout is highly parallel to Devanagari; the main deviations are the
// absent short e/o at the +0x0E/+0x12 slots and a single ઑ (U+0A91) for
// "ô" used in English loans.
var Gujarati = func() *Script {
	sc := newStandardScript("gujarati", 0x0A80)
	delete(sc.IndependentVowel, 0x0A8E)
	delete(sc.IndependentVowel, 0x0A92)
	delete(sc.VowelSign, 0x0AC6)
	delete(sc.VowelSign, 0x0ACA)
	sc.IndependentVowel[0x0A8F] = "e"
	sc.IndependentVowel[0x0A93] = "o"
	sc.VowelSign[0x0AC7] = "e"
	sc.VowelSign[0x0ACB] = "o"
	sc.IndependentVowel[0x0A8D] = "ê"
	sc.IndependentVowel[0x0A91] = "ô"
	sc.VowelSign[0x0AC5] = "ê"
	sc.VowelSign[0x0AC9] = "ô"
	delete(sc.ConsonantBase, 0x0AB1)
	delete(sc.ConsonantBase, 0x0AB4)
	return sc
}()
