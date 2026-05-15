package brahmic

// Oriya / Odia (U+0B00..U+0B7F).
//
// Notes:
//   - Mostly parallel layout. No short e/o at the canonical offsets.
//   - Oriya has its own ୱ (U+0B71) "wa" and ୟ (U+0B5F) "ẏa".
//   - Oriya digits at U+0B66..U+0B6F.
var Oriya = func() *Script {
	sc := newStandardScript("oriya", 0x0B00)
	delete(sc.IndependentVowel, 0x0B0E)
	delete(sc.IndependentVowel, 0x0B12)
	delete(sc.VowelSign, 0x0B46)
	delete(sc.VowelSign, 0x0B4A)
	sc.IndependentVowel[0x0B0F] = "e"
	sc.IndependentVowel[0x0B13] = "o"
	sc.VowelSign[0x0B47] = "e"
	sc.VowelSign[0x0B4B] = "o"
	sc.ConsonantBase[0x0B5C] = "ṛ"  // ଡ଼
	sc.ConsonantBase[0x0B5D] = "ṛh" // ଢ଼
	sc.ConsonantBase[0x0B5F] = "ẏ"  // ୟ
	sc.ConsonantBase[0x0B71] = "w"  // ୱ
	delete(sc.ConsonantBase, 0x0B31)
	delete(sc.ConsonantBase, 0x0B34)
	// Nukta combining (decomposed form).
	sc.Nukta = 0x0B3C
	sc.NuktaCombine = map[rune]string{
		0x0B21: "ṛ",  // ଡ + ◌଼
		0x0B22: "ṛh", // ଢ + ◌଼
		0x0B2F: "ẏ",  // ଯ + ◌଼
	}
	return sc
}()
